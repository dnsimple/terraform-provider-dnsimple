package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = RecordType{}

type RecordType struct{}

func (v RecordType) Description(ctx context.Context) string {
	return "record type must be specified in UPPERCASE"
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (v RecordType) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate runs the main validation logic of the validator, reading
// configuration data out of `req` and updating `resp` with diagnostics.
func (v RecordType) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var recordType types.String
	resp.Diagnostics.Append(tfsdk.ValueAs(ctx, req.ConfigValue, &recordType)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if recordType.IsUnknown() || recordType.IsNull() {
		return
	}

	recordTypeValue := recordType.ValueString()
	recordTypeUpper := strings.ToUpper(recordTypeValue)
	if recordTypeUpper != recordTypeValue {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Record type must be uppercase",
			fmt.Sprintf("Record type must be specified in UPPERCASE, but got %q. Use %q instead.", recordTypeValue, recordTypeUpper),
		)
		return
	}
}
