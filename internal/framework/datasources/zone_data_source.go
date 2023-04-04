package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ZoneDataSource{}

func NewZoneDataSource() datasource.DataSource {
	return &ZoneDataSource{}
}

// ZoneDataSource defines the data source implementation.
type ZoneDataSource struct {
	config *common.DnsimpleProviderConfig
}

// ZoneDataSourceModel describes the data source data model.
type ZoneDataSourceModel struct {
	Id        types.Int64  `tfsdk:"id"`
	AccountId types.Int64  `tfsdk:"account_id"`
	Name      types.String `tfsdk:"name"`
	Reverse   types.Bool   `tfsdk:"reverse"`
}

func (d *ZoneDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

func (d *ZoneDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple zone data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Zone ID",
				Computed:            true,
			},
			"account_id": schema.Int64Attribute{
				MarkdownDescription: "DNSimple Account ID to which the zone belongs to",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Zone Name",
				Required:            true,
			},
			"reverse": schema.BoolAttribute{
				MarkdownDescription: "True if the zone is a reverse zone",
				Computed:            true,
			},
		},
	}
}

func (d *ZoneDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ZoneDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ZoneDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := d.config.Client.Zones.GetZone(ctx, d.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to find DNSimple Zone",
			err.Error(),
		)
		return
	}

	data.Id = types.Int64Value(response.Data.ID)
	data.AccountId = types.Int64Value(response.Data.AccountID)
	data.Name = types.StringValue(response.Data.Name)
	data.Reverse = types.BoolValue(response.Data.Reverse)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
