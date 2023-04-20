package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &LetsEncryptCertificateResource{}
	_ resource.ResourceWithConfigure   = &LetsEncryptCertificateResource{}
	_ resource.ResourceWithImportState = &LetsEncryptCertificateResource{}
)

func NewLetsEncryptCertificateResource() resource.Resource {
	return &LetsEncryptCertificateResource{}
}

// LetsEncryptCertificateResource defines the resource implementation.
type LetsEncryptCertificateResource struct {
	config *common.DnsimpleProviderConfig
}

// LetsEncryptCertificateResourceModel describes the resource data model.
type LetsEncryptCertificateResourceModel struct {
	Id                  types.Int64  `tfsdk:"id"`
	DomainId            types.String `tfsdk:"domain_id"`
	Name                types.String `tfsdk:"name"`
	AlternateNames      types.List   `tfsdk:"alternate_names"`
	Years               types.Int64  `tfsdk:"years"`
	State               types.String `tfsdk:"state"`
	AuthorityIdentifier types.String `tfsdk:"authority_identifier"`
	AutoRenew           types.Bool   `tfsdk:"auto_renew"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
	ExpiresAt           types.String `tfsdk:"expires_at"`
	Csr                 types.String `tfsdk:"csr"`
	SignatureAlgorithm  types.String `tfsdk:"signature_algorithm"`
}

func (r *LetsEncryptCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lets_encrypt_certificate"
}

func (r *LetsEncryptCertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple Let's Encrypt certificate resource",
		Attributes: map[string]schema.Attribute{
			"id": common.IDInt64Attribute(),
			"domain_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"alternate_names": schema.ListAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				ElementType: types.StringType,
			},
			"years": schema.Int64Attribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"authority_identifier": schema.StringAttribute{
				Computed: true,
			},
			"auto_renew": schema.BoolAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"expires_at": schema.StringAttribute{
				Computed: true,
			},
			"csr": schema.StringAttribute{
				Computed: true,
			},
			"signature_algorithm": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *LetsEncryptCertificateResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LetsEncryptCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *LetsEncryptCertificateResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	alternateNames := make([]string, len(data.AlternateNames.Elements()))
	resp.Diagnostics.Append(data.AlternateNames.ElementsAs(ctx, alternateNames, false)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.LetsencryptCertificateAttributes{
		AutoRenew:          data.AutoRenew.ValueBool(),
		Name:               data.Name.ValueString(),
		AlternateNames:     alternateNames,
		SignatureAlgorithm: data.SignatureAlgorithm.ValueString(),
	}

	tflog.Debug(ctx, "creating DNSimple LetsEncryptCertificate", map[string]interface{}{"attributes": domainAttributes})

	response, err := r.config.Client.Certificates.PurchaseLetsencryptCertificate(ctx, r.config.AccountID, data.DomainId.ValueString(), domainAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to purchase Let's Encrypt Certificate",
			err.Error(),
		)
		return
	}

	certificateId := response.Data.CertificateID

	issueResponse, issueErr := r.config.Client.Certificates.IssueLetsencryptCertificate(ctx, r.config.AccountID, data.DomainId.ValueString(), certificateId)

	if issueErr != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(issueErr, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to issue Let's Encrypt Certificate",
			issueErr.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(issueResponse.Data, data)

	tflog.Info(ctx, "purchased Let's Encrypt Certificate", map[string]interface{}{"id": data.Id})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LetsEncryptCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *LetsEncryptCertificateResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Certificates.GetCertificate(ctx, r.config.AccountID, data.DomainId.ValueString(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple LetsEncryptCertificate: %d", data.Id.ValueInt64()),
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LetsEncryptCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No-op
	tflog.Info(ctx, "DNSimple does not support updating Let's Encrypt certificates")
}

func (r *LetsEncryptCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No-op
	tflog.Info(ctx, "DNSimple does not support deleting Let's Encrypt certificates")
}

func (r *LetsEncryptCertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// No-op
	tflog.Info(ctx, "DNSimple does not support importing Let's Encrypt certificates")
}

func (r *LetsEncryptCertificateResource) updateModelFromAPIResponse(cert *dnsimple.Certificate, data *LetsEncryptCertificateResourceModel) {
	// Do not set data.DomainId to cert.DomainID, as that will cause Terraform to reject it with an inconsistent state error. Neither the domain name nor ID will ever change anyway.
	data.Id = types.Int64Value(cert.ID)
	data.Years = types.Int64Value(int64(cert.Years))
	data.State = types.StringValue(cert.State)
	data.AuthorityIdentifier = types.StringValue(cert.AuthorityIdentifier)
	data.AutoRenew = types.BoolValue(cert.AutoRenew)
	data.CreatedAt = types.StringValue(cert.CreatedAt)
	data.UpdatedAt = types.StringValue(cert.UpdatedAt)
	data.ExpiresAt = types.StringValue(cert.ExpiresAt)
	data.Csr = types.StringValue(cert.CertificateRequest)
}
