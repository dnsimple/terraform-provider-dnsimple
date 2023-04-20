package registered_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &RegisteredDomainResource{}
	_ resource.ResourceWithConfigure   = &RegisteredDomainResource{}
	_ resource.ResourceWithImportState = &RegisteredDomainResource{}
)

func NewRegisteredDomainResource() resource.Resource {
	return &RegisteredDomainResource{}
}

// RegisteredDomainResource defines the resource implementation.
type RegisteredDomainResource struct {
	config *common.DnsimpleProviderConfig
}

// DomainResourceModel describes the resource data model.
type RegisteredDomainResourceModel struct {
	Name                types.String `tfsdk:"name"`
	AccountId           types.Int64  `tfsdk:"account_id"`
	UnicodeName         types.String `tfsdk:"unicode_name"`
	State               types.String `tfsdk:"state"`
	AutoRenewEnabled    types.Bool   `tfsdk:"auto_renew_enabled"`
	WhoisPrivacyEnabled types.Bool   `tfsdk:"whois_privacy_enabled"`
	DNSSECEnabled       types.Bool   `tfsdk:"dnssec_enabled"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	PremiumPrice        types.String `tfsdk:"premium_price"`
	DomainRegistration  types.Object `tfsdk:"domain_registration"`
	Timeouts            types.Object `tfsdk:"timeouts"`
	Id                  types.Int64  `tfsdk:"id"`
}

func (r *RegisteredDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registered_domain"
}

func (r *RegisteredDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.DomainName{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_id": schema.Int64Attribute{
				Computed: true,
			},
			"unicode_name": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"auto_renew_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"whois_privacy_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"dnssec_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"contact_id": schema.Int64Attribute{
				Required: true,
			},
			"expires_at": schema.StringAttribute{
				Computed: true,
			},
			"extended_attributes": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"premium_price": schema.StringAttribute{
				Optional: true,
			},
			"domain_registration": schema.SingleNestedAttribute{
				Description: "The domain registration details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"period": schema.Int64Attribute{
						Computed: true,
					},
					"state": schema.StringAttribute{
						Computed: true,
						Optional: true,
					},
					"id": common.IDInt64Attribute(),
				},
				PlanModifiers: []planmodifier.Object{
					DomainRegistrationState(),
				},
			},
			"timeouts": schema.SingleNestedAttribute{
				MarkdownDescription: "Timeouts for operations, given as a parsable string as in `10m` or `30s`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						Optional:    true,
						Description: "Create timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("10m"),
						},
					},
					"update": schema.StringAttribute{
						Optional:    true,
						Description: "Update timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("30s"),
						},
					},
					"delete": schema.StringAttribute{
						Optional:    true,
						Description: "Delete timeout (currently unused).",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("30s"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"id": common.IDInt64Attribute(),
		},
	}
}
