package registered_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
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

// RegistrantChangeResourceModel describes the resource data model.
type RegistrantChangeResourceModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	AccountId           types.Int64  `tfsdk:"account_id"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	DomainId            types.String `tfsdk:"domain_id"`
	State               types.String `tfsdk:"state"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	RegistryOwnerChange types.Bool   `tfsdk:"registry_owner_change"`
	IrtLockLiftedBy     types.String `tfsdk:"irt_lock_lifted_by"`
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
	TransferLockEnabled types.Bool   `tfsdk:"transfer_lock_enabled"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	PremiumPrice        types.String `tfsdk:"premium_price"`
	DomainRegistration  types.Object `tfsdk:"domain_registration"`
	RegistrantChange    types.Object `tfsdk:"registrant_change"`
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
			"transfer_lock_enabled": schema.BoolAttribute{
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
			"registrant_change": schema.SingleNestedAttribute{
				Description: "The registrant change details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"id": common.IDInt64Attribute(),
					"account_id": schema.Int64Attribute{
						MarkdownDescription: "DNSimple Account ID to which the registrant change belongs to",
						Computed:            true,
					},
					"contact_id": schema.Int64Attribute{
						MarkdownDescription: "DNSimple contact ID for which the registrant change is being performed",
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.RequiresReplace(),
						},
					},
					"domain_id": schema.StringAttribute{
						MarkdownDescription: "DNSimple domain ID for which the registrant change is being performed",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.RequiresReplace(),
						},
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
						Computed:            true,
						PlanModifiers: []planmodifier.Map{
							mapplanmodifier.RequiresReplaceIfConfigured(),
						},
					},
					"registry_owner_change": schema.BoolAttribute{
						MarkdownDescription: "True if the registrant change will result in a registry owner change",
						Computed:            true,
					},
					"irt_lock_lifted_by": schema.StringAttribute{
						MarkdownDescription: "Date when the registrant change lock was lifted for the domain",
						Computed:            true,
					},
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
