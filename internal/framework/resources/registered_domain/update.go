package registered_domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func (r *RegisteredDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		configData *RegisteredDomainResourceModel
		planData   *RegisteredDomainResourceModel
		stateData  *RegisteredDomainResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.ContactId.ValueInt64() != stateData.ContactId.ValueInt64() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("contact_id change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
			"contact_id change not supported by the DNSimple API",
		)
		return
	}

	if !planData.ExtendedAttributes.Equal(stateData.ExtendedAttributes) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("extended_attributes change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
			"extended_attributes change not supported by the DNSimple API",
		)
		return
	}

	domainRegistration, diags := getDomainRegistration(ctx, stateData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if domainRegistration.Id.IsNull() {
		resp.Diagnostics.AddError(
			"failed to read domain registration from state",
			"domain registration is nil. This should not happen. Please remove the resource from the state and try importing.",
		)
		return
	}

	response, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, configData.Name.ValueString(), strconv.Itoa(int(domainRegistration.Id.ValueInt64())))

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

	// Check if the domain is in registered state
	if response.Data.State != consts.DomainStateRegistered {
		timeouts, diags := getTimeouts(ctx, planData)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		err := utils.RetryWithTimeout(ctx, func() (error, bool) {

			domainRegistration, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, configData.Name.ValueString(), strconv.Itoa(int(response.Data.ID)))

			if err != nil {
				return err, false
			}

			if domainRegistration.Data.State == consts.DomainStateFailed {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", configData.Name.ValueString()),
					"domain registration failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
				)
				return nil, true
			}

			if domainRegistration.Data.State == consts.DomainStateCancelling || domainRegistration.Data.State == consts.DomainStateCancelled {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", configData.Name.ValueString()),
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
			domainRegistration, secondaryErr := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, configData.Name.ValueString(), strconv.Itoa(int(response.Data.ID)))

			if secondaryErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Domain Registration: %s", configData.Name.ValueString()),
					secondaryErr.Error(),
				)
				return
			}

			if domainRegistration.Data.State == consts.DomainStateFailed {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", configData.Name.ValueString()),
					"domain registration failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
				)
				return
			}

			if domainRegistration.Data.State == consts.DomainStateCancelling || domainRegistration.Data.State == consts.DomainStateCancelled {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to register DNSimple Domain: %s", configData.Name.ValueString()),
					"domain registration was cancelled, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
				)
				return
			}

			// We attempted to converge on the domain registration, but the domain registration was not ready
			// user needs to run terraform again to try and converge the domain registration
			resp.Diagnostics.AddError(
				"failed to converge on domain registration",
				err.Error(),
			)
			return
		}
	}

	if planData.AutoRenewEnabled.ValueBool() != stateData.AutoRenewEnabled.ValueBool() {

		diags := r.setAutoRenewal(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if planData.WhoisPrivacyEnabled.ValueBool() != stateData.WhoisPrivacyEnabled.ValueBool() {

		diags := r.setWhoisPrivacy(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if planData.DNSSECEnabled.ValueBool() != stateData.DNSSECEnabled.ValueBool() {

		diags := r.setDNSSEC(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, planData.Name.ValueString(), domainRegistrationId)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", planData.Name.ValueString(), domainRegistration.Id.ValueInt64()),
			err.Error(),
		)
		return
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, planData.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", planData.Name.ValueString()),
			err.Error(),
		)
		return
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, planData.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", planData.Name.ValueString()),
			err.Error(),
		)
		return
	}

	diags = r.updateModelFromAPIResponse(ctx, planData, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated planData into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
