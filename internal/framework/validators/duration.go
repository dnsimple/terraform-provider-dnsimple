package validators

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = Duration{}

type Duration struct{}

func (v Duration) Description(ctx context.Context) string {
	return "a duration given as a parsable string as in 60m or 2h"
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (v Duration) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading
// configuration data out of `req` and updating `resp` with diagnostics.
func (v Duration) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {

	var duration types.String
	resp.Diagnostics.Append(tfsdk.ValueAs(ctx, req.ConfigValue, &duration)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if duration.IsUnknown() || duration.IsNull() {
		return
	}

	_, err := time.ParseDuration(duration.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid duration",
			err.Error(),
		)
		return
	}
}
