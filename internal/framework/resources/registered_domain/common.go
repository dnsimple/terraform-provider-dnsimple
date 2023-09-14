package registered_domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

const (
	RegistrationConverged          = "registration_converged"
	RegistrationConvergenceTimeout = "registration_converged_timeout"
	RegistrationFailed             = "registration_failed"

	RegistrantChangeConverged          = "registrant_change_converged"
	RegistrantChangeConvergenceTimeout = "registrant_change_converged_timeout"
	RegistrantChangeFailed             = "registrant_change_failed"
)

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

func (r *RegisteredDomainResource) setTransferLock(ctx context.Context, data *RegisteredDomainResourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	tflog.Debug(ctx, fmt.Sprintf("setting transfer_lock_enabled to %t", data.TransferLockEnabled.ValueBool()))

	if data.TransferLockEnabled.ValueBool() {
		_, err := r.config.Client.Registrar.EnableDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())

		if err != nil {
			diagnostics.AddError(
				fmt.Sprintf("failed to enable DNSimple Domain transfer lock: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
				err.Error(),
			)
		}
		return diagnostics
	}

	_, err := r.config.Client.Registrar.DisableDomainTransferLock(ctx, r.config.AccountID, data.Name.ValueString())
	if err != nil {
		diagnostics.AddError(
			fmt.Sprintf("failed to disable DNSimple Domain transfer lock: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
			err.Error(),
		)
	}

	return diagnostics
}

func (r *RegisteredDomainResource) updateModelFromAPIResponse(ctx context.Context, data *RegisteredDomainResourceModel, domainRegistration *dnsimple.DomainRegistration, domain *dnsimple.Domain, dnssec *dnsimple.Dnssec, transferLock *dnsimple.DomainTransferLock) diag.Diagnostics {
	if domainRegistration != nil {
		domainRegistrationObject, diags := r.domainRegistrationAPIResponseToObject(ctx, domainRegistration)

		if diags.HasError() {
			return diags
		}

		data.DomainRegistration = domainRegistrationObject
	}

	if domain != nil {
		data.Id = types.Int64Value(domain.ID)
		data.AutoRenewEnabled = types.BoolValue(domain.AutoRenew)
		data.WhoisPrivacyEnabled = types.BoolValue(domain.PrivateWhois)
		data.State = types.StringValue(domain.State)
		data.UnicodeName = types.StringValue(domain.UnicodeName)
		data.AccountId = types.Int64Value(domain.AccountID)
		data.ExpiresAt = types.StringValue(domain.ExpiresAt)
		data.Name = types.StringValue(domain.Name)
	}

	if dnssec != nil {
		data.DNSSECEnabled = types.BoolValue(dnssec.Enabled)
	}

	if transferLock != nil {
		data.TransferLockEnabled = types.BoolValue(transferLock.Enabled)
	}

	return nil
}

func (r *RegisteredDomainResource) updateModelFromAPIResponsePartialCreate(ctx context.Context, data *RegisteredDomainResourceModel, domainRegistration *dnsimple.DomainRegistration, domain *dnsimple.Domain) *diag.Diagnostics {
	domainRegistrationObject, diags := r.domainRegistrationAPIResponseToObject(ctx, domainRegistration)

	if diags.HasError() {
		return &diags
	}

	data.DomainRegistration = domainRegistrationObject

	if domain != nil {
		data.Id = types.Int64Value(domain.ID)

		if data.AutoRenewEnabled.IsNull() {
			data.AutoRenewEnabled = types.BoolValue(domain.AutoRenew)
		}

		if data.WhoisPrivacyEnabled.IsNull() {
			data.WhoisPrivacyEnabled = types.BoolValue(domain.PrivateWhois)
		}

		data.State = types.StringValue(domain.State)
		data.UnicodeName = types.StringValue(domain.UnicodeName)
		data.AccountId = types.Int64Value(domain.AccountID)
		data.ExpiresAt = types.StringValue(domain.ExpiresAt)
		data.Name = types.StringValue(domain.Name)
	}

	data.DNSSECEnabled = types.BoolValue(false)
	data.TransferLockEnabled = types.BoolValue(false)

	return nil
}

func (r *RegisteredDomainResource) domainRegistrationAPIResponseToObject(ctx context.Context, domainRegistration *dnsimple.DomainRegistration) (basetypes.ObjectValue, diag.Diagnostics) {
	domainRegistrationData := common.DomainRegistration{
		Id:     types.Int64Value(domainRegistration.ID),
		Period: types.Int64Value(int64(domainRegistration.Period)),
		State:  types.StringValue(domainRegistration.State),
	}

	return types.ObjectValueFrom(ctx, common.DomainRegistrationAttrType, domainRegistrationData)
}

func (r *RegisteredDomainResource) registrantChangeAPIResponseToObject(ctx context.Context, registrantChange *dnsimple.RegistrantChange) (basetypes.ObjectValue, diag.Diagnostics) {
	elements, diags := types.MapValueFrom(ctx, types.StringType, registrantChange.ExtendedAttributes)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	registrantChangeData := common.RegistrantChange{
		Id:                  types.Int64Value(int64(registrantChange.Id)),
		AccountId:           types.Int64Value(int64(registrantChange.AccountId)),
		ContactId:           types.Int64Value(int64(registrantChange.ContactId)),
		DomainId:            types.StringValue(fmt.Sprintf("%d", registrantChange.DomainId)),
		State:               types.StringValue(registrantChange.State),
		ExtendedAttributes:  elements,
		RegistryOwnerChange: types.BoolValue(registrantChange.RegistryOwnerChange),
		IrtLockLiftedBy:     types.StringValue(registrantChange.IrtLockLiftedBy),
	}

	return types.ObjectValueFrom(ctx, common.RegistrantChangeAttrType, registrantChangeData)
}

func getTimeouts(ctx context.Context, model *RegisteredDomainResourceModel) (*common.Timeouts, diag.Diagnostics) {
	timeouts := &common.Timeouts{}
	diags := model.Timeouts.As(ctx, timeouts, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return timeouts, diags
}

func getDomainRegistration(ctx context.Context, data *RegisteredDomainResourceModel) (*common.DomainRegistration, diag.Diagnostics) {
	domainRegistration := &common.DomainRegistration{}
	diags := data.DomainRegistration.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return domainRegistration, diags
}

func getRegistrantChange(ctx context.Context, data *RegisteredDomainResourceModel) (*common.RegistrantChange, diag.Diagnostics) {
	registrantChange := &common.RegistrantChange{}
	diags := data.RegistrantChange.As(ctx, registrantChange, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return registrantChange, diags
}

func tryToConvergeRegistration(ctx context.Context, data *RegisteredDomainResourceModel, diagnostics *diag.Diagnostics, r *RegisteredDomainResource, registrationID string) (string, error) {
	timeouts, diags := getTimeouts(ctx, data)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return RegistrationFailed, nil
	}

	err := utils.RetryWithTimeout(ctx, func() (error, bool) {
		domainRegistration, err := r.config.Client.Registrar.GetDomainRegistration(ctx, r.config.AccountID, data.Name.ValueString(), registrationID)

		if err != nil {
			return err, false
		}

		if domainRegistration.Data.State == consts.DomainStateFailed {
			diagnostics.AddError(
				fmt.Sprintf("failed to register DNSimple Domain: %s", data.Name.ValueString()),
				"domain registration failed, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if domainRegistration.Data.State == consts.DomainStateCancelling || domainRegistration.Data.State == consts.DomainStateCancelled {
			diagnostics.AddError(
				fmt.Sprintf("failed to register DNSimple Domain: %s", data.Name.ValueString()),
				"domain registration was cancelled, please investigate why this happened. If you need assistance, please contact support at support@dnsimple.com",
			)
			return nil, true
		}

		if domainRegistration.Data.State != consts.DomainStateRegistered {
			tflog.Info(ctx, fmt.Sprintf("[RETRYING] Domain registration is not complete, current state: %s", domainRegistration.Data.State))

			return fmt.Errorf("domain registration is not complete, current state: %s. You can try to run terraform again to try and converge the domain registration", domainRegistration.Data.State), false
		}

		return nil, false
	}, timeouts.CreateDuration(), 20*time.Second)

	if diagnostics.HasError() {
		// If we have diagnostic errors, we suspended the retry loop because the domain is in a bad state, and cannot converge.
		return RegistrationFailed, nil
	}

	if err != nil {
		// If we have an error, it means the retry loop timed out, and we cannot converge during this run.
		return RegistrationConvergenceTimeout, err
	}

	return RegistrationConverged, nil
}

func tryToConvergeRegistrantChange(ctx context.Context, data *RegisteredDomainResourceModel, diagnostics *diag.Diagnostics, r *RegisteredDomainResource, registrantChangeId int) (string, error) {
	timeouts, diags := getTimeouts(ctx, data)
	diagnostics.Append(diags...)
	if diagnostics.HasError() {
		return RegistrantChangeFailed, nil
	}

	err := utils.RetryWithTimeout(ctx, func() (error, bool) {
		registrantChangeResponse, err := r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, registrantChangeId)

		if err != nil {
			return err, false
		}

		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelled || registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelling {
			diagnostics.AddError(
				fmt.Sprintf("failed to change registrant for DNSimple Domain: %s, registrant change id: %d", data.Name.ValueString(), registrantChangeId),
				"Registrant change was cancelled, please investigate why this happened. You can refer to our support article https://support.dnsimple.com/articles/changing-domain-contact/ to get started and if you need assistance, please contact support at support@dnsimple.com.",
			)
			return nil, true
		}

		if registrantChangeResponse.Data.State != consts.DomainStateRegistered {
			tflog.Info(ctx, fmt.Sprintf("[RETRYING] Registrant change is not complete, current state: %s", registrantChangeResponse.Data.State))

			return fmt.Errorf("registrant change is not complete, current state: %s. You can try to run terraform again to try and converge the registrant change", registrantChangeResponse.Data.State), false
		}

		return nil, false
	}, timeouts.CreateDuration(), 20*time.Second)

	if diagnostics.HasError() {
		// If we have diagnostic errors, we suspended the retry loop because the domain is in a bad state, and cannot converge.
		return RegistrantChangeFailed, nil
	}

	if err != nil {
		// If we have an error, it means the retry loop timed out, and we cannot converge during this run.
		return RegistrantChangeConvergenceTimeout, err
	}

	return RegistrantChangeConverged, nil
}

func createRegistrantChange(ctx context.Context, data *RegisteredDomainResourceModel, r *RegisteredDomainResource, resp *resource.UpdateResponse) {
	registrantChangeAttributes := dnsimple.CreateRegistrantChangeInput{
		DomainId:  fmt.Sprintf("%d", data.Id.ValueInt64()),
		ContactId: fmt.Sprintf("%d", data.ContactId.ValueInt64()),
	}

	if !data.ExtendedAttributes.IsNull() {
		extendedAttrs := make(map[string]string)
		diags := data.ExtendedAttributes.ElementsAs(ctx, &extendedAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		registrantChangeAttributes.ExtendedAttributes = extendedAttrs
	}

	registrantChangeResponse, err := r.config.Client.Registrar.CreateRegistrantChange(ctx, r.config.AccountID, &registrantChangeAttributes)

	if err != nil {
		var errorResponse *dnsimple.ErrorResponse
		if errors.As(err, &errorResponse) {
			resp.Diagnostics.Append(utils.AttributeErrorsToDiagnostics(errorResponse)...)
			return
		}

		resp.Diagnostics.AddError(
			"failed to create DNSimple registrant change",
			err.Error(),
		)
		return
	}

	if registrantChangeResponse.Data.State != consts.RegistrantChangeStateCompleted {
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelled || registrantChangeResponse.Data.State == consts.RegistrantChangeStateCancelling {
			resp.Diagnostics.AddError(
				fmt.Sprintf("failed to create DNSimple registrant change current state: %s", registrantChangeResponse.Data.State),
				"",
			)
			return
		}

		// Registrant change has been created, but is not yet completed
		if registrantChangeResponse.Data.State == consts.RegistrantChangeStateNew || registrantChangeResponse.Data.State == consts.RegistrantChangeStatePending {
			convergenceState, err := tryToConvergeRegistrantChange(ctx, data, &resp.Diagnostics, r, registrantChangeResponse.Data.Id)
			if convergenceState == RegistrantChangeFailed {
				// Response is already populated with the error we can safely return
				return
			}

			if convergenceState == RegistrantChangeConvergenceTimeout {
				// We attempted to converge on the registrant change, but the registrant change was not ready
				// user needs to run terraform again to try and converge the registrant change

				registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
				resp.Diagnostics.Append(diags...)

				// Update the data with the current registrant change
				data.RegistrantChange = registrantChangeObject

				// Save data into Terraform state
				resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

				// Exit with warning to prevent the state from being tainted
				resp.Diagnostics.AddWarning(
					"failed to converge on registrant change",
					err.Error(),
				)
				return
			}

			registrantChangeResponse, err = r.config.Client.Registrar.GetRegistrantChange(ctx, r.config.AccountID, int(data.Id.ValueInt64()))

			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("failed to read DNSimple Registrant Change for: %s, %d", data.Name.ValueString(), data.Id.ValueInt64()),
					err.Error(),
				)
				return
			}
		}
	}

	// Registrant change was successful, we can now update the resource model
	registrantChangeObject, diags := r.registrantChangeAPIResponseToObject(ctx, registrantChangeResponse.Data)
	resp.Diagnostics.Append(diags...)

	// Update the data with the current registrant change
	data.RegistrantChange = registrantChangeObject
}
