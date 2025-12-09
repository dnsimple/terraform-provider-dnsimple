package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &DomainDelegationResource{}
	_ resource.ResourceWithConfigure   = &DomainDelegationResource{}
	_ resource.ResourceWithImportState = &DomainDelegationResource{}
)

func NewDomainDelegationResource() resource.Resource {
	return &DomainDelegationResource{}
}

// DomainDelegationResource defines the resource implementation.
type DomainDelegationResource struct {
	config *common.DnsimpleProviderConfig
}

// DomainDelegationResourceModel describes the resource data model.
type DomainDelegationResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	NameServers types.Set    `tfsdk:"name_servers"`
}

func (r *DomainDelegationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_delegation"
}

func (r *DomainDelegationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain delegation resource",
		Attributes: map[string]schema.Attribute{
			"id": common.IDStringAttribute(),
			"domain": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name_servers": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					modifiers.SetTrimSuffixValue(),
				},
			},
		},
	}
}

func (r *DomainDelegationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DomainDelegationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DomainDelegationResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	nameServers := dnsimple.Delegation{}
	resp.Diagnostics.Append(data.NameServers.ElementsAs(ctx, &nameServers, false)...)

	tflog.Debug(ctx, "creating domain delegation", map[string]interface{}{"name servers": nameServers})

	_, err := r.config.Client.Registrar.ChangeDomainDelegation(ctx, r.config.AccountID, data.Domain.ValueString(), &nameServers)
	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create domain delegation",
			err.Error(),
		)
		return
	}

	data.Id = data.Domain

	tflog.Info(ctx, "created domain delegation", map[string]interface{}{"domain": data.Domain})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainDelegationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Read Terraform prior state data into the model.
	var data *DomainDelegationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Registrar.GetDomainDelegation(ctx, r.config.AccountID, data.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read domain delegation for domain %s", data.Domain.ValueString()),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(r.updateModelFromAPIResponse(ctx, response.Data, data)...)

	// Save updated data into Terraform state.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainDelegationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DomainDelegationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nameServers := dnsimple.Delegation{}
	resp.Diagnostics.Append(data.NameServers.ElementsAs(ctx, &nameServers, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Registrar.ChangeDomainDelegation(ctx, r.config.AccountID, data.Domain.ValueString(), &nameServers)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to update domain delegation for domain %s", data.Domain.ValueString()),
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "domain delegation updated", map[string]interface{}{"data": response.Data})

	data.Id = data.Domain

	resp.Diagnostics.Append(r.updateModelFromAPIResponse(ctx, response.Data, data)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DomainDelegationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
	tflog.Info(ctx, "deleting a domain delegation simply deletes the Terraform state, and the current domain delegation will remain as is and no longer be managed by Terraform")
}

func (r *DomainDelegationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	domainId := req.ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain"), domainId)...)
}

func (r *DomainDelegationResource) updateModelFromAPIResponse(ctx context.Context, delegation *dnsimple.Delegation, data *DomainDelegationResourceModel) diag.Diagnostics {
	nameServers, diag := types.SetValueFrom(ctx, types.StringType, delegation)
	data.NameServers = nameServers
	return diag
}
