package registered_domain_contact

import (
	"context"
	"fmt"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

const (
	RegistrantChangeConverged          = "registrant_change_converged"
	RegistrantChangeConvergenceTimeout = "registrant_change_converged_timeout"
	RegistrantChangeFailed             = "registrant_change_failed"
)

func (r *RegisteredDomainContactResource) updateModelFromAPIResponse(registrantChange *dnsimple.RegistrantChange, data *RegisteredDomainContactResourceModel) {
	data.Id = types.Int64Value(int64(registrantChange.Id))
	data.AccountId = types.Int64Value(int64(registrantChange.AccountId))
	data.ContactId = types.Int64Value(int64(registrantChange.ContactId))
	data.State = types.StringValue(registrantChange.State)
	data.RegistryOwnerChange = types.BoolValue(registrantChange.RegistryOwnerChange)
	data.IrtLockLiftedBy = types.StringValue(registrantChange.IrtLockLiftedBy)
}

func getTimeouts(ctx context.Context, model *RegisteredDomainContactResourceModel) (*common.Timeouts, diag.Diagnostics) {
	timeouts := &common.Timeouts{}
	diags := model.Timeouts.As(ctx, timeouts, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	return timeouts, diags
}

func tryToConvergeRegistration(ctx context.Context, data *RegisteredDomainContactResourceModel, diagnostics *diag.Diagnostics, r *RegisteredDomainContactResource, registrantChangeId int) (string, error) {
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
				fmt.Sprintf("failed to change registrant for DNSimple Domain: %s, registrant change id: %d", data.DomainId.ValueString(), registrantChangeId),
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
