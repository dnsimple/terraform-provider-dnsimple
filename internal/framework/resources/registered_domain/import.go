package registered_domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *RegisteredDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("resource import invalid ID", fmt.Sprintf("wrong format of import ID (%s), use: '<domain-name>_<domain-registration-id>'", req.ID))
		return
	}
	domainName := parts[0]
	domainRegistrationID := parts[1]

	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, domainName, domainRegistrationID)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to find DNSimple Domain Registration ID: %s", domainRegistrationID),
			err.Error(),
		)
		return
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, domainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected error when trying to find DNSimple Domain ID: %s", domainName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainResponse.Data.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), domainResponse.Data.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_registration").AtName("id"), domainRegistrationResponse.Data.ID)...)
}
