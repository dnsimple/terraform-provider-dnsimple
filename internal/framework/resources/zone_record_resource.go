package resources

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &ZoneRecordResource{}
	_ resource.ResourceWithConfigure      = &ZoneRecordResource{}
	_ resource.ResourceWithImportState    = &ZoneRecordResource{}
	_ resource.ResourceWithValidateConfig = &ZoneRecordResource{}
	_ resource.ResourceWithModifyPlan     = &ZoneRecordResource{}
)

func NewZoneRecordResource() resource.Resource {
	return &ZoneRecordResource{}
}

// ZoneRecordResource defines the resource implementation.
type ZoneRecordResource struct {
	config *common.DnsimpleProviderConfig
}

// ZoneRecordResourceModel describes the resource data model.
type ZoneRecordResourceModel struct {
	ZoneName        types.String `tfsdk:"zone_name"`
	ZoneId          types.String `tfsdk:"zone_id"`
	Name            types.String `tfsdk:"name"`
	NameNormalized  types.String `tfsdk:"name_normalized"`
	QualifiedName   types.String `tfsdk:"qualified_name"`
	Type            types.String `tfsdk:"type"`
	Regions         types.List   `tfsdk:"regions"`
	Value           types.String `tfsdk:"value"`
	ValueNormalized types.String `tfsdk:"value_normalized"`
	TTL             types.Int64  `tfsdk:"ttl"`
	Priority        types.Int64  `tfsdk:"priority"`
	Id              types.Int64  `tfsdk:"id"`
}

func (r *ZoneRecordResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone_record"
}

// checkPriorityAgainstRecordType reports a diagnostic when a record carries a priority
// its type cannot hold.
//
// The DNSimple API stores a priority for MX and SRV records. For every other type it
// discards the submitted value and returns null, which the client decodes as 0. Left
// unchecked, the apply reaches Terraform's consistency check with a planned priority the
// API never accepted and fails with "Provider produced inconsistent result after apply",
// which reads as a provider defect rather than a configuration mistake.
//
// A priority of 0 passes for every type. It already matches what the API reports back, so
// it applies cleanly today and rejecting it would break working configurations without
// fixing anything.
//
// The MX/SRV set is about the API's separate priority attribute, which is a legacy of
// splitting priority out of the record content. Newer priority-bearing types such as SVCB
// and HTTPS carry their priority inside the content instead, so they do not belong here
// and are correctly rejected. Revisit this set only if the API starts populating the
// priority attribute for a further type; see dnsimple/terraform-provider-dnsimple#360.
//
// Either value may still be unknown when this runs, in which case there is nothing to
// compare and the caller closer to apply performs the check instead.
func checkPriorityAgainstRecordType(recordType types.String, priority types.Int64) diag.Diagnostics {
	var diags diag.Diagnostics

	if priority.IsNull() || priority.IsUnknown() || recordType.IsNull() || recordType.IsUnknown() {
		return diags
	}

	if priority.ValueInt64() == 0 {
		return diags
	}

	switch strings.ToUpper(recordType.ValueString()) {
	case "MX", "SRV":
		return diags
	}

	diags.AddAttributeError(
		path.Root("priority"),
		"priority is not supported for this record type",
		fmt.Sprintf(
			"DNSimple only stores a priority for MX and SRV records, and discards it for %s records. Remove the priority attribute from this resource.",
			recordType.ValueString(),
		),
	)

	return diags
}

// ValidateConfig catches the common case, where both values are written literally, and
// reports it before any plan is produced.
func (r *ZoneRecordResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ZoneRecordResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(checkPriorityAgainstRecordType(data.Type, data.Priority)...)
}

// ModifyPlan catches values that were unknown during configuration validation.
//
// That covers references to resources that already exist, which resolve by planning time,
// and also references to resources created in the same run: Terraform plans each resource
// again during apply once its dependencies have resolved, so the real value reaches this
// hook before the API is called. Create and Update therefore need no separate guard.
func (r *ZoneRecordResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// A destroy plan has no planned values to check.
	if req.Plan.Raw.IsNull() {
		return
	}

	var data ZoneRecordResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(checkPriorityAgainstRecordType(data.Type, data.Priority)...)
}

func (r *ZoneRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain resource",
		Version:             1,
		Attributes: map[string]schema.Attribute{
			"zone_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"zone_id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"name_normalized": schema.StringAttribute{
				Computed: true,
			},
			"qualified_name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.RecordType{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"regions": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"value": schema.StringAttribute{
				Required: true,
			},
			"value_normalized": schema.StringAttribute{
				Computed: true,
			},
			"ttl": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					modifiers.Int64DefaultValue(3600),
				},
			},
			"priority": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *ZoneRecordResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ZoneRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ZoneRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	regions := make([]string, len(data.Regions.Elements()))
	resp.Diagnostics.Append(data.Regions.ElementsAs(ctx, &regions, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordAttributes := dnsimple.ZoneRecordAttributes{
		Name:    dnsimple.String(data.Name.ValueString()),
		Type:    data.Type.ValueString(),
		Content: data.Value.ValueString(),
		Regions: regions,
		TTL:     int(data.TTL.ValueInt64()),
	}

	if !data.Priority.IsNull() {
		recordAttributes.Priority = int(data.Priority.ValueInt64())
	}

	tflog.Debug(ctx, "DNSimple Zone Record recordAttributes", map[string]interface{}{
		"attributes": recordAttributes,
	})

	response, err := r.config.Client.Zones.CreateRecord(
		ctx,
		r.config.AccountID,
		data.ZoneName.ValueString(),
		recordAttributes,
	)
	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple Zone Record",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	tflog.Info(ctx, "DNSimple Record ID", map[string]interface{}{"id": data.Id})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var (
		data *ZoneRecordResourceModel

		record dnsimple.ZoneRecord

		skip_prefetch_cache bool = false
	)

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check if we should skip the cache prefetch
	if value, diags := req.Private.GetKey(ctx, "skip_prefetch_cache"); diags.HasError() {
		resp.Diagnostics.Append(diags...)
	} else {
		if string(value) == "true" {
			skip_prefetch_cache = true
		}
	}

	if r.config.Prefetch && !skip_prefetch_cache {
		if _, ok := r.config.ZoneRecordCache.Get(data.ZoneName.ValueString()); !ok {
			err := r.config.ZoneRecordCache.Hydrate(ctx, r.config.Client, r.config.AccountID, data.ZoneName.ValueString(), nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"failed to hydrate zone record cache",
					err.Error(),
				)
				return
			}
		}

		var lookupName string
		if data.NameNormalized.IsNull() || data.NameNormalized.IsUnknown() {
			lookupName = data.Name.ValueString()
		} else {
			lookupName = data.NameNormalized.ValueString()
		}

		cacheRecord, ok := r.config.ZoneRecordCache.Find(data.ZoneName.ValueString(), lookupName, data.Type.ValueString(), data.ValueNormalized.ValueString())
		if !ok {
			resp.Diagnostics.AddError(
				"failed to read DNSimple Zone Record",
				fmt.Sprintf("Zone record not found in cache for qualified name '%s'", data.QualifiedName.ValueString()),
			)
			return
		}

		tflog.Debug(ctx, "DNSimple Zone Record cache hit", map[string]interface{}{
			"zone_name": data.ZoneName.ValueString(),
		})

		record = cacheRecord
	} else {
		tflog.Debug(ctx, "DNSimple Zone Record cache miss", map[string]interface{}{
			"zone_name": data.ZoneName.ValueString(),
		})

		response, err := r.config.Client.Zones.GetRecord(ctx, r.config.AccountID, data.ZoneName.ValueString(), data.Id.ValueInt64())
		if err != nil {
			var errorResponse *dnsimple.ErrorResponse
			if errors.As(err, &errorResponse) {
				if errorResponse.Response.HTTPResponse.StatusCode == http.StatusNotFound {
					tflog.Warn(ctx, "removing zone record from state because it is not present in the remote")
					resp.State.RemoveResource(ctx)
					return
				}
			}

			resp.Diagnostics.AddError(
				"failed to read DNSimple Zone Record",
				fmt.Sprintf("Unable to read zone record with ID %d: %s", data.Id.ValueInt64(), err.Error()),
			)
			return
		}

		record = *response.Data
	}

	if record.Content != data.ValueNormalized.ValueString() {
		// If the record content has changed, we need to update the record in the remote
		tflog.Debug(ctx, "DNSimple Zone Record content changed")
		data.Value = types.StringValue(record.Content)
	}

	r.updateModelFromAPIResponse(&record, data)

	// Clear the private key to avoid reusing it in the next request after the import
	if skip_prefetch_cache {
		resp.Private.SetKey(ctx, "skip_prefetch_cache", nil)
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ZoneRecordResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	regions := make([]string, len(data.Regions.Elements()))
	resp.Diagnostics.Append(data.Regions.ElementsAs(ctx, &regions, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	recordAttributes := dnsimple.ZoneRecordAttributes{
		Name:    dnsimple.String(data.Name.ValueString()),
		Type:    data.Type.ValueString(),
		Content: data.Value.ValueString(),
		Regions: regions,
		TTL:     int(data.TTL.ValueInt64()),
	}

	if !data.Priority.IsNull() {
		recordAttributes.Priority = int(data.Priority.ValueInt64())
	}

	tflog.Debug(ctx, fmt.Sprintf("DNSimple Zone Record updateRecordAttributes: %+v", recordAttributes))

	response, err := r.config.Client.Zones.UpdateRecord(
		ctx,
		r.config.AccountID,
		data.ZoneName.ValueString(),
		data.Id.ValueInt64(),
		recordAttributes,
	)
	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		fmt.Printf("error: %+v", err)

		resp.Diagnostics.AddError(
			"failed to update DNSimple Zone Record",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ZoneRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ZoneRecordResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple Record: %s, %d", data.ZoneName, data.Id))

	_, err := r.config.Client.Zones.DeleteRecord(ctx, r.config.AccountID, data.ZoneName.ValueString(), data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to delete DNSimple Zone Record",
			fmt.Sprintf("Unable to delete zone record '%s' (ID: %d): %s", data.Name.ValueString(), data.Id.ValueInt64(), err.Error()),
		)
		return
	}
}

func (r *ZoneRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"invalid import ID",
			fmt.Sprintf("Invalid import ID format '%s'. Expected format: '<zone-name>_<record-id>'", req.ID),
		)
		return
	}
	zoneName := parts[0]
	recordID := parts[1]

	id, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"invalid import ID",
			fmt.Sprintf("Unable to parse record ID '%s' as integer. Expected a numeric ID", recordID),
		)
		return
	}

	resp.Private.SetKey(ctx, "skip_prefetch_cache", []byte(`true`))

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("zone_name"), zoneName)...)
}

func (r *ZoneRecordResource) updateModelFromAPIResponse(record *dnsimple.ZoneRecord, data *ZoneRecordResourceModel) {
	data.Id = types.Int64Value(record.ID)
	data.ZoneId = types.StringValue(record.ZoneID)
	data.NameNormalized = types.StringValue(record.Name)
	data.Type = types.StringValue(record.Type)
	data.ValueNormalized = types.StringValue(record.Content)
	data.TTL = types.Int64Value(int64(record.TTL))
	data.Priority = types.Int64Value(int64(record.Priority))

	if data.Name.IsNull() || data.Name.IsUnknown() {
		// This can happen during a resource import, where the name is not in the state
		data.Name = types.StringValue(record.Name)
	}

	if record.Name == "" {
		data.QualifiedName = data.ZoneName
	} else {
		data.QualifiedName = types.StringValue(fmt.Sprintf("%s.%s", record.Name, data.ZoneName.ValueString()))
	}
}
