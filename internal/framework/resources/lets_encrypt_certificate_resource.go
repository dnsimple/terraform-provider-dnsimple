package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
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
	DomainId            types.Int64  `tfsdk:"domain_id"`
	Name                types.String `tfsdk:"name"`
	Years               types.Int64  `tfsdk:"years"`
	State               types.String `tfsdk:"state"`
	AuthorityIdentifier types.String `tfsdk:"authority_identifier"`
	AutoRenew           types.Bool   `tfsdk:"auto_renew"`
	CreatedAt           types.String `tfsdk:"created_at"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
	ExpiresOn           types.String `tfsdk:"expires_on"`
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
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"domain_id": schema.StringAttribute{
				Required: true,
			},
			"contact_id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Required: true,
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
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
			"expires_on": schema.StringAttribute{
				Computed: true,
			},
			"csr": schema.StringAttribute{
				Computed: true,
			},
			"signature_algorithm": schema.StringAttribute{
				Optional: true,
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

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.LetsencryptCertificateAttributes{
		AutoRenew:          data.AutoRenew.ValueBool(),
		Name:               data.Name.ValueString(),
		SignatureAlgorithm: data.SignatureAlgorithm.ValueString(),
	}

	tflog.Debug(ctx, "creating DNSimple LetsEncryptCertificate", map[string]interface{}{"attributes": domainAttributes})

	response, err := r.config.Client.Certificates.PurchaseLetsencryptCertificate(ctx, r.config.AccountID, data.DomainId.String(), domainAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(attributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to purchase Let's Encrypt Certificate",
			err.Error(),
		)
		return
	}

	certificateId := response.Data.CertificateID

	issueResponse, issueErr := r.config.Client.Certificates.IssueLetsencryptCertificate(ctx, r.config.AccountID, data.DomainId.String(), certificateId)

	if issueErr != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(issueErr, &errorResponse) {
			resp.Diagnostics.Append(attributeErrorsToDiagnostics(errorResponse)...)
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

	response, err := r.config.Client.Certificates.GetCertificate(ctx, r.config.AccountID, data.DomainId.String(), data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple LetsEncryptCertificate: %d", data.Id.ValueInt64()),
			err.Error(),
		)
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *LetsEncryptCertificateResource) updateModelFromAPIResponse(cert *dnsimple.Certificate, data *LetsEncryptCertificateResourceModel) {
	data.Id = types.Int64Value(cert.ID)
	data.DomainId = types.Int64Value(cert.DomainID)
	data.Years = types.Int64Value(int64(cert.Years))
	data.State = types.StringValue(cert.State)
	data.AuthorityIdentifier = types.StringValue(cert.AuthorityIdentifier)
	data.AutoRenew = types.BoolValue(cert.AutoRenew)
	data.CreatedAt = types.StringValue(cert.CreatedAt)
	data.UpdatedAt = types.StringValue(cert.UpdatedAt)
	data.Csr = types.StringValue(cert.CertificateRequest)
}
