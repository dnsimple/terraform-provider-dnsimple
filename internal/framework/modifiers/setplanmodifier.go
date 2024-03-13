package modifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setTrimSuffix struct {
}

func SetTrimSuffixValue() planmodifier.Set {
	return setTrimSuffix{}
}

func (m setTrimSuffix) Description(context.Context) string {
	return "Trim suffix from configuration value when comparing with state"
}

func (m setTrimSuffix) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m setTrimSuffix) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.ConfigValue.IsNull() || req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		return
	}

	var planValue []string
	resp.Diagnostics.Append(req.PlanValue.ElementsAs(ctx, &planValue, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, element := range planValue {
		planValue[i] = strings.TrimSuffix(element, ".")
	}

	serializedPlanValue, diags := types.SetValueFrom(ctx, types.StringType, planValue)
	resp.Diagnostics.Append(diags...)

	resp.PlanValue = serializedPlanValue
}
