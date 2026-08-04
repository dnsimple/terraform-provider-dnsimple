package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DomainResource{}
	_ resource.ResourceWithConfigure   = &DomainResource{}
	_ resource.ResourceWithImportState = &DomainResource{}
)

func NewDomainResource() resource.Resource {
	return &DomainResource{}
}

// DomainResource defines the resource implementation.
type DomainResource struct {
	config *common.DnsimpleProviderConfig
}

// DomainResourceModel describes the resource data model.
type DomainResourceModel struct {
	Name          types.String `tfsdk:"name"`
	PreventDelete types.Bool   `tfsdk:"prevent_delete"`
	AccountId     types.Int64  `tfsdk:"account_id"`
	RegistrantId  types.Int64  `tfsdk:"registrant_id"`
	UnicodeName   types.String `tfsdk:"unicode_name"`
	State         types.String `tfsdk:"state"`
	AutoRenew     types.Bool   `tfsdk:"auto_renew"`
	PrivateWhois  types.Bool   `tfsdk:"private_whois"`
	Trustee       types.Bool   `tfsdk:"trustee"`
	Id            types.Int64  `tfsdk:"id"`
}

func (r *DomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

func (r *DomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			/*
			 * Deleting a domain removes its registration, which is not recoverable in
			 * the usual sense. This flag lets practitioners opt into a guard against
			 * that. It defaults to false so the resource destroys like any other
			 * Terraform resource unless protection is asked for explicitly.
			 */
			"prevent_delete": schema.BoolAttribute{
				MarkdownDescription: "Whether to block `terraform destroy` from deleting the domain registration. Defaults to `false`. When set to `true`, destroying this resource fails until it is set back to `false` and applied.",
				Optional:            true,
				Default:             booldefault.StaticBool(false),
				Computed:            true,
			},
			"account_id": schema.Int64Attribute{
				Computed: true,
			},
			"registrant_id": schema.Int64Attribute{
				Computed: true,
			},
			"unicode_name": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"auto_renew": schema.BoolAttribute{
				Computed: true,
			},
			"private_whois": schema.BoolAttribute{
				Computed: true,
			},
			"trustee": schema.BoolAttribute{
				Computed: true,
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *DomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*common.DnsimpleProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *common.DnsimpleProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.config = config
}

func (r *DomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.Domain{
		Name: data.Name.ValueString(),
	}

	response, err := r.config.Client.Domains.CreateDomain(ctx, r.config.AccountID, domainAttributes)
	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple Domain",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to read DNSimple Domain",
			fmt.Sprintf("Unable to read domain '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var preventDeleteFlag types.Bool

	/*
	 * Save updated data into Terraform state.
	 *
	 * Since the single argument that can be updated is `prevent_delete` and that
	 * value is only meant to be stored in the state and acted on when a domain might
	 * be deleted, we only need to ensure it is stored when changed.
	 */
	var priorPreventDeleteFlag types.Bool

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("prevent_delete"), &preventDeleteFlag)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("prevent_delete"), &priorPreventDeleteFlag)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Only warn on an actual true -> false transition. Warning on any planned false
	// would also fire for state written before this attribute existed, where the plan is
	// null -> false and no protection was ever in place to lose.
	if priorPreventDeleteFlag.ValueBool() && !preventDeleteFlag.ValueBool() {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("prevent_delete"),
			"disabling domain deletion protection endangers domain registration.",
			"Domain registration is lost when deleting domain resources. Only disable deletion protection if you fully understand the consequences of doing so.",
		)
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prevent_delete"), preventDeleteFlag)...)
}

func (r *DomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// ValueBool() reports false for a null or unknown value, which matches the schema
	// default, so state written before this attribute existed destroys as it always did.
	if data.PreventDelete.ValueBool() {
		resp.Diagnostics.AddError(
			"failed to delete DNSimple Domain",
			"Domain deletion protection enabled.",
		)
		resp.Diagnostics.AddWarning(
			"domain registration is lost when deleting domain resources.",
			fmt.Sprintf("Disabling deletion protection and destroying this resource will DELETE YOUR DOMAIN REGISTRATION with DNSimple for the domain %s. Note this also blocks the destroy half of a replacement, so changing 'name' fails until protection is disabled.", data.Name.ValueString()),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple Domain: %s, %s", data.Name, data.Id))

	_, err := r.config.Client.Domains.DeleteDomain(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete DNSimple Domain",
			fmt.Sprintf("Unable to delete domain '%s': %s", data.Name.ValueString(), err.Error()),
		)
		return
	}
}

func (r *DomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to import DNSimple Domain",
			fmt.Sprintf("Unable to find domain '%s': %s", req.ID, err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.Data.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), response.Data.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prevent_delete"), types.BoolValue(false))...)
}

func (r *DomainResource) updateModelFromAPIResponse(domain *dnsimple.Domain, data *DomainResourceModel) {
	// prevent_delete is stored only in state and has no API counterpart, so state
	// written before it existed carries no value for it. Settle it to the schema default
	// on refresh, otherwise every previously managed domain shows a null -> false diff
	// on the first plan after upgrading.
	if data.PreventDelete.IsNull() || data.PreventDelete.IsUnknown() {
		data.PreventDelete = types.BoolValue(false)
	}

	data.Id = types.Int64Value(domain.ID)
	data.Name = types.StringValue(domain.Name)
	data.AccountId = types.Int64Value(domain.AccountID)
	data.RegistrantId = types.Int64Value(domain.RegistrantID)
	data.UnicodeName = types.StringValue(domain.UnicodeName)
	data.State = types.StringValue(domain.State)
	data.AutoRenew = types.BoolValue(domain.AutoRenew)
	data.PrivateWhois = types.BoolValue(domain.PrivateWhois)
	data.Trustee = types.BoolValue(domain.Trustee)
}
