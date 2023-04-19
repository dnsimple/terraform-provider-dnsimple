package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ContactResource{}
	_ resource.ResourceWithConfigure   = &ContactResource{}
	_ resource.ResourceWithImportState = &ContactResource{}
)

func NewContactResource() resource.Resource {
	return &ContactResource{}
}

// ContactResource defines the resource implementation.
type ContactResource struct {
	config *common.DnsimpleProviderConfig
}

// ContactResourceModel describes the resource data model.
type ContactResourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	AccountId        types.Int64  `tfsdk:"account_id"`
	Label            types.String `tfsdk:"label"`
	FirstName        types.String `tfsdk:"first_name"`
	LastName         types.String `tfsdk:"last_name"`
	OrganizationName types.String `tfsdk:"organization_name"`
	JobTitle         types.String `tfsdk:"job_title"`
	Address1         types.String `tfsdk:"address1"`
	Address2         types.String `tfsdk:"address2"`
	City             types.String `tfsdk:"city"`
	StateProvince    types.String `tfsdk:"state_province"`
	PostalCode       types.String `tfsdk:"postal_code"`
	Country          types.String `tfsdk:"country"`
	Phone            types.String `tfsdk:"phone"`
	PhoneNormalized  types.String `tfsdk:"phone_normalized"`
	Fax              types.String `tfsdk:"fax"`
	FaxNormalized    types.String `tfsdk:"fax_normalized"`
	Email            types.String `tfsdk:"email"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
}

func (r *ContactResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_contact"
}

func (r *ContactResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple contact resource",
		Attributes: map[string]schema.Attribute{
			"id": common.IDInt64Attribute(),
			"account_id": schema.Int64Attribute{
				Computed: true,
			},
			"label": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{modifiers.StringDefaultValue("")},
			},
			"first_name": schema.StringAttribute{
				Required: true,
			},
			"last_name": schema.StringAttribute{
				Required: true,
			},
			"organization_name": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{modifiers.StringDefaultValue("")},
			},
			"job_title": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{modifiers.StringDefaultValue("")},
			},
			"address1": schema.StringAttribute{
				Required: true,
			},
			"address2": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{modifiers.StringDefaultValue("")},
			},
			"city": schema.StringAttribute{
				Required: true,
			},
			"state_province": schema.StringAttribute{
				Required: true,
			},
			"postal_code": schema.StringAttribute{
				Required: true,
			},
			"country": schema.StringAttribute{
				Required: true,
			},
			"phone": schema.StringAttribute{
				Required: true,
			},
			"phone_normalized": schema.StringAttribute{
				Computed: true,
			},
			"fax": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				PlanModifiers: []planmodifier.String{modifiers.StringDefaultValue("")},
			},
			"fax_normalized": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Required: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ContactResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ContactResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ContactResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	contactAttributes := dnsimple.Contact{
		Label:         data.Label.ValueString(),
		FirstName:     data.FirstName.ValueString(),
		LastName:      data.LastName.ValueString(),
		JobTitle:      data.JobTitle.ValueString(),
		Organization:  data.OrganizationName.ValueString(),
		Address1:      data.Address1.ValueString(),
		Address2:      data.Address2.ValueString(),
		City:          data.City.ValueString(),
		StateProvince: data.StateProvince.ValueString(),
		PostalCode:    data.PostalCode.ValueString(),
		Country:       data.Country.ValueString(),
		Phone:         data.Phone.ValueString(),
		Fax:           data.Fax.ValueString(),
		Email:         data.Email.ValueString(),
	}

	response, err := r.config.Client.Contacts.CreateContact(ctx, r.config.AccountID, contactAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(attributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple Contact",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ContactResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, err := r.config.Client.Contacts.GetContact(ctx, r.config.AccountID, data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Contact: %d", data.Id.ValueInt64()),
			err.Error(),
		)
	}

	r.updateModelFromAPIResponse(response.Data, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ContactResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	contactAttributes := dnsimple.Contact{
		Label:         data.Label.ValueString(),
		FirstName:     data.FirstName.ValueString(),
		LastName:      data.LastName.ValueString(),
		JobTitle:      data.JobTitle.ValueString(),
		Organization:  data.OrganizationName.ValueString(),
		Address1:      data.Address1.ValueString(),
		Address2:      data.Address2.ValueString(),
		City:          data.City.ValueString(),
		StateProvince: data.StateProvince.ValueString(),
		PostalCode:    data.PostalCode.ValueString(),
		Country:       data.Country.ValueString(),
		Phone:         data.Phone.ValueString(),
		Fax:           data.Fax.ValueString(),
		Email:         data.Email.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("DNSimple Contact updateContactAttributes: %+v", contactAttributes))

	response, err := r.config.Client.Contacts.UpdateContact(
		ctx,
		r.config.AccountID,
		data.Id.ValueInt64(),
		contactAttributes,
	)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(attributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		fmt.Printf("error: %+v", err)

		resp.Diagnostics.AddError(
			"failed to update DNSimple Contact",
			err.Error(),
		)
		return
	}

	r.updateModelFromAPIResponse(response.Data, data)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ContactResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ContactResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting DNSimple Contact: %s, %s", data.Label, data.Id))

	_, err := r.config.Client.Contacts.DeleteContact(ctx, r.config.AccountID, data.Id.ValueInt64())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to delete DNSimple Contact: %d", data.Id.ValueInt64()),
			err.Error(),
		)
	}
}

func (r *ContactResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("resource import invalid ID", fmt.Sprintf("failed to parse contact ID (%s) as integer", req.ID))
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func (r *ContactResource) updateModelFromAPIResponse(contact *dnsimple.Contact, data *ContactResourceModel) {
	data.Id = types.Int64Value(contact.ID)
	data.AccountId = types.Int64Value(contact.AccountID)
	data.Label = types.StringValue(contact.Label)
	data.FirstName = types.StringValue(contact.FirstName)
	data.LastName = types.StringValue(contact.LastName)
	data.OrganizationName = types.StringValue(contact.Organization)
	data.JobTitle = types.StringValue(contact.JobTitle)
	data.Address1 = types.StringValue(contact.Address1)
	data.Address2 = types.StringValue(contact.Address2)
	data.City = types.StringValue(contact.City)
	data.StateProvince = types.StringValue(contact.StateProvince)
	data.PostalCode = types.StringValue(contact.PostalCode)
	data.Country = types.StringValue(contact.Country)
	data.PhoneNormalized = types.StringValue(contact.Phone)
	data.FaxNormalized = types.StringValue(contact.Fax)
	data.Email = types.StringValue(contact.Email)
	data.CreatedAt = types.StringValue(contact.CreatedAt)
	data.UpdatedAt = types.StringValue(contact.UpdatedAt)
}
