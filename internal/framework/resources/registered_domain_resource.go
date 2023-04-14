package resources

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &RegisteredDomainResource{}
	_ resource.ResourceWithConfigure   = &RegisteredDomainResource{}
	_ resource.ResourceWithImportState = &RegisteredDomainResource{}
)

func NewRegisteredDomainResource() resource.Resource {
	return &RegisteredDomainResource{}
}

// RegisteredDomainResource defines the resource implementation.
type RegisteredDomainResource struct {
	config *common.DnsimpleProviderConfig
}

// DomainResourceModel describes the resource data model.
type RegisteredDomainResourceModel struct {
	Name                types.String `tfsdk:"name"`
	AccountId           types.Int64  `tfsdk:"account_id"`
	UnicodeName         types.String `tfsdk:"unicode_name"`
	State               types.String `tfsdk:"state"`
	AutoRenewEnabled    types.Bool   `tfsdk:"auto_renew_enabled"`
	WhoisPrivacyEnabled types.Bool   `tfsdk:"whois_privacy_enabled"`
	DNSSECEnabled       types.Bool   `tfsdk:"dnssec_enabled"`
	ContactId           types.Int64  `tfsdk:"contact_id"`
	ExtendedAttributes  types.Map    `tfsdk:"extended_attributes"`
	PremiumPrice        types.String `tfsdk:"premium_price"`
	DomainRegistration  types.Object `tfsdk:"domain_registration"`
	Timeouts            types.Object `tfsdk:"timeouts"`
	Id                  types.Int64  `tfsdk:"id"`
}

var DomainRegistrationAttrType = map[string]attr.Type{
	"period": types.Int64Type,
	"state":  types.StringType,
	"id":     types.Int64Type,
}

type DomainRegistration struct {
	Period types.Int64  `tfsdk:"period"`
	State  types.String `tfsdk:"state"`
	Id     types.Int64  `tfsdk:"id"`
}

func (r *RegisteredDomainResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registered_domain"
}

func (r *RegisteredDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "DNSimple domain resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					validators.DomainName{},
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"account_id": schema.Int64Attribute{
				Computed: true,
			},
			"unicode_name": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"auto_renew_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"whois_privacy_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"dnssec_enabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"contact_id": schema.Int64Attribute{
				Required: true,
			},
			"extended_attributes": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"premium_price": schema.StringAttribute{
				Optional: true,
			},
			"domain_registration": schema.SingleNestedAttribute{
				Description: "The domain registration details.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"period": schema.Int64Attribute{
						Computed: true,
					},
					"state": schema.StringAttribute{
						Computed: true,
					},
					"id": common.IDInt64Attribute(),
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"timeouts": schema.SingleNestedAttribute{
				MarkdownDescription: "Timeouts for operations, given as a parsable string as in `10m` or `30s`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						Optional:    true,
						Description: "Create timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("10m"),
						},
					},
					"update": schema.StringAttribute{
						Optional:    true,
						Description: "Update timeout.",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("30s"),
						},
					},
					"delete": schema.StringAttribute{
						Optional:    true,
						Description: "Delete timeout (currently unused).",
						Validators: []validator.String{
							validators.Duration{},
						},
						PlanModifiers: []planmodifier.String{
							modifiers.StringDefaultValue("30s"),
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"id": common.IDInt64Attribute(),
		},
	}
}

func (r *RegisteredDomainResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RegisteredDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainAttributes := dnsimple.RegisterDomainInput{
		RegistrantID: int(data.ContactId.ValueInt64()),
	}

	if !data.AutoRenewEnabled.IsNull() {
		domainAttributes.EnableAutoRenewal = data.AutoRenewEnabled.ValueBool()
	}

	if !data.WhoisPrivacyEnabled.IsNull() {
		domainAttributes.EnableWhoisPrivacy = data.WhoisPrivacyEnabled.ValueBool()
	}

	if !data.PremiumPrice.IsNull() {
		domainAttributes.PremiumPrice = data.PremiumPrice.ValueString()
	}

	if !data.ExtendedAttributes.IsNull() {
		extendedAttrs := make(map[string]string)
		diags := data.ExtendedAttributes.ElementsAs(ctx, &extendedAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		domainAttributes.ExtendedAttributes = extendedAttrs
	}

	lowerCaseDomainName := strings.ToLower(data.Name.ValueString())
	response, err := r.config.Client.Registrar.RegisterDomain(ctx, r.config.AccountID, lowerCaseDomainName, &domainAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(attributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to register DNSimple Domain",
			err.Error(),
		)
		return
	}

	if response.Data.State == consts.DomainStateHosted {
		resp.Diagnostics.AddError(
			"failed to register DNSimple Domain",
			"Domain added to DNSimple as hosted, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
		)
		return
	}

	if response.Data.State != consts.DomainStateRegistered {
		timeouts, diags := getTimeouts(ctx, data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		utils.RetryWithTimeout(ctx, func() error {

			domainRegistration, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, lowerCaseDomainName, strconv.Itoa(int(response.Data.ID)))

			if err != nil {
				return err
			}

			if domainRegistration.Data.State == consts.DomainStateHosted {
				return fmt.Errorf("domain added to DNSimple as hosted, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com")
			}

			if domainRegistration.Data.State != consts.DomainStateRegistered {
				tflog.Info(ctx, fmt.Sprintf("[RETRYING] Domain registration is not complete, current state: %s", domainRegistration.Data.State))

				return fmt.Errorf("domain registration is not complete, current state: %s", domainRegistration.Data.State)
			}

			return nil
		}, timeouts.CreateDuration(), 20*time.Second)
	}

	if !data.DNSSECEnabled.IsNull() && data.DNSSECEnabled.ValueBool() {

		diags := r.setDNSSEC(ctx, data)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, lowerCaseDomainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", lowerCaseDomainName),
			err.Error(),
		)
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, lowerCaseDomainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", lowerCaseDomainName),
			err.Error(),
		)
	}

	diags := r.updateModelFromAPIResponse(ctx, data, response.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(*diags...)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegisteredDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	domainRegistration := &DomainRegistration{}
	if diags := data.DomainRegistration.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), domainRegistrationId)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", data.Name.ValueString(), domainRegistration.Id.ValueInt64()),
			err.Error(),
		)
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", data.Name.ValueString()),
			err.Error(),
		)
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, data.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", data.Name.ValueString()),
			err.Error(),
		)
	}

	diags := r.updateModelFromAPIResponse(ctx, data, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(*diags...)
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RegisteredDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var (
		configData *RegisteredDomainResourceModel
		planData   *RegisteredDomainResourceModel
		stateData  *RegisteredDomainResourceModel
	)

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &configData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planData.ContactId.ValueInt64() != stateData.ContactId.ValueInt64() {
		resp.Diagnostics.AddError(
			fmt.Sprintf("contact_id change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
			"contact_id change not supported by the DNSimple API",
		)
		return
	}

	if !planData.ExtendedAttributes.Equal(stateData.ExtendedAttributes) {
		resp.Diagnostics.AddError(
			fmt.Sprintf("extended_attributes change not supported: %s, %d", planData.Name.ValueString(), planData.Id.ValueInt64()),
			"extended_attributes change not supported by the DNSimple API",
		)
		return
	}

	if planData.AutoRenewEnabled.ValueBool() != stateData.AutoRenewEnabled.ValueBool() {

		diags := r.setAutoRenewal(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if planData.WhoisPrivacyEnabled.ValueBool() != stateData.WhoisPrivacyEnabled.ValueBool() {

		diags := r.setWhoisPrivacy(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if planData.DNSSECEnabled.ValueBool() != stateData.DNSSECEnabled.ValueBool() {

		diags := r.setDNSSEC(ctx, planData)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	domainRegistration := &DomainRegistration{}
	if diags := planData.DomainRegistration.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	domainRegistrationId := strconv.Itoa(int(domainRegistration.Id.ValueInt64()))
	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, planData.Name.ValueString(), domainRegistrationId)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain Registration: %s, %d", planData.Name.ValueString(), domainRegistration.Id.ValueInt64()),
			err.Error(),
		)
		return
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, planData.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain: %s", planData.Name.ValueString()),
			err.Error(),
		)
		return
	}

	dnssecResponse, err := r.config.Client.Domains.GetDnssec(ctx, r.config.AccountID, planData.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to read DNSimple Domain DNSSEC status: %s", planData.Name.ValueString()),
			err.Error(),
		)
		return
	}

	diags := r.updateModelFromAPIResponse(ctx, planData, domainRegistrationResponse.Data, domainResponse.Data, dnssecResponse.Data)
	if diags != nil && diags.HasError() {
		resp.Diagnostics.Append(*diags...)
		return
	}

	// Save updated planData into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *RegisteredDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RegisteredDomainResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Warn(ctx, fmt.Sprintf("Removing DNSimple Registered Domain from Terraform state only: %s, %s", data.Name, data.Id))
}

func (r *RegisteredDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "_")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("resource import invalid ID", fmt.Sprintf("wrong format of import ID (%s), use: '<domain-name>_<domain-registration-id>'", req.ID))
	}
	domainName := parts[0]
	domainRegistrationID := parts[1]

	domainRegistrationResponse, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, domainName, domainRegistrationID)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("failed to find DNSimple Domain Registration ID: %s", domainRegistrationID),
			err.Error(),
		)
		return
	}

	domainResponse, err := r.config.Client.Domains.GetDomain(ctx, r.config.AccountID, domainName)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("unexpected error when trying to find DNSimple Domain ID: %s", domainName),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), domainResponse.Data.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), domainResponse.Data.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_registration").AtName("id"), domainRegistrationResponse.Data.ID)...)
}

func (r *RegisteredDomainResource) setAutoRenewal(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting auto_renew_enabled to %t", data.AutoRenewEnabled.ValueBool()))

	if data.AutoRenewEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableDomainAutoRenewal(ctx, r.config.AccountID, data.Name.ValueString())
		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain Auto Renewal: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableDomainAutoRenewal(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain Auto Renewal: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) setWhoisPrivacy(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting whois_privacy_enabled to %t", data.WhoisPrivacyEnabled.ValueBool()))

	if data.WhoisPrivacyEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableWhoisPrivacy(ctx, r.config.AccountID, data.Name.ValueString())
		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain Whois Privacy: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableWhoisPrivacy(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain Whois Privacy: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) setDNSSEC(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting dnssec_enabled to %t", data.DNSSECEnabled.ValueBool()))

	if data.DNSSECEnabled.ValueBool() {
		_, err := r.config.Client.Domains.EnableDnssec(ctx, r.config.AccountID, data.Name.ValueString())

		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain DNSSEC: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Domains.DisableDnssec(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain DNSSEC: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) updateModelFromAPIResponse(ctx context.Context, data *RegisteredDomainResourceModel, domainRegistration *dnsimple.DomainRegistration, domain *dnsimple.Domain, dnssec *dnsimple.Dnssec) *diag.Diagnostics {
	domainRegistrationData := DomainRegistration{
		Id:     types.Int64Value(domainRegistration.ID),
		Period: types.Int64Value(int64(domainRegistration.Period)),
		State:  types.StringValue(domainRegistration.State),
	}

	obj, diags := types.ObjectValueFrom(ctx, DomainRegistrationAttrType, domainRegistrationData)

	if diags.HasError() {
		return &diags
	}

	data.DomainRegistration = obj

	data.Id = types.Int64Value(domain.ID)
	data.AutoRenewEnabled = types.BoolValue(domain.AutoRenew)
	data.WhoisPrivacyEnabled = types.BoolValue(domain.PrivateWhois)
	data.State = types.StringValue(domain.State)
	data.UnicodeName = types.StringValue(domain.UnicodeName)
	data.AccountId = types.Int64Value(domain.AccountID)
	data.ContactId = types.Int64Value(domain.RegistrantID)
	data.Name = types.StringValue(domain.Name)

	data.DNSSECEnabled = types.BoolValue(dnssec.Enabled)

	return nil
}

func getTimeouts(ctx context.Context, model *RegisteredDomainResourceModel) (*common.Timeouts, diag.Diagnostics) {
	timeouts := &common.Timeouts{}
	diags := model.Timeouts.As(ctx, timeouts, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return timeouts, diags
}
