package registered_domain_contact

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func (r *RegisteredDomainContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RegisteredDomainContactResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	registrantChangeAttributes := dnsimple.CreateRegistrantChangeInput{
		DomainId:  data.DomainId.ValueString(),
		ContactId: fmt.Sprintf("%d", data.ContactId.ValueInt64()),
	}

	if !data.ExtendedAttributes.IsNull() {
		extendedAttrs := make(map[string]string)
		diags := data.ExtendedAttributes.ElementsAs(ctx, &extendedAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		registrantChangeAttributes.ExtendedAttributes = extendedAttrs
	}

	registrantChangeResponse, err := r.config.Client.Registrar.CreateRegistrantChange(ctx, r.config.AccountID, &registrantChangeAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple registrant change",
			err.Error(),
		)
		return
	}

	if registrantChangeResponse.Data.State != consts.RegistrantChangeStateCompleted {
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelled || registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelling {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to create DNSimple registrant change current state: %s", registrantChangeResponse.Data.State),
				"",
			)
			return
		}

		// Registrant change has been created, but is not yet completed
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateNew || registrantChangeResponse.Data.State == consts.RegistrantChangeStatePending {
			convergenceState, err := tryToConvergeRegistration(ctx, data, &resp.Diagnostics, r, registrantChangeResponse.Data.Id)
			if convergenceState == RegistrantChangeFailed {
				// Response is already populated with the error we can safely return
				return
			}

			if convergenceState == RegistrantChangeConvergenceTimeout {
				// We attempted to converge on the registrant change, but the registrant change was not ready
				// user needs to run terraform again to try and converge the registrant change

				// Update the data with the current registrant change
				r.updateModelFromAPIResponse(registrantChangeResponse.Data, data)

				// Save data into Terraform state
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

				// Exit with warning to prevent the state from being tainted
				resp.Diagnostics.AddWarning(
					"failed to converge on registrant change",
					err.Error(),
				)
				return
			}

			registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(data.Id.ValueInt64()))

			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Registrant Change for: %s, %d", data.DomainId.ValueString(), data.Id.ValueInt64()),
					err.Error(),
				)
				return
			}
		}
	}

	// Registrant change was successful, we can now update the resource model
	r.updateModelFromAPIResponse(registrantChangeResponse.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
