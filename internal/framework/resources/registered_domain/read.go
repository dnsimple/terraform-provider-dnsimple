package registered_domain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Bounds on retrying a transient registrar failure during a refresh.
//
// The budget follows the 30s convention the create and update paths already use, and is
// hardcoded rather than exposed as a `read` timeout: this retrying exists only until the
// API gives transient failures a retryable status code (dnsimple/dnsimple-app#36024), and
// a schema attribute added for it would outlive the reason for it.
//
// Backing off from 1s to 8s spends the budget on roughly six attempts, which covers a
// momentary blip without holding a large workspace open on a registrar outage that is not
// going to clear.
const (
	readRetryBudget       = 30 * time.Second
	readRetryInitialDelay = 1 * time.Second
	readRetryMaxDelay     = 8 * time.Second
)

// retryTransient runs an idempotent read against the DNSimple API, retrying it while it
// fails with a transient registrar-connection error and returning every other error
// unchanged on the first attempt.
//
// Read calls GetDomainRegistration, GetRegistrantChange, GetDomain, GetDnssec and
// GetDomainTransferLock in sequence and treats any error as fatal, so a single blip on
// any one of them aborts the whole plan. All five are reads with no side effects, which
// is what makes retrying them safe.
//
// The deadline is passed in rather than derived from the budget here, so that one refresh
// spends one budget across all five calls. Giving each call its own would let a single
// domain add five times the budget to a plan during a flapping outage, and a plan
// refreshes every registered domain in the workspace.
func retryTransient[T any](ctx context.Context, deadline time.Time, fn func() (T, error)) (T, error) {
	var result T

	err := utils.RetryWithBackoff(ctx, func() (error, bool) {
		var callErr error
		result, callErr = fn()

		// Suspend on anything that is not a transient registrar failure, so ordinary
		// errors still surface immediately.
		return callErr, callErr != nil && !utils.IsTransientRegistrarError(callErr)
	}, time.Until(deadline), readRetryInitialDelay, readRetryMaxDelay)

	return result, err
}

// warnDomainNoLongerRegistered reports that a registrar call was rejected because the
// domain has lapsed, without failing the refresh. Callers return immediately afterwards,
// which leaves prior state untouched.
//
// A lapsed domain must not fail the whole plan, which is the point of this handling. It
// must not silently drop the resource from state either: Terraform plans a fresh create
// for anything still in the configuration, and creating this resource re-registers the
// domain, which is billable and would run unattended under `apply -auto-approve`.
// Reporting and leaving state alone keeps the decision with the practitioner.
func warnDomainNoLongerRegistered(ctx context.Context, domain string, resp *resource.ReadResponse) {
	tflog.Warn(ctx, "registrar call rejected because the domain is no longer registered or has expired", map[string]interface{}{"domain": domain})

	resp.Diagnostics.AddWarning(
		fmt.Sprintf("could not refresh %s because it is no longer registered", domain),
		fmt.Sprintf("The DNSimple API reports that '%s' is not registered or has expired, so its registrar details could not be refreshed. Terraform state has been left as-is rather than dropping the resource, which would have planned a new registration. Renew or restore the domain to resume managing it, or remove the resource from your configuration and state if you no longer want it managed. Any plan produced while the domain stays lapsed is unlikely to apply cleanly.", domain),
	)
}

func (r *RegisteredDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// One budget for the whole refresh, shared by every call below.
	retryDeadline := time.Now().Add(readRetryBudget)

	domainRegistration, diags := getDomainRegistration(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var domainRegistrationResponse *dnsimple.DomainRegistrationResponse
	var err error
	if !domainRegistration.Id.IsNull() {
		domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
		domainRegistrationResponse, err = retryTransient(ctx, retryDeadline, func() (*dnsimple.DomainRegistrationResponse, error) {
			return r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), domainRegistrationId)
		})
		if err != nil {
			if utils.IsDomainNotRegisteredOrExpiredError(err) {
				warnDomainNoLongerRegistered(ctx, data.Name.ValueString(), resp)
				return
			}

			resp.Diagnostics.AddError(
				"failed to read DNSimple Domain Registration",
				fmt.Sprintf("Unable to read domain registration for domain '%s' (registration ID: %d): %s", data.Name.ValueString(), domainRegistration.Id.ValueInt64(), err.Error()),
			)
			return
		}
	}

	registrantChange, diags := getRegistrantChange(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var registrantChangeResponse *dnsimple.RegistrantChangeResponse
	if !registrantChange.Id.IsNull() {
		registrantChangeResponse, err = retryTransient(ctx, retryDeadline, func() (*dnsimple.RegistrantChangeResponse, error) {
			return r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(registrantChange.Id.ValueInt64()))
		})
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to read DNSimple Registrant Change",
				fmt.Sprintf("Unable to read registrant change with ID %d: %s", registrantChange.Id.ValueInt64(), err.Error()),
			)
			return
		}

		registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		data.RegistrantChange = registrantChangeObject
	}

	domainResponse, err := retryTransient(ctx, retryDeadline, func() (*dnsimple.DomainResponse, error) {
		return r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read DNSimple Domain",
			fmt.Sprintf("Unable to read domain '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}

	var dnssec *dnsimple.Dnssec
	var transferLock *dnsimple.DomainTransferLock

	if domainResponse.Data.State == consts.DomainStateRegistered {
		dnssecResponse, err := retryTransient(ctx, retryDeadline, func() (*dnsimple.DnssecResponse, error) {
			return r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, data.Name.ValueString())
		})
		if err != nil {
			if utils.IsDomainNotRegisteredOrExpiredError(err) {
				warnDomainNoLongerRegistered(ctx, data.Name.ValueString(), resp)
				return
			}

			resp.Diagnostics.AddError(
				"failed to read DNSimple Domain DNSSEC status",
				fmt.Sprintf("Unable to read DNSSEC status for domain '%s': %s", data.Name.ValueString(), err.Error()),
			)
			return
		}
		dnssec = dnssecResponse.Data

		transferLockResponse, err := retryTransient(ctx, retryDeadline, func() (*dnsimple.DomainTransferLockResponse, error) {
			return r.config.Client.Registrar.GetDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())
		})
		if err != nil {
			if utils.IsDomainNotRegisteredOrExpiredError(err) {
				warnDomainNoLongerRegistered(ctx, data.Name.ValueString(), resp)
				return
			}

			resp.Diagnostics.AddError(
				"failed to read DNSimple Domain Transfer Lock status",
				fmt.Sprintf("Unable to read transfer lock status for domain '%s': %s", data.Name.ValueString(), err.Error()),
			)
			return
		}
		transferLock = transferLockResponse.Data
	}

	if domainRegistrationResponse == nil {
		diags = r.updateModelFromAPIResponse(ctx, data, nil, domainResponse.Data, dnssec, transferLock)
	} else {
		diags = r.updateModelFromAPIResponse(ctx, data, domainRegistrationResponse.Data, domainResponse.Data, dnssec, transferLock)
	}

	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
