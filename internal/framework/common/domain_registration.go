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

type RegistrantChange struct {
	Id                  types.Int64  `tfsdk:"id"`
	AccountId           types.Int64  `tfsdk:"account_id"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	DomainId            types.String `tfsdk:"domain_id"`
	State               types.String `tfsdk:"state"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	RegistryOwnerChange types.Bool   `tfsdk:"registry_owner_change"`
	IrtLockLiftedBy     types.String `tfsdk:"irt_lock_lifted_by"`
}

var RegistrantChangeAttrType = map[string]attr.Type{
	"id":         types.Int64Type,
	"account_id": types.Int64Type,
	"contact_id": types.Int64Type,
	"domain_id":  types.StringType,
	"state":      types.StringType,
	"extended_attributes": types.MapType{
		ElemType: types.StringType,
	},
	"registry_owner_change": types.BoolType,
	"irt_lock_lifted_by":    types.StringType,
}
