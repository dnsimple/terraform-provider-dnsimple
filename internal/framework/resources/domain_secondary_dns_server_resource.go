package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DomainSecondaryServerResource{}
	_ resource.ResourceWithConfigure   = &DomainSecondaryServerResource{}
)

func NewDomainSecondaryServerResource() resource.Resource {
	return &DomainSecondaryServerResource{}
}

// DomainSecondaryServerResource defines the resource implementation.
type DomainSecondaryServerResource struct {
	config *common.DnsimpleProviderConfig
}

// DomainSecondaryServerResourceModel describes the resource data model.
type DomainSecondaryServerResourceModel struct {
	Name         types.String `tfsdk:"name"`
	IPAddress    types.String `tfsdk:"ip_address"`
	Port         types.Int64  `tfsdk:"port"`
	ID           types.Int64  `tfsdk:"id"`
	Zones        types.Set    `tfsdk:"zones"`
}

func (r *DomainSecondaryServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_secondary_server"
}

func (r *DomainSecondaryServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain secondary server resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"port": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"zones": schema.SetAttribute{
				Required: true,
				ElementType: types.StringType,
				// For now, we do in-place replace, but we can technicaly link one-by-one
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *DomainSecondaryServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainSecondaryServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DomainSecondaryServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverAttributes := dnsimple.SecondaryServer{
		Name: data.Name.ValueString(),
		IP: data.IPAddress.ValueString(),
		Port: uint64(data.Port.ValueInt64()),
	}

	response, err := r.config.Client.SecondaryDNS.CreatePrimaryServer(ctx, r.config.AccountID, serverAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple secondary DNS primary server",
			err.Error(),
		)
		return
	}

	for _, zoneValue := range data.Zones.Elements() {
		tfv, err := zoneValue.ToTerraformValue(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to get zone name as string when creating DNSimple secondary DNS server",
				err.Error(),
			)
			return
		}
		var zone string
		if err := tfv.As(&zone); err != nil {
			resp.Diagnostics.AddError(
				"failed to convert zone name to string when creating DNSimple secondary DNS server",
				err.Error(),
			)
			return
		}

		if _, err := r.config.Client.SecondaryDNS.LinkPrimaryServerToSecondaryZone(ctx, r.config.AccountID, string(data.ID.ValueInt64()), zone); err != nil {
			var errorResponse *dnsimple.ErrorResponse
			if errors.As(err, &errorResponse) {
				resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
				return
			}

			resp.Diagnostics.AddError(
				"failed to link zone to DNSimple secondary DNS server",
				err.Error(),
			)
			return
		}
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainSecondaryServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DomainSecondaryServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.SecondaryDNS.GetPrimaryServer(ctx, r.config.AccountID, string(data.ID.ValueInt64()))

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainSecondaryServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op
}

func (r *DomainSecondaryServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DomainSecondaryServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple secondary DNS primary server: %s, %s", data.Name, data.ID))

	_, err := r.config.Client.SecondaryDNS.DeletePrimaryServer(ctx, r.config.AccountID, string(data.ID.ValueInt64()))

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to delete DNSimple secondary DNS primary server: %s", data.Name.ValueString()),
			err.Error(),
		)
		return
	}
}

func (r *DomainSecondaryServerResource) updateModelFromAPIResponse(server *dnsimple.SecondaryServer, data *DomainSecondaryServerResourceModel) {
	data.ID = types.Int64Value(server.ID)
	data.Name = types.StringValue(server.Name)
	data.IPAddress = types.StringValue(server.IP)
	data.Port = types.Int64Value(int64(server.Port))
	zones := make([]attr.Value, 0, len(server.LinkedSecondaryZones))
	for _, z := range server.LinkedSecondaryZones {
		zones = append(zones, basetypes.NewStringValue(z))
	}
	data.Zones = basetypes.NewSetValueMust(types.StringType, zones)
}
