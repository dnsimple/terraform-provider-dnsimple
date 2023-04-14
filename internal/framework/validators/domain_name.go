package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

var _ validator.String = DomainName{}

type DomainName struct{}

func (v DomainName) Description(ctx context.Context) string {
	return "a domain name should always be lower case"
}

// MarkdownDescription returns a markdown formatted description of the
// validator's behavior, suitable for a practitioner to understand its impact.
func (v DomainName) MarkdownDescription(ctx context.Context) string {
	return "a domain name should always be lower case"
}

// Validate runs the main validation logic of the validator, reading
// configuration data out of `req` and updating `resp` with diagnostics.
func (v DomainName) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {

	var domainName types.String
	resp.Diagnostics.Append(tfsdk.ValueAs(ctx, req.ConfigValue, &domainName)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if domainName.IsUnknown() || domainName.IsNull() {
		return
	}

	domainNameLower := strings.ToLower(domainName.ValueString())
	if domainNameLower != domainName.ValueString() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Domain name should be lower case",
			fmt.Sprintf("Domain name should be lower case, but got %s", domainName.ValueString()),
		)
		return
	}

	if utils.HasUnicodeChars(domainNameLower) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Domain name should not contain unicode characters please use punycode",
			fmt.Sprintf("Domain name should not contain unicode characters, but got %s", domainName.ValueString()),
		)
		return
	}
}
