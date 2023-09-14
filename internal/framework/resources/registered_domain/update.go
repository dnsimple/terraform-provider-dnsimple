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

	if planData.ContactId.ValueInt64() != stateData.ContactId.ValueInt64() {
		if stateData.State.ValueString() != consts.DomainStateRegistered {
			resp.Diagnostics.AddError(
				fmt.Sprintf("contact_id change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
				"contact_id change not supported for domains that are not in registered state",
			)
			return
		}
	}

	if !planData.ExtendedAttributes.Equal(stateData.ExtendedAttributes) {
		if stateData.State.ValueString() != consts.DomainStateRegistered {
			resp.Diagnostics.AddError(
				fmt.Sprintf("extended_attributes change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
				"extended_attributes change not supported for domains that are not in registered state",
			)
			return
		}
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

	registrantChange, diags := getRegistrantChange(ctx, planData)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var registrantChangeResponse *dnsimple.RegistrantChangeResponse
	if planData.ContactId.ValueInt64() != stateData.ContactId.ValueInt64() {
		if !registrantChange.Id.IsNull() {
			convergenceState, err := tryToConvergeRegistrantChange(ctx, planData, &resp.Diagnostics, r, int(registrantChange.Id.ValueInt64()))
			if convergenceState == RegistrantChangeFailed {
				// Response is already populated with the error we can safely return
				return
			}

			if convergenceState == RegistrantChangeConvergenceTimeout {
				// We attempted to converge on the registrant change, but the registrant change was not ready
				// user needs to run terraform again to try and converge the registrant change

				// Update the data with the current registrant change
				registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
				planData.RegistrantChange = registrantChangeObject

				// Save data into Terraform state
				resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

				// Exit with warning to prevent the state from being tainted
				resp.Diagnostics.AddWarning(
					"failed to converge on registrant change",
					err.Error(),
				)
				return
			}

		} else {
			// Create a new registrant change and handle any errors
			createRegistrantChange(ctx, planData, r, resp)
		}
	} else if !registrantChange.Id.IsNull() {
		registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(registrantChange.Id.ValueInt64()))

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Registrant Change Id: %d", registrantChange.Id.ValueInt64()),
				err.Error(),
			)
			return
		}

		if registrantChangeResponse.Data.State != consts.RegistrantChangeStateCompleted {
			convergenceState, err := tryToConvergeRegistrantChange(ctx, planData, &resp.Diagnostics, r, int(registrantChange.Id.ValueInt64()))
			if convergenceState == RegistrantChangeFailed {
				// Response is already populated with the error we can safely return
				return
			}

			if convergenceState == RegistrantChangeConvergenceTimeout {
				// We attempted to converge on the registrant change, but the registrant change was not ready
				// user needs to run terraform again to try and converge the registrant change

				// Update the data with the current registrant change
				registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
				planData.RegistrantChange = registrantChangeObject

				// Save data into Terraform state
				resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

				// Exit with warning to prevent the state from being tainted
				resp.Diagnostics.AddError(
					"failed to converge on registrant change",
					err.Error(),
				)
				return
			}
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

	if planData.TransferLockEnabled.ValueBool() != stateData.TransferLockEnabled.ValueBool() {

		diags := r.setTransferLock(ctx, planData)
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

	transferLockResponse, err := r.config.Client.Registrar.GetDomainTransferLock(ctx, r.config.AccountID, planData.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain transfer lock status: %s", planData.Name.ValueString()),
			err.Error(),
		)
		return
	}

	if domainRegistrationResponse == nil {
		diags = r.updateModelFromAPIResponse(ctx, planData, nil, domainResponse.Data, dnssecResponse.Data, transferLockResponse.Data)
	} else {
		diags = r.updateModelFromAPIResponse(ctx, planData, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data, transferLockResponse.Data)
	}
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated planData into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
