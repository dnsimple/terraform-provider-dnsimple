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

func (r *RegisteredDomainContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		configData *RegisteredDomainContactResourceModel
		planData   *RegisteredDomainContactResourceModel
		stateData  *RegisteredDomainContactResourceModel
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

	var registrantChangeResponse *dnsimple.RegistrantChangeResponse
	var err error
	registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(planData.Id.ValueInt64()))

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to read DNSimple registrant change",
			err.Error(),
		)
		return
	}

	// Check if the registrant change is completed
	if registrantChangeResponse.Data.State != consts.RegistrantChangeStateCompleted {
		convergenceState, err := tryToConvergeRegistration(ctx, planData, &resp.Diagnostics, r, registrantChangeResponse.Data.Id)
		if convergenceState == RegistrantChangeFailed {
			// Response is already populated with the error we can safely return
			return
		}

		if convergenceState == RegistrantChangeConvergenceTimeout {
			// We attempted to converge on the registrant change, but the registrant change was not ready
			// user needs to run terraform again to try and converge the registrant change

			// Exit with warning to prevent the state from being tainted
			resp.Diagnostics.AddError(
				"failed to converge on registrant change",
				err.Error(),
			)
			return
		}

		registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(planData.Id.ValueInt64()))

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Registrant Change for: %s, %d", planData.DomainId.ValueString(), planData.Id.ValueInt64()),
				err.Error(),
			)
			return
		}
	}

	r.updateModelFromAPIResponse(registrantChangeResponse.Data, planData)

	// Save updated planData into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}
