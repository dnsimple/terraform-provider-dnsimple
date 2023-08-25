package registered_domain_contact

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &RegisteredDomainContactResource{}
	_ resource.ResourceWithConfigure   = &RegisteredDomainContactResource{}
	_ resource.ResourceWithImportState = &RegisteredDomainContactResource{}
)

func NewRegisteredDomainContactResource() resource.Resource {
	return &RegisteredDomainContactResource{}
}

// RegisteredDomainContactResource defines the resource implementation.
type RegisteredDomainContactResource struct {
	config *common.DnsimpleProviderConfig
}

// RegisteredDomainContactResourceModel describes the resource data model.
type RegisteredDomainContactResourceModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	AccountId           types.Int64  `tfsdk:"account_id"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	DomainId            types.String `tfsdk:"domain_id"`
	State               types.String `tfsdk:"state"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	RegistryOwnerChange types.Bool   `tfsdk:"registry_owner_change"`
	IrtLockLiftedBy     types.String `tfsdk:"irt_lock_lifted_by"`
	Timeouts            types.Object `tfsdk:"timeouts"`
}

func (r *RegisteredDomainContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registered_domain_contact"
}

func (r *RegisteredDomainContactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple registered domain contact resource",
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
				PlanModifiers: []planmodifier.String{
					RegistrantChangeState(),
				},
				Computed: true,
				Optional: true,
			},
			"extended_attributes": schema.MapAttribute{
				MarkdownDescription: "Extended attributes for the registrant change",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"registry_owner_change": schema.BoolAttribute{
				MarkdownDescription: "True if the registrant change will result in a registry owner change",
				Computed:            true,
			},
			"irt_lock_lifted_by": schema.StringAttribute{
				MarkdownDescription: "Date when the registrant change lock was lifted for the domain",
				Computed:            true,
			},
			"timeouts": schema.SingleNestedAttribute{
				MarkdownDescription: "Timeouts for operations, given as a parsable string as in `10m` or `30s`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Create timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("5m"),
						},
					},
					"update": schema.StringAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Update timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("1m"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
