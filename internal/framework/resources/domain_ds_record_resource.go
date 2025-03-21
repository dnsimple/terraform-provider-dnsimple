package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/v4/dnsimple"
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
	_ resource.Resource                = &DsRecordResource{}
	_ resource.ResourceWithConfigure   = &DsRecordResource{}
	_ resource.ResourceWithImportState = &DsRecordResource{}
)

func NewDsRecordResource() resource.Resource {
	return &DsRecordResource{}
}

// DsRecordResource defines the resource implementation.
type DsRecordResource struct {
	config *common.DnsimpleProviderConfig
}

// DsRecordResourceModel describes the resource data model.
type DsRecordResourceModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Domain     types.String `tfsdk:"domain"`
	Algorithm  types.String `tfsdk:"algorithm"`
	Digest     types.String `tfsdk:"digest"`
	DigestType types.String `tfsdk:"digest_type"`
	Keytag     types.String `tfsdk:"keytag"`
	PublicKey  types.String `tfsdk:"public_key"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (r *DsRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ds_record"
}

func (r *DsRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain delegation signer record resource",
		Attributes: map[string]schema.Attribute{
			"id": common.IDInt64Attribute(),
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"digest": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"digest_type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"keytag": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"public_key": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *DsRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DsRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DsRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	dsAttributes := dnsimple.DelegationSignerRecord{
		Algorithm:  data.Algorithm.ValueString(),
		Digest:     data.Digest.ValueString(),
		DigestType: data.DigestType.ValueString(),
		Keytag:     data.Keytag.ValueString(),
		PublicKey:  data.PublicKey.ValueString(),
	}

	response, err := r.config.Client.Domains.CreateDelegationSignerRecord(ctx, r.config.AccountID, data.Domain.ValueString(), dsAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple domain delegation signer record",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DsRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DsRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Domains.GetDelegationSignerRecord(ctx, r.config.AccountID, data.Domain.ValueString(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple domain delegation signer record: %d", data.Id.ValueInt64()),
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DsRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op
	tflog.Info(ctx, "delegation signer records cannot be updated")
}

func (r *DsRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DsRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple domain delegation signer record: %s", data.Id))

	_, err := r.config.Client.Domains.DeleteDelegationSignerRecord(ctx, r.config.AccountID, data.Domain.ValueString(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to delete DNSimple domain delegation signer record: %d", data.Id.ValueInt64()),
			err.Error(),
		)
	}
}

func (r *DsRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"resource import invalid ID",
			fmt.Sprintf("wrong format of import ID (%s), use: '<domain-name>_<delegation-signer-record-id>'", req.ID),
		)
	}
	domain := parts[0]
	dsIdRaw := parts[1]

	dsId, err := strconv.ParseInt(dsIdRaw, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to import DNSimple domain delegation signer record",
			fmt.Sprintf("invalid ID: %s", dsIdRaw),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), dsId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domain)...)
}

func (r *DsRecordResource) updateModelFromAPIResponse(ds *dnsimple.DelegationSignerRecord, data *DsRecordResourceModel) {
	data.Id = types.Int64Value(ds.ID)
	data.Algorithm = types.StringValue(ds.Algorithm)
	data.Digest = types.StringValue(ds.Digest)
	data.DigestType = types.StringValue(ds.DigestType)
	data.Keytag = types.StringValue(ds.Keytag)
	data.PublicKey = types.StringValue(ds.PublicKey)
	data.CreatedAt = types.StringValue(ds.CreatedAt)
	data.UpdatedAt = types.StringValue(ds.UpdatedAt)
}
