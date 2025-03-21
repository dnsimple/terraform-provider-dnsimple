package registered_domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dnsimple/dnsimple-go/v4/dnsimple"
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

	var domainRegistrationResponse *dnsimple.DomainRegistrationResponse
	var err error
	if !domainRegistration.Id.IsNull() {
		domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
		domainRegistrationResponse, err = r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), domainRegistrationId)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", data.Name.ValueString(), domainRegistration.Id.ValueInt64()),
				err.Error(),
			)
			return
		}
	}

	registrantChange, diags := getRegistrantChange(ctx, data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	var registrantChangeResponse *dnsimple.RegistrantChangeResponse
	if !registrantChange.Id.IsNull() {
		registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(registrantChange.Id.ValueInt64()))

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to read DNSimple Registrant Change Id: %d", registrantChange.Id.ValueInt64()),
				err.Error(),
			)
			return
		}

		registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		data.RegistrantChange = registrantChangeObject
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	transferLockResponse, err := r.config.Client.Registrar.GetDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain transfer lock status: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	if domainRegistrationResponse == nil {
		diags = r.updateModelFromAPIResponse(ctx, data, nil, domainResponse.Data, dnssecResponse.Data, transferLockResponse.Data)
	} else {
		diags = r.updateModelFromAPIResponse(ctx, data, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data, transferLockResponse.Data)
	}

	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
