package registered_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

type domainRegistrationState struct {
}

// DomainRegistrationState return a object plan modifier that sets the specified value if the planned value is Null.
func DomainRegistrationState() planmodifier.Object {
	return domainRegistrationState{}
}

func (m domainRegistrationState) Description(context.Context) string {
	return "If the domain registration state is not registered, set it to registered. Unless the state is a failing state"
}

func (m domainRegistrationState) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m domainRegistrationState) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.ConfigValue.IsNull() || req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	// Check if the domain registrstion status is as expected
	domainRegistration := &common.DomainRegistration{}
	if diags := req.PlanValue.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// If the domain registration state is a failing state, do not attempt to update it
	if domainRegistration.State.ValueString() == consts.DomainStateFailed {
		return
	}

	// If the domain registration state is a cancelled state, do not attempt to update it
	if domainRegistration.State.ValueString() == consts.DomainStateCancelling || domainRegistration.State.ValueString() == consts.DomainStateCancelled {
		return
	}

	// If the domain registration state is not registered, set it to registered
	// this will trigger a plan change and result in an update so we can attempt to sync
	if domainRegistration.State.ValueString() == consts.DomainStateRegistered {
		return
        }
	domainRegistration.State = types.StringValue(consts.DomainStateRegistered)

	obj, diags := types.ObjectValueFrom(ctx, common.DomainRegistrationAttrType, domainRegistration)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.PlanValue = obj
}
