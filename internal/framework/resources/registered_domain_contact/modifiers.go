package registered_domain_contact

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
)

type registrantChangeState struct {
}

// RegistrantChangeState return a object plan modifier that sets the specified value if the planned value is Null.
func RegistrantChangeState() planmodifier.String {
	return registrantChangeState{}
}

func (m registrantChangeState) Description(context.Context) string {
	return "If the registrant change state is not completed, set it to completed. Unless the state is a cancelled or cancelling state"
}

func (m registrantChangeState) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m registrantChangeState) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.ConfigValue.IsNull() || req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	state := req.PlanValue.ValueString()

	// If the registrant change state is a cancelled state, do not attempt to update it
	if state == consts.RegistrantChangeStateCancelling || state == consts.RegistrantChangeStateCancelled {
		return
	}

	// If the registrant change state is not completed, set it to completed
	// this will trigger a plan change and result in an update so we can attempt to sync
	if state == consts.RegistrantChangeStateCompleted {
		return
	}

	resp.PlanValue = types.StringValue(consts.RegistrantChangeStateCompleted)
}
