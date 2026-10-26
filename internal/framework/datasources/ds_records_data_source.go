package datasources

import (
	"context"
	"fmt"
	"sort"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DsRecordsDataSource{}

func NewDsRecordsDataSource() datasource.DataSource {
	return &DsRecordsDataSource{}
}

// DsRecordsDataSource defines the data source implementation.
type DsRecordsDataSource struct {
	config *common.DnsimpleProviderConfig
}

// DsRecordsDataSourceModel describes the data source data model.
type DsRecordsDataSourceModel struct {
	Id      types.String    `tfsdk:"id"`
	Domain  types.String    `tfsdk:"domain"`
	Records []DsRecordModel `tfsdk:"records"`
}

// DsRecordModel describes a single delegation signer record in the data source.
// Attribute names and types mirror dnsimple_ds_record so data source results
// can be compared against that resource's attributes.
type DsRecordModel struct {
	Id         types.Int64  `tfsdk:"id"`
	Algorithm  types.String `tfsdk:"algorithm"`
	Digest     types.String `tfsdk:"digest"`
	DigestType types.String `tfsdk:"digest_type"`
	Keytag     types.String `tfsdk:"key_tag"`
	PublicKey  types.String `tfsdk:"public_key"`
	CreatedAt  types.String `tfsdk:"created_at"`
	UpdatedAt  types.String `tfsdk:"updated_at"`
}

func (d *DsRecordsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ds_records"
}

func (d *DsRecordsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple delegation signer (DS) records for a domain. " +
			"Use this to copy the DS records of a DNSSEC-signed zone hosted at DNSimple to the domain's registrar.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The domain name or ID",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The name or ID of the domain",
				Required:            true,
			},
			"records": schema.ListNestedAttribute{
				MarkdownDescription: "The delegation signer records for the domain",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							MarkdownDescription: "The DS record ID in DNSimple",
							Computed:            true,
						},
						"algorithm": schema.StringAttribute{
							MarkdownDescription: "The DNSSEC algorithm number used to generate the signature (e.g. `8`, `13`)",
							Computed:            true,
						},
						"digest": schema.StringAttribute{
							MarkdownDescription: "The digest of the DS record",
							Computed:            true,
						},
						"digest_type": schema.StringAttribute{
							MarkdownDescription: "The digest type of the DS record (e.g. `2` for SHA-256)",
							Computed:            true,
						},
						"key_tag": schema.StringAttribute{
							MarkdownDescription: "The key tag of the DS record",
							Computed:            true,
						},
						"public_key": schema.StringAttribute{
							MarkdownDescription: "The public key of the DS record",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "When the DS record was created at DNSimple",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "When the DS record was last updated at DNSimple",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DsRecordsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*common.DnsimpleProviderConfig)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *common.DnsimpleProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.config = config
}

func (d *DsRecordsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DsRecordsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// A KSK rollover temporarily leaves two DS records on the domain, so fetch
	// every page before returning. An empty result (unsigned domain) is valid.
	options := &dnsimple.ListOptions{PerPage: dnsimple.Int(100)}

	var records []dnsimple.DelegationSignerRecord

	for {
		response, err := d.config.Client.Domains.ListDelegationSignerRecords(ctx, d.config.AccountID, data.Domain.ValueString(), options)
		if err != nil {
			resp.Diagnostics.AddError(
				"failed to read DNSimple DS Records",
				fmt.Sprintf("Unable to list DS records for domain %s: %s", data.Domain.ValueString(), err.Error()),
			)
			return
		}

		records = append(records, response.Data...)

		if response.Pagination.CurrentPage >= response.Pagination.TotalPages {
			break
		}

		options.Page = dnsimple.Int(response.Pagination.CurrentPage + 1)
	}

	// Sort by ID ascending so the result is deterministic regardless of API ordering.
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID < records[j].ID
	})

	data.Records = make([]DsRecordModel, 0, len(records))
	for _, record := range records {
		data.Records = append(data.Records, DsRecordModel{
			Id:         types.Int64Value(record.ID),
			Algorithm:  types.StringValue(record.Algorithm),
			Digest:     types.StringValue(record.Digest),
			DigestType: types.StringValue(record.DigestType),
			Keytag:     types.StringValue(record.Keytag),
			PublicKey:  types.StringValue(record.PublicKey),
			CreatedAt:  types.StringValue(record.CreatedAt),
			UpdatedAt:  types.StringValue(record.UpdatedAt),
		})
	}

	data.Id = data.Domain

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
