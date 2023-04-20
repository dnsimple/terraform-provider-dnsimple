package provider

import (
	"context"
	"fmt"

	"github.com/dnsimple/dnsimple-go/dnsimple"
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
	Token          types.String `tfsdk:"token"`
	Account        types.String `tfsdk:"account"`
	Sandbox        types.Bool   `tfsdk:"sandbox"`
	Prefetch       types.Bool   `tfsdk:"prefetch"`
	UserAgentExtra types.String `tfsdk:"user_agent"`
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
			"must provide api token",
			"must define a token for the dnsimple provider or set the DNSIMPLE_TOKEN environment variable",
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
			"must provide account",
			"must define an account for the dnsimple provider or set the DNSIMPLE_ACCOUNT environment variable",
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

	if data.Prefetch.IsNull() {
		prefetchValue := utils.GetDefaultFromEnv("DNSIMPLE_PREFETCH", "")
		if prefetchValue == "" {
			prefetch = false
		}
	} else {
		prefetch = data.Prefetch.ValueBool()
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := dnsimple.NewClient(tc)

	userAgent := fmt.Sprintf("terraform/%s terraform-provider-dnsimple/%s", req.TerraformVersion, p.version)
	if data.UserAgentExtra.ValueString() != "" {
		userAgent = fmt.Sprintf("%s %s", userAgent, data.UserAgentExtra.ValueString())
	}
	client.SetUserAgent(userAgent)

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
		resources.NewDsRecordResource,
		resources.NewEmailForwardResource,
		resources.NewLetsEncryptCertificateResource,
		resources.NewZoneRecordResource,
	}
}

// DataSources returns the data sources supported by this provider.
func (p *DnsimpleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewCertificateDataSource,
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
