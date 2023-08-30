package registered_domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

	var domainRegistrationResponse *dnsimple.DomainRegistrationResponse
	var err error
	if !domainRegistration.Id.IsNull() {
		registerDomainResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, configData.Name.ValueString(), strconv.Itoa(int(domainRegistration.Id.ValueInt64())))

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
		if registerDomainResponse.Data.State != consts.DomainStateRegistered {
			convergenceState, err := tryToConvergeRegistration(ctx, planData, &resp.Diagnostics, r, strconv.Itoa(int(registerDomainResponse.Data.ID)))
			if convergenceState == RegistrationFailed {
				// Response is already populated with the error we can safely return
				return
			}

			if convergenceState == RegistrationConvergenceTimeout {
				// We attempted to converge on the domain registration, but the domain registration was not ready
				// user needs to run terraform again to try and converge the domain registration
				resp.Diagnostics.AddError(
					"failed to converge on domain registration",
					err.Error(),
				)
				return
			}
		}

		domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
		domainRegistrationResponse, err = r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, planData.Name.ValueString(), domainRegistrationId)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", planData.Name.ValueString(), domainRegistration.Id.ValueInt64()),
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

	if domainRegistrationResponse == nil {
		diags = r.updateModelFromAPIResponse(ctx, planData, nil, domainResponse.Data, dnssecResponse.Data)
	} else {
		diags = r.updateModelFromAPIResponse(ctx, planData, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data)
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated planData into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
