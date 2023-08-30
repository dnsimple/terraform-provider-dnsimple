package registered_domain_contact

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *RegisteredDomainContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", "The ID must be an integer")
		return
	}

	registrantChangeResponse, err := r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read DNSimple Registrant Change",
			err.Error(),
		)
		return
	}

	domainId := strconv.Itoa(registrantChangeResponse.Data.DomainId)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(id))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), domainId)...)
}
