package registrant_change

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *RegistrantChangeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RegistrantChangeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	registrantChange, diags := getRegistrantChange(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var registrantChangeResponse *dnsimple.RegistrantChangeResponse
	var err error
	if !registrantChange.Id.IsNull() {
		registrantChangeId := strconv.Itoa(int(registrantChange.Id.ValueInt64()))
		registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, registrantChangeId)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Registrant Change: %s, %d", data.Name.ValueString(), registrantChange.Id.ValueInt64()),
				err.Error(),
			)
			return
		}
	}

	if registrantChangeResponse == nil {
		diags = r.updateModelFromAPIResponse(ctx, data, registrantChangeResponse.Data)
	}

	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
