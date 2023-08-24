package registrant_change

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

func (r *RegistrantChangeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RegistrantChangeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	registrantChangeAttributes := dnsimple.CreateRegistrantChangeInput{
		DomainID:  int(data.DomainId.ValueInt64()),
		ContactID: int(data.ContactId.ValueInt64()),
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
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateFailed {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to create DNSimple registrant change: %s", registrantChangeResponse.Data.FailureReason),
				"",
			)
			return
		}

		// Registrant change has been created, but is not yet completed
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateNew {


		// Registrant change is pending
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStatePending {

		}
	}

	// Registrant change was successful, we can now update the resource model
	diags := r.updateModelFromAPIResponse(ctx, data, registrantChangeResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
