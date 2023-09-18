package registered_domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

func (r *RegisteredDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	domainName := parts[0]

	usingDomainAndRegistrationID := len(parts) == 2
	if usingDomainAndRegistrationID {
		domainRegistrationID := parts[1]

		domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, domainName, domainRegistrationID)

		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to find DNSimple Domain Registration ID: %s", domainRegistrationID),
				err.Error(),
			)
			return
		}
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_registration").AtName("id"), domainRegistrationResponse.Data.ID)...)
	} else {
		resp.Private.SetKey(ctx, "skip_domain_registration", []byte(`true`))
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_registration"), basetypes.NewObjectNull(common.DomainRegistrationAttrType))...)
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, domainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected error when trying to find DNSimple Domain ID: %s", domainName),
			err.Error(),
		)
		return
	}

	if domainResponse.Data.State != consts.DomainStateRegistered && !usingDomainAndRegistrationID {
		resp.Diagnostics.AddError(
			fmt.Sprintf("domain %s is not registered", domainName),
			"domain must be registered before it can be imported",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainResponse.Data.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), domainResponse.Data.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("contact_id"), domainResponse.Data.RegistrantID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("registrant_change"), types.ObjectNull(common.RegistrantChangeAttrType))...)
}
