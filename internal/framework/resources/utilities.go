package resources

import (
	"fmt"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func attributeErrorsToDiagnostics(err *dnsimple.ErrorResponse) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	diagnostics.AddError(
		"API returned an error",
		err.Message,
	)

	for field, errors := range err.AttributeErrors {
		terraformField := translateFieldFromAPIToTerraform(field)

		diagnostics.AddAttributeError(
			path.Root(terraformField),
			fmt.Sprintf("API returned a Validation Error for: %s", terraformField),
			strings.Join(errors, ", "),
		)
	}

	return diagnostics
}

func (r *DomainResource) updateModelFromAPIResponse(domain *dnsimple.Domain, data *DomainResourceModel) {
	data.Id = types.Int64Value(domain.ID)
	data.Name = types.StringValue(domain.Name)
	data.AccountId = types.Int64Value(domain.AccountID)
	data.RegistrantId = types.Int64Value(domain.RegistrantID)
	data.UnicodeName = types.StringValue(domain.UnicodeName)
	data.State = types.StringValue(domain.State)
	data.AutoRenew = types.BoolValue(domain.AutoRenew)
	data.PrivateWhois = types.BoolValue(domain.PrivateWhois)
}

func translateFieldFromAPIToTerraform(field string) string {
	switch field {
	case "record_type":
		return "type"
	case "content":
		return "value"
	default:
		return field
	}
}
