package registrant_change

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &RegistrantChangeResource{}
	_ resource.ResourceWithConfigure   = &RegistrantChangeResource{}
	_ resource.ResourceWithImportState = &RegistrantChangeResource{}
)

func NewRegistrantChangeResource() resource.Resource {
	return &RegistrantChangeResource{}
}

// RegistrantChangeResource defines the resource implementation.
type RegistrantChangeResource struct {
	config *common.DnsimpleProviderConfig
}

// RegistrantChangeResourceModel describes the resource data model.
type RegistrantChangeResourceModel struct {
	Id                  types.String        `tfsdk:"id"`
	AccountId           types.Int64         `tfsdk:"account_id"`
	ContactId           types.Int64         `tfsdk:"contact_id"`
	DomainId            types.String        `tfsdk:"domain_id"`
	State               types.String        `tfsdk:"state"`
	ExtendedAttributes  []ExtendedAttribute `tfsdk:"extended_attributes"`
	RegistryOwnerChange types.Bool          `tfsdk:"registry_owner_change"`
	IrtLockLiftedBy     types.Date          `tfsdk:"irt_lock_lifted_by"`
}

type ExtendedAttribute struct {
	Name        types.String              `tfsdk:"name"`
	Description types.String              `tfsdk:"description"`
	Required    types.Bool                `tfsdk:"required"`
	Options     []ExtendedAttributeOption `tfsdk:"options"`
}

type ExtendedAttributeOption struct {
	Title       types.String `tfsdk:"title"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

func (r *RegistrantChangeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registrant_change"
}

func (r *RegistrantChangeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple registrant change resource",
		Attributes: map[string]schema.Attribute{
			"id": common.IDInt64Attribute(),
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "DNSimple Account ID to which the registrant change belongs to",
				Computed:            true,
			},
			"contact_id": schema.Int64Attribute{
				MarkdownDescription: "DNSimple contact ID for which the registrant change is being performed",
				Required:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "DNSimple domain ID for which the registrant change is being performed",
				Required:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "State of the registrant change",
				Computed:            true,
			},
			"extended_attributes": schema.ListAttribute{
				MarkdownDescription: "Extended attributes for the registrant change",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"description": types.StringType,
						"required":    types.BoolType,
						"options": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"title":       types.StringType,
									"value":       types.StringType,
									"description": types.StringType,
								},
							},
						},
					},
				},
			},
			"registry_owner_change": schema.BoolAttribute{
				MarkdownDescription: "True if the registrant change will result in a registry owner change",
				Computed:            true,
			},
			"irt_lock_lifted_by": schema.DateAttribute{
				MarkdownDescription: "Date when the registrant change lock was lifted for the domain",
				Comuted:             true,
			},
		},
	}
}
