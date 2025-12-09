package registered_domain

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func (r *RegisteredDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.RegisterDomainInput{
		RegistrantID: int(data.ContactId.ValueInt64()),
	}

	if !data.AutoRenewEnabled.IsNull() {
		domainAttributes.EnableAutoRenewal = data.AutoRenewEnabled.ValueBool()
	}

	if !data.WhoisPrivacyEnabled.IsNull() {
		domainAttributes.EnableWhoisPrivacy = data.WhoisPrivacyEnabled.ValueBool()
	}

	if !data.PremiumPrice.IsNull() {
		domainAttributes.PremiumPrice = data.PremiumPrice.ValueString()
	}

	if !data.ExtendedAttributes.IsNull() {
		extendedAttrs := make(map[string]string)
		diags := data.ExtendedAttributes.ElementsAs(ctx, &extendedAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		domainAttributes.ExtendedAttributes = extendedAttrs
	}

	registerDomainResponse, err := r.config.Client.Registrar.RegisterDomain(ctx, r.config.AccountID, data.Name.ValueString(), &domainAttributes)
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

	if registerDomainResponse.Data.State != consts.DomainStateRegistered {
		convergenceState, err := tryToConvergeRegistration(ctx, data, &resp.Diagnostics, r, strconv.Itoa(int(registerDomainResponse.Data.ID)))
		if convergenceState == RegistrationFailed {
			// Response is already populated with the error we can safely return
			return
		}

		if convergenceState == RegistrationConvergenceTimeout {
			domainRegistration, secondaryErr := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), strconv.Itoa(int(registerDomainResponse.Data.ID)))

			if secondaryErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Domain Registration: %s", data.Name.ValueString()),
					secondaryErr.Error(),
				)
				return
			}

			domainResponse, secondaryErr := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())

			if secondaryErr != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Domain: %s", data.Name.ValueString()),
					secondaryErr.Error(),
				)
				return
			}

			// Commit the partial state to the model
			r.updateModelFromAPIResponsePartialCreate(ctx, data, domainRegistration.Data, domainResponse.Data)
			// Save data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

			// Exit with warning to prevent the state from being tainted
			resp.Diagnostics.AddWarning(
				"failed to converge on domain registration",
				err.Error(),
			)
			return
		}
	}

	// Domain registration was successful, we can now proceed with the rest of the resource

	if !data.DNSSECEnabled.IsNull() && data.DNSSECEnabled.ValueBool() {
		diags := r.setDNSSEC(ctx, data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !data.TransferLockEnabled.IsNull() && data.TransferLockEnabled.ValueBool() {
		diags := r.setTransferLock(ctx, data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
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

	diags := r.updateModelFromAPIResponse(ctx, data, registerDomainResponse.Data, domainResponse.Data, dnssecResponse.Data, transferLockResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.RegistrantChange = types.ObjectNull(common.RegistrantChangeAttrType)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
