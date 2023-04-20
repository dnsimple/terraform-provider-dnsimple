package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
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
	_ resource.Resource                = &EmailForwardResource{}
	_ resource.ResourceWithConfigure   = &EmailForwardResource{}
	_ resource.ResourceWithImportState = &EmailForwardResource{}
)

func NewEmailForwardResource() resource.Resource {
	return &EmailForwardResource{}
}

// EmailForwardResource defines the resource implementation.
type EmailForwardResource struct {
	config *common.DnsimpleProviderConfig
}

// EmailForwardResourceModel describes the resource data model.
type EmailForwardResourceModel struct {
	Domain           types.String `tfsdk:"domain"`
	AliasName        types.String `tfsdk:"alias_name"`
	AliasEmail       types.String `tfsdk:"alias_email"`
	DestinationEmail types.String `tfsdk:"destination_email"`
	Id               types.Int64  `tfsdk:"id"`
}

func (r *EmailForwardResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_email_forward"
}

func (r *EmailForwardResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple email forward resource",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alias_email": schema.StringAttribute{
				Computed: true,
			},
			"destination_email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *EmailForwardResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmailForwardResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *EmailForwardResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.EmailForward{
		From: data.AliasName.ValueString(),
		To:   data.DestinationEmail.ValueString(),
	}

	tflog.Debug(ctx, "creating DNSimple EmailForward", map[string]interface{}{"attributes": domainAttributes})

	response, err := r.config.Client.Domains.CreateEmailForward(ctx, r.config.AccountID, data.Domain.ValueString(), domainAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple EmailForward",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	tflog.Info(ctx, "created DNSimple EmailForward", map[string]interface{}{"id": data.Id.ValueInt64()})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailForwardResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *EmailForwardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Domains.GetEmailForward(ctx, r.config.AccountID, data.Domain.ValueString(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple EmailForward: %d", data.Id.ValueInt64()),
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EmailForwardResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op
	tflog.Info(ctx, "DNSimple does not support updating email forwards")
}

func (r *EmailForwardResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *EmailForwardResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple EmailForward: %d", data.Id.ValueInt64()))

	_, err := r.config.Client.Domains.DeleteEmailForward(ctx, r.config.AccountID, data.Domain.ValueString(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to delete DNSimple EmailForward: %d", data.Id.ValueInt64()),
			err.Error(),
		)
		return
	}
}

func (r *EmailForwardResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("resource import invalid ID", fmt.Sprintf("wrong format of import ID (%s), use: '<domain-name>_<email-forward-id>'", req.ID))
		return
	}
	domainName := parts[0]
	recordID := parts[1]

	id, err := strconv.ParseInt(recordID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("resource import invalid ID", fmt.Sprintf("failed to parse email forward ID (%s) as integer", recordID))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domainName)...)
}

func (r *EmailForwardResource) updateModelFromAPIResponse(emailForward *dnsimple.EmailForward, data *EmailForwardResourceModel) {
	data.Id = types.Int64Value(emailForward.ID)
	data.AliasName = types.StringValue(strings.Split(emailForward.From, "@")[0])
	data.AliasEmail = types.StringValue(emailForward.From)
	data.DestinationEmail = types.StringValue(emailForward.To)
}
