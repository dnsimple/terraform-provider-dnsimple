package registered_domain

import (
	"context"
	"fmt"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

const (
	RegistrationConverged          = "registration_converged"
	RegistrationConvergenceTimeout = "registration_converged_timeout"
	RegistrationFailed             = "registration_failed"
)

func (r *RegisteredDomainResource) setAutoRenewal(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting auto_renew_enabled to %t", data.AutoRenewEnabled.ValueBool()))

	if data.AutoRenewEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableDomainAutoRenewal(ctx, r.config.AccountID, data.Name.ValueString())
		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain Auto Renewal: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableDomainAutoRenewal(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain Auto Renewal: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) setWhoisPrivacy(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting whois_privacy_enabled to %t", data.WhoisPrivacyEnabled.ValueBool()))

	if data.WhoisPrivacyEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableWhoisPrivacy(ctx, r.config.AccountID, data.Name.ValueString())
		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain Whois Privacy: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableWhoisPrivacy(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain Whois Privacy: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) setDNSSEC(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting dnssec_enabled to %t", data.DNSSECEnabled.ValueBool()))

	if data.DNSSECEnabled.ValueBool() {
		_, err := r.config.Client.Domains.EnableDnssec(ctx, r.config.AccountID, data.Name.ValueString())

		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain DNSSEC: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Domains.DisableDnssec(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain DNSSEC: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) setTransferLock(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting transfer_lock_enabled to %t", data.TransferLockEnabled.ValueBool()))

	if data.TransferLockEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())

		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain transfer lock: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain transfer lock: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) updateModelFromAPIResponse(ctx context.Context, data *RegisteredDomainResourceModel, domainRegistration *dnsimple.DomainRegistration, domain *dnsimple.Domain, dnssec *dnsimple.Dnssec, transferLock *dnsimple.DomainTransferLock) diag.Diagnostics {
	if domainRegistration != nil {
		domainRegistrationObject, diags := r.domainRegistrationAPIResponseToObject(ctx, domainRegistration)

		if diags.HasError() {
			return diags
		}

		data.DomainRegistration = domainRegistrationObject
	}

	if domain != nil {
		data.Id = types.Int64Value(domain.ID)
		data.AutoRenewEnabled = types.BoolValue(domain.AutoRenew)
		data.WhoisPrivacyEnabled = types.BoolValue(domain.PrivateWhois)
		data.State = types.StringValue(domain.State)
		data.UnicodeName = types.StringValue(domain.UnicodeName)
		data.AccountId = types.Int64Value(domain.AccountID)
		data.ContactId = types.Int64Value(domain.RegistrantID)
		data.ExpiresAt = types.StringValue(domain.ExpiresAt)
		data.Name = types.StringValue(domain.Name)
	}

	if dnssec != nil {
		data.DNSSECEnabled = types.BoolValue(dnssec.Enabled)
	}

	if transferLock != nil {
		data.TransferLockEnabled = types.BoolValue(transferLock.Enabled)
	}

	return nil
}

func (r *RegisteredDomainResource) updateModelFromAPIResponsePartialCreate(ctx context.Context, data *RegisteredDomainResourceModel, domainRegistration *dnsimple.DomainRegistration, domain *dnsimple.Domain) *diag.Diagnostics {
	domainRegistrationObject, diags := r.domainRegistrationAPIResponseToObject(ctx, domainRegistration)

	if diags.HasError() {
		return &diags
	}

	data.DomainRegistration = domainRegistrationObject

	if domain != nil {
		data.Id = types.Int64Value(domain.ID)

		if data.AutoRenewEnabled.IsNull() {
			data.AutoRenewEnabled = types.BoolValue(domain.AutoRenew)
		}

		if data.WhoisPrivacyEnabled.IsNull() {
			data.WhoisPrivacyEnabled = types.BoolValue(domain.PrivateWhois)
		}

		data.State = types.StringValue(domain.State)
		data.UnicodeName = types.StringValue(domain.UnicodeName)
		data.AccountId = types.Int64Value(domain.AccountID)
		data.ContactId = types.Int64Value(domain.RegistrantID)
		data.ExpiresAt = types.StringValue(domain.ExpiresAt)
		data.Name = types.StringValue(domain.Name)
	}

	data.DNSSECEnabled = types.BoolValue(false)
	data.TransferLockEnabled = types.BoolValue(false)

	return nil
}

func (r *RegisteredDomainResource) domainRegistrationAPIResponseToObject(ctx context.Context, domainRegistration *dnsimple.DomainRegistration) (basetypes.ObjectValue, diag.Diagnostics) {
	domainRegistrationData := common.DomainRegistration{
		Id:     types.Int64Value(domainRegistration.ID),
		Period: types.Int64Value(int64(domainRegistration.Period)),
		State:  types.StringValue(domainRegistration.State),
	}

	return types.ObjectValueFrom(ctx, common.DomainRegistrationAttrType, domainRegistrationData)
}

func getTimeouts(ctx context.Context, model *RegisteredDomainResourceModel) (*common.Timeouts, diag.Diagnostics) {
	timeouts := &common.Timeouts{}
	diags := model.Timeouts.As(ctx, timeouts, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return timeouts, diags
}

func getDomainRegistration(ctx context.Context, data *RegisteredDomainResourceModel) (*common.DomainRegistration, diag.Diagnostics) {
	domainRegistration := &common.DomainRegistration{}
	diags := data.DomainRegistration.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return domainRegistration, diags
}

func tryToConvergeRegistration(ctx context.Context, data *RegisteredDomainResourceModel, diagnostics *diag.Diagnostics, r *RegisteredDomainResource, registrationID string) (string, error) {
	timeouts, diags := getTimeouts(ctx, data)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return RegistrationFailed, nil
	}

	err := utils.RetryWithTimeout(ctx, func() (error, bool) {
		domainRegistration, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), registrationID)

		if err != nil {
			return err, false
		}

		if domainRegistration.Data.State == consts.DomainStateFailed {
			diagnostics.AddError(
				fmt.Sprintf("failed to register DNSimple Domain: %s", data.Name.ValueString()),
				"domain registration failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if domainRegistration.Data.State == consts.DomainStateCancelling || domainRegistration.Data.State == consts.DomainStateCancelled {
			diagnostics.AddError(
				fmt.Sprintf("failed to register DNSimple Domain: %s", data.Name.ValueString()),
				"domain registration was cancelled, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if domainRegistration.Data.State != consts.DomainStateRegistered {
			tflog.Info(ctx, fmt.Sprintf("[RETRYING] Domain registration is not complete, current state: %s", domainRegistration.Data.State))

			return fmt.Errorf("domain registration is not complete, current state: %s. You can try to run terraform again to try and converge the domain registration", domainRegistration.Data.State), false
		}

		return nil, false
	}, timeouts.CreateDuration(), 20*time.Second)

	if diagnostics.HasError() {
		// If we have diagnostic errors, we suspended the retry loop because the domain is in a bad state, and cannot converge.
		return RegistrationFailed, nil
	}

	if err != nil {
		// If we have an error, it means the retry loop timed out, and we cannot converge during this run.
		return RegistrationConvergenceTimeout, err
	}

	return RegistrationConverged, nil
}
