package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RegistrantChangeCheckDataSource{}

func NewRegistrantChangeCheckDataSource() datasource.DataSource {
	return &RegistrantChangeCheckDataSource{}
}

// RegistrantChangeCheckDataSource defines the data source implementation.
type RegistrantChangeCheckDataSource struct {
	config *common.DnsimpleProviderConfig
}

// RegistrantChangeCheckDataSourceModel describes the data source data model.
type RegistrantChangeCheckDataSourceModel struct {
	Id                  types.String        `tfsdk:"id"`
	ContactId           types.Int64         `tfsdk:"contact_id"`
	DomainId            types.String        `tfsdk:"domain_id"`
	ExtendedAttributes  []ExtendedAttribute `tfsdk:"extended_attributes"`
	RegistryOwnerChange types.Bool          `tfsdk:"registry_owner_change"`
}

type ExtendedAttribute struct {
	Name        types.String              `tfsdk:"name"`
	Description types.String              `tfsdk:"description"`
	Required    types.Bool                `tfsdk:"required"`
	Options     []ExtendedAttributeOption `tfsdk:"options"`
}

type ExtendedAttributeOption struct {
	Title       types.String `tfsdk:"title"`
	Value       types.String `tfsdk:"value"`
	Description types.String `tfsdk:"description"`
}

func (d *RegistrantChangeCheckDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registrant_change_check"
}

func (d *RegistrantChangeCheckDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple registrant change check data source",

		Attributes: map[string]schema.Attribute{
			"id": common.IDInt64Attribute(),
			"contact_id": schema.Int64Attribute{
				MarkdownDescription: "DNSimple contact ID for which the registrant change check is being performed",
				Required:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "DNSimple domain ID for which the registrant change check is being performed",
				Required:            true,
			},
			"extended_attributes": schema.ListAttribute{
				MarkdownDescription: "Extended attributes for the registrant change",
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name":        types.StringType,
						"description": types.StringType,
						"required":    types.BoolType,
						"options": types.ListType{
							ElemType: types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"title":       types.StringType,
									"value":       types.StringType,
									"description": types.StringType,
								},
							},
						},
					},
				},
			},
			"registry_owner_change": schema.BoolAttribute{
				MarkdownDescription: "True if the registrant change will result in a registry owner change",
				Computed:            true,
			},
		},
	}
}

func (d *RegistrantChangeCheckDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RegistrantChangeCheckDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RegistrantChangeCheckDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Replace with official DNSimple Go client
	response, err := d.config.TempClient.CheckRegistrantChange(data.DomainId.ValueString(), data.ContactId.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			"failed to check registrant change",
			err.Error(),
		)
		return
	}

	contactIdString := fmt.Sprintf("%d", response.Data.ContactID)
	data.Id = types.StringValue(data.DomainId.ValueString() + ":" + contactIdString)
	data.ContactId = types.Int64Value(response.Data.ContactID)
	data.DomainId = types.StringValue(fmt.Sprintf("%d", response.Data.DomainID))
	data.ExtendedAttributes = make([]ExtendedAttribute, len(response.Data.ExtendedAttributes))
	for i, extendedAttribute := range response.Data.ExtendedAttributes {
		data.ExtendedAttributes[i].Name = types.StringValue(extendedAttribute["name"].(string))
		data.ExtendedAttributes[i].Description = types.StringValue(extendedAttribute["description"].(string))
		data.ExtendedAttributes[i].Required = types.BoolValue(extendedAttribute["required"].(bool))
		data.ExtendedAttributes[i].Options = make([]ExtendedAttributeOption, len(extendedAttribute["options"].([]map[string]interface{})))
		for j, option := range extendedAttribute["options"].([]map[string]interface{}) {
			data.ExtendedAttributes[i].Options[j].Title = types.StringValue(option["title"].(string))
			data.ExtendedAttributes[i].Options[j].Value = types.StringValue(option["value"].(string))
			data.ExtendedAttributes[i].Options[j].Description = types.StringValue(option["description"].(string))
		}
	}
	data.RegistryOwnerChange = types.BoolValue(response.Data.RegistryOwnerChange)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
