package datasources

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

const (
	CertificateConverged = "certificate_converged"
	CertificateFailed    = "certificate_failed"
	CertificateTimeout   = "certificate_timeout"
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
	Id                types.String   `tfsdk:"id"`
	CertificateId     types.Int64    `tfsdk:"certificate_id"`
	Domain            types.String   `tfsdk:"domain"`
	ServerCertificate types.String   `tfsdk:"server_certificate"`
	RootCertificate   types.String   `tfsdk:"root_certificate"`
	CertificateChain  types.List     `tfsdk:"certificate_chain"`
	PrivateKey        types.String   `tfsdk:"private_key"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

func (d *CertificateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate"
}

func (d *CertificateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple certificate data source",

		Attributes: map[string]schema.Attribute{
			"id": common.IDStringAttribute(),
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
			"certificate_chain": schema.ListAttribute{
				MarkdownDescription: "Certificate chain",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"private_key": schema.StringAttribute{
				MarkdownDescription: "Private key",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx),
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
	var data *CertificateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	convergenceState, err := tryToConvergeCertificate(ctx, data, &resp.Diagnostics, d, data.CertificateId.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"failed to get certificate state",
			err.Error(),
		)
		return
	}

	if convergenceState == CertificateFailed || convergenceState == CertificateTimeout {
		// Response is already populated with the error we can safely return
		return
	}

	if convergenceState == CertificateConverged {

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
		if diag.HasError() {
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
		data.Id = types.StringValue(idFromCertificateChain(data.ServerCertificate.ValueString(), data.RootCertificate.ValueString(), response.Data.IntermediateCertificates))

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	}
}

func tryToConvergeCertificate(ctx context.Context, data *CertificateDataSourceModel, diagnostics *diag.Diagnostics, d *CertificateDataSource, certificateID int64) (string, error) {
	readTimeout, diags := data.Timeouts.Read(ctx, 5*time.Minute)

	diagnostics.Append(diags...)

	if diagnostics.HasError() {
		return CertificateFailed, nil
	}

	err := utils.RetryWithTimeout(ctx, func() (error, bool) {
		certificate, err := d.config.Client.Certificates.GetCertificate(ctx, d.config.AccountID, data.Domain.ValueString(), data.CertificateId.ValueInt64())
		if err != nil {
			return err, false
		}

		if certificate.Data.State == consts.CertificateStateFailed {
			diagnostics.AddError(
				fmt.Sprintf("failed to issue certificate: %s", data.Domain.ValueString()),
				"certificate order failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if certificate.Data.State == consts.CertificateStateCancelled || certificate.Data.State == consts.CertificateStateRefunded {
			diagnostics.AddError(
				fmt.Sprintf("failed to issue certificate: %s", data.Domain.ValueString()),
				"certificate order failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if certificate.Data.State != consts.CertificateStateIssued {
			tflog.Info(ctx, fmt.Sprintf("[RETRYING] Certificate order is not complete, current state: %s", certificate.Data.State))

			return fmt.Errorf("certificate has not been issued, current state: %s. You can try to run terraform again to try and converge the certificate", certificate.Data.State), false
		}

		return nil, false
	}, readTimeout, 20*time.Second)

	if diagnostics.HasError() {
		// If we have diagnostic errors, we suspended the retry loop because the certificate is in a bad state, and cannot converge.
		return CertificateFailed, nil
	}

	if err != nil {
		// If we have an error, it means the retry loop timed out, and we cannot converge during this run.
		return CertificateTimeout, err
	}

	return CertificateConverged, nil
}

// idFromCertificateChain generates a SHA1 hash from the certificate chain.
func idFromCertificateChain(ServerCertificate, rootCertificate string, intermediateCertificateChain []string) string {
	// Concatenate all certificates into a single string
	certChain := ServerCertificate + rootCertificate + strings.Join(intermediateCertificateChain, "")

	// Create a new SHA1 hash.
	h := sha1.New()

	// Write the certificate chain string to the hash.
	h.Write([]byte(certChain))
	hashedCertChain := hex.EncodeToString(h.Sum(nil))

	return hashedCertChain
}
