package datasources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CertificateDataSource{}

func NewCertificateDataSource() datasource.DataSource {
	return &CertificateDataSource{}
}

// CertificateDataSource defines the data source implementation.
type CertificateDataSource struct {
	config *common.DnsimpleProviderConfig
}

// CertificateDataSourceModel describes the data source data model.
type CertificateDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	CertificateId     types.Int64  `tfsdk:"certificate_id"`
	Domain            types.String `tfsdk:"domain"`
	ServerCertificate types.String `tfsdk:"server_certificate"`
	RootCertificate   types.String `tfsdk:"root_certificate"`
	CertificateChain  types.List   `tfsdk:"certificate_chain"`
	PrivateKey        types.String `tfsdk:"private_key"`
}

func (d *CertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (d *CertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple certificate data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Certificate data ID",
				Computed:            true,
			},
			"certificate_id": schema.Int64Attribute{
				MarkdownDescription: "Certificate ID",
				Required:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "Domain name",
				Required:            true,
			},
			"server_certificate": schema.StringAttribute{
				MarkdownDescription: "Server certificate",
				Computed:            true,
			},
			"root_certificate": schema.StringAttribute{
				MarkdownDescription: "Root certificate",
				Computed:            true,
			},
			"certificate_chain": schema.StringAttribute{
				MarkdownDescription: "Certificate chain",
				Computed:            true,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Private key",
				Computed:            true,
			},
		},
	}
}

func (d *CertificateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.config = config
}

func (d *CertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CertificateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.config.Client.Certificates.DownloadCertificate(ctx, d.config.AccountID, data.Domain.ValueString(), data.CertificateId.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to download DNSimple Certificate",
			err.Error(),
		)
		return
	}

	data.ServerCertificate = types.StringValue(response.Data.ServerCertificate)
	data.RootCertificate = types.StringValue(response.Data.RootCertificate)
	chain, diag := types.ListValueFrom(ctx, types.StringType, response.Data.IntermediateCertificates)
	if err != nil {
		resp.Diagnostics.Append(diag...)
		return
	}
	data.CertificateChain = chain

	response, err = d.config.Client.Certificates.GetCertificatePrivateKey(ctx, d.config.AccountID, data.Domain.ValueString(), data.CertificateId.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to download DNSimple Certificate private key",
			err.Error(),
		)
		return
	}

	data.PrivateKey = types.StringValue(response.Data.PrivateKey)
	data.Id = types.StringValue(time.Now().UTC().String())

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
