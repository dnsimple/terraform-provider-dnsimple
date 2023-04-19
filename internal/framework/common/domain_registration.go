package common

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DomainRegistration struct {
	Period types.Int64  `tfsdk:"period"`
	State  types.String `tfsdk:"state"`
	Id     types.Int64  `tfsdk:"id"`
}

var DomainRegistrationAttrType = map[string]attr.Type{
	"period": types.Int64Type,
	"state":  types.StringType,
	"id":     types.Int64Type,
}
