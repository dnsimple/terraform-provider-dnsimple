package provider

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/datasources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources/registered_domain"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
	"golang.org/x/oauth2"
)

var _ provider.Provider = &DnsimpleProvider{}

type DnsimpleProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DnsimpleProviderModel describes the provider data model.
type DnsimpleProviderModel struct {
	Token              types.String `tfsdk:"token"`
	Account            types.String `tfsdk:"account"`
	Sandbox            types.Bool   `tfsdk:"sandbox"`
	Prefetch           types.Bool   `tfsdk:"prefetch"`
	UserAgentExtra     types.String `tfsdk:"user_agent"`
	DebugTransportFile types.String `tfsdk:"debug_transport_file"`
}

// debugTransport is an HTTP transport that logs requests and responses
type debugTransport struct {
	Base     http.RoundTripper
	FilePath string
}

func (t *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Open debug file for appending
	f, err := os.OpenFile(t.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()

		// Log request
		fmt.Fprintf(f, "\n=== HTTP REQUEST ===\n")
		fmt.Fprintf(f, "Method: %s\n", req.Method)
		fmt.Fprintf(f, "URL: %s\n", req.URL.String())
		fmt.Fprintf(f, "Headers: %v\n", req.Header)

		// Log request body if present
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Fprintf(f, "Error reading request body: %v\n", err)
			} else {
				fmt.Fprintf(f, "Request Body: %s\n", string(bodyBytes))
				// Restore the body for the actual request
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
	}

	// Perform the request
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	resp, respErr := base.RoundTrip(req)

	if respErr != nil {
		if f != nil {
			fmt.Fprintf(f, "=== HTTP ERROR ===\n")
			fmt.Fprintf(f, "Error: %v\n", respErr)
		}
		return resp, respErr
	}

	if f != nil {
		// Log response
		fmt.Fprintf(f, "=== HTTP RESPONSE ===\n")
		fmt.Fprintf(f, "Status: %d %s\n", resp.StatusCode, resp.Status)
		fmt.Fprintf(f, "Headers: %v\n", resp.Header)

		// Read and log body
		if resp.Body != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Fprintf(f, "Error reading response body: %v\n", err)
			} else {
				fmt.Fprintf(f, "Body: %s\n", string(bodyBytes))
				// Restore the body for the actual client to read
				resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}
		fmt.Fprintf(f, "=== END RESPONSE ===\n")
	}

	return resp, respErr
}

// Metadata returns information about the provider.
func (p *DnsimpleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dnsimple"

	resp.Version = p.version
}

// Schema returns the schema for this provider's configuration.
func (p *DnsimpleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The API v2 token for API operations.",
				Sensitive:           true,
			},
			"account": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The account for API operations.",
			},
			"sandbox": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Flag to enable the sandbox API.",
			},
			"prefetch": schema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Flag to enable the prefetch of zone records.",
			},
			"user_agent": schema.StringAttribute{
				Optional:    true,
				Description: "Custom string to append to the user agent used for sending HTTP requests to the API.",
			},
			"debug_transport_file": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "File path to enable HTTP request/response debugging. When set, all HTTP requests and responses will be logged to this file.",
			},
		},
		MarkdownDescription: "The DNSimple provider is used to interact with the various services that DNSimple offers. " +
			"The provider needs to be configured with the proper credentials before it can be used.",
	}
}

// Configure configures the provider.
func (p *DnsimpleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var (
		data DnsimpleProviderModel

		token    string
		account  string
		sandbox  bool
		prefetch bool
	)

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Fallback to env vars if client settings not set
	if data.Token.IsNull() {
		token = utils.GetDefaultFromEnv("DNSIMPLE_TOKEN", "")
	} else {
		token = data.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"failed to configure DNSimple provider",
			"API token is required. Set the token in the provider configuration or the DNSIMPLE_TOKEN environment variable",
		)
		return
	}

	if data.Account.IsNull() {
		account = utils.GetDefaultFromEnv("DNSIMPLE_ACCOUNT", "")
	} else {
		account = data.Account.ValueString()
	}

	if account == "" {
		resp.Diagnostics.AddError(
			"failed to configure DNSimple provider",
			"Account ID is required. Set the account in the provider configuration or the DNSIMPLE_ACCOUNT environment variable",
		)
		return
	}

	if data.Sandbox.IsNull() {
		sandboxValue := utils.GetDefaultFromEnv("DNSIMPLE_SANDBOX", "")
		if sandboxValue == "true" {
			sandbox = true
		}
	} else {
		sandbox = data.Sandbox.ValueBool()
	}

	if data.Prefetch.IsNull() || data.Prefetch.IsUnknown() {
		prefetchValue := utils.GetDefaultFromEnv("DNSIMPLE_PREFETCH", "")
		prefetch = prefetchValue != ""
	} else {
		prefetch = data.Prefetch.ValueBool()
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// Add debug transport to log HTTP requests/responses if debug_transport_file is set
	if !data.DebugTransportFile.IsNull() && !data.DebugTransportFile.IsUnknown() {
		debugFile := data.DebugTransportFile.ValueString()
		if debugFile != "" {
			tc.Transport = &debugTransport{Base: tc.Transport, FilePath: debugFile}
		}
	}

	client := dnsimple.NewClient(tc)

	userAgent := fmt.Sprintf("terraform/%s terraform-provider-dnsimple/%s", req.TerraformVersion, p.version)
	if data.UserAgentExtra.ValueString() != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, data.UserAgentExtra.ValueString())
	}
	client.SetUserAgent(userAgent)
	client.Debug = true

	if sandbox {
		client.BaseURL = consts.BaseURLSandbox
	}

	providerData := &common.DnsimpleProviderConfig{
		Client:          client,
		AccountID:       account,
		Prefetch:        prefetch,
		ZoneRecordCache: common.ZoneRecordCache{},
	}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *DnsimpleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewContactResource,
		resources.NewDomainDelegationResource,
		resources.NewDomainResource,
		registered_domain.NewRegisteredDomainResource,
		resources.NewDsRecordResource,
		resources.NewEmailForwardResource,
		resources.NewLetsEncryptCertificateResource,
		resources.NewZoneRecordResource,
		resources.NewZoneResource,
	}
}

// DataSources returns the data sources supported by this provider.
func (p *DnsimpleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewCertificateDataSource,
		datasources.NewRegistrantChangeCheckDataSource,
		datasources.NewZoneDataSource,
	}
}

// New returns a new provider factory for the DNSimple provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DnsimpleProvider{
			version: version,
		}
	}
}

func NewProto6ProviderFactory() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"dnsimple": providerserver.NewProtocol6WithError(New("test")()),
	}
}
