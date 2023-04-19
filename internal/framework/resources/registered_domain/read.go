package registered_domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *RegisteredDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainRegistration, diags := getDomainRegistration(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), domainRegistrationId)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", data.Name.ValueString(), domainRegistration.Id.ValueInt64()),
			err.Error(),
		)
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", data.Name.ValueString()),
			err.Error(),
		)
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", data.Name.ValueString()),
			err.Error(),
		)
	}

	diags = r.updateModelFromAPIResponse(ctx, data, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
