package registered_domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func (r *RegisteredDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.RegisterDomainInput{
		RegistrantID: int(data.ContactId.ValueInt64()),
	}

	if !data.AutoRenewEnabled.IsNull() {
		domainAttributes.EnableAutoRenewal = data.AutoRenewEnabled.ValueBool()
	}

	if !data.WhoisPrivacyEnabled.IsNull() {
		domainAttributes.EnableWhoisPrivacy = data.WhoisPrivacyEnabled.ValueBool()
	}

	if !data.PremiumPrice.IsNull() {
		domainAttributes.PremiumPrice = data.PremiumPrice.ValueString()
	}

	if !data.ExtendedAttributes.IsNull() {
		extendedAttrs := make(map[string]string)
		diags := data.ExtendedAttributes.ElementsAs(ctx, &extendedAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		domainAttributes.ExtendedAttributes = extendedAttrs
	}

	lowerCaseDomainName := strings.ToLower(data.Name.ValueString())
	registerDomainResponse, err := r.config.Client.Registrar.RegisterDomain(ctx, r.config.AccountID, lowerCaseDomainName, &domainAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to register DNSimple Domain",
			err.Error(),
		)
		return
	}

	if registerDomainResponse.Data.State == consts.DomainStateHosted {
		resp.Diagnostics.AddError(
			"failed to register DNSimple Domain",
			"Domain added to DNSimple as hosted, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
		)
		return
	}

	if registerDomainResponse.Data.State != consts.DomainStateRegistered {
		timeouts, diags := getTimeouts(ctx, data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		err := utils.RetryWithTimeout(ctx, func() (error, bool) {
			domainRegistration, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, lowerCaseDomainName, strconv.Itoa(int(registerDomainResponse.Data.ID)))

			if err != nil {
				return err, false
			}

			if domainRegistration.Data.State == consts.DomainStateFailed {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", lowerCaseDomainName),
					"domain registration failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
				)
				return nil, true
			}

			if domainRegistration.Data.State == consts.DomainStateCancelling || domainRegistration.Data.State == consts.DomainStateCancelled {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", lowerCaseDomainName),
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

		if resp.Diagnostics.HasError() {
			// If we have an error, we likely suspended the retry loop, due to bad state
			return
		}

		if err != nil {
			domainRegistration, secondaryErr := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, lowerCaseDomainName, strconv.Itoa(int(registerDomainResponse.Data.ID)))

			if secondaryErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Domain Registration: %s", lowerCaseDomainName),
					secondaryErr.Error(),
				)
				return
			}

			domainResponse, secondaryErr := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, lowerCaseDomainName)

			if secondaryErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Domain: %s", lowerCaseDomainName),
					secondaryErr.Error(),
				)
				return
			}

			// Commit the partial state to the model
			r.updateModelFromAPIResponsePartialCreate(ctx, data, domainRegistration.Data, domainResponse.Data)
			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

			// Exit with warning to prevent the state from being tainted
			resp.Diagnostics.AddWarning(
				"failed to converge on domain registration",
				err.Error(),
			)
			return
		}
	}

	if !data.DNSSECEnabled.IsNull() && data.DNSSECEnabled.ValueBool() {
		diags := r.setDNSSEC(ctx, data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, lowerCaseDomainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", lowerCaseDomainName),
			err.Error(),
		)
		return
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, lowerCaseDomainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", lowerCaseDomainName),
			err.Error(),
		)
		return
	}

	diags := r.updateModelFromAPIResponse(ctx, data, registerDomainResponse.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
