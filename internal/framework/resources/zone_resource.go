package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ZoneResource{}
	_ resource.ResourceWithConfigure   = &ZoneResource{}
	_ resource.ResourceWithImportState = &ZoneResource{}
)

func NewZoneResource() resource.Resource {
	return &ZoneResource{}
}

// ZoneResource defines the resource implementation.
type ZoneResource struct {
	config *common.DnsimpleProviderConfig
}

// ZoneResourceModel describes the resource data model.
type ZoneResourceModel struct {
	Name              types.String `tfsdk:"name"`
	AccountId         types.Int64  `tfsdk:"account_id"`
	Reverse           types.Bool   `tfsdk:"reverse"`
	Secondary         types.Bool   `tfsdk:"secondary"`
	Active            types.Bool   `tfsdk:"active"`
	LastTransferredAt types.String `tfsdk:"last_transferred_at"`
	Id                types.Int64  `tfsdk:"id"`
}

func (r *ZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (r *ZoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple zone resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_id": schema.Int64Attribute{
				Computed: true,
			},
			"reverse": schema.BoolAttribute{
				Computed: true,
			},
			"secondary": schema.BoolAttribute{
				Computed: true,
			},
			"active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"last_transferred_at": schema.StringAttribute{
				Computed: true,
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *ZoneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*common.DnsimpleProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *provider.DnsimpleProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.config = config
}

func (r *ZoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ZoneResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Zones.GetZone(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to retrieve DNSimple Zone",
			err.Error(),
		)
		return
	}

	if !(data.Active.IsUnknown() || data.Active.IsNull()) && data.Active.ValueBool() != response.Data.Active {
		zone, diags := r.setActiveState(ctx, data)

		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		r.updateModelFromAPIResponse(zone, data)

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ZoneResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Zones.GetZone(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Zone: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		configData *ZoneResourceModel
		planData   *ZoneResourceModel
		stateData  *ZoneResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !(planData.Active.IsUnknown() || planData.Active.IsNull()) && planData.Active.ValueBool() != stateData.Active.ValueBool() {
		zone, diags := r.setActiveState(ctx, planData)

		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}

		r.updateModelFromAPIResponse(zone, planData)

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

		return
	}
}

func (r *ZoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ZoneResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, fmt.Sprintf("Removing DNSimple Zone from Terraform state only: %s, %s", data.Name, data.Id))
}

func (r *ZoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	response, err := r.config.Client.Zones.GetZone(ctx, r.config.AccountID, req.ID)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to find DNSimple Zone ID: %s", req.ID),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), response.Data.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), response.Data.Name)...)
}

func (r *ZoneResource) updateModelFromAPIResponse(zone *dnsimple.Zone, data *ZoneResourceModel) {
	data.Id = types.Int64Value(zone.ID)
	data.Name = types.StringValue(zone.Name)
	data.AccountId = types.Int64Value(zone.AccountID)
	data.Reverse = types.BoolValue(zone.Reverse)
	data.Secondary = types.BoolValue(zone.Secondary)
	data.Active = types.BoolValue(zone.Active)
	data.LastTransferredAt = types.StringValue(zone.LastTransferredAt)
}

func (r *ZoneResource) setActiveState(ctx context.Context, data *ZoneResourceModel) (*dnsimple.Zone, diag.Diagnostics) {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting active to %t", data.Active.ValueBool()))

	if data.Active.ValueBool() {
		zoneResponse, err := r.config.Client.Zones.ActivateZoneDns(ctx, r.config.AccountID, data.Name.ValueString())
		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to activate DNSimple Zone: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return zoneResponse.Data, diagnostics
	}

	zoneResponse, err := r.config.Client.Zones.DeactivateZoneDns(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to deactivate DNSimple Zone: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return zoneResponse.Data, diagnostics
}
