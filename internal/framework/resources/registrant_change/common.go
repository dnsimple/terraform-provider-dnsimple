package registrant_change

import (
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	RegistrantChangeConverged          = "registrant_change_converged"
	RegistrantChangeConvergenceTimeout = "registrant_change_converged_timeout"
	RegistrantChangeFailed             = "registrant_change_failed"
)

func (r *RegistrantChangeResource) updateModelFromAPIResponse(registrantChange *dnsimple.RegistrantChange, data *RegistrantChangeResourceModel) {
	data.Id = types.Int64Value(registrantChange.ID)
	data.AccountId = types.Int64Value(registrantChange.AccountID)
	data.ContactId = types.Int64Value(registrantChange.ContactID)
	data.DomainId = types.Int64Value(registrantChange.DomainID)
	data.State = types.StringValue(domain.State)
	data.ExtendedAttributes = []ExtendedAttribute{}
	data.RegistryOwnerChange = types.BoolValue(registrantChange.RegistryOwnerChange)
	data.IrtLockLiftedBy = types.DateValue(registrantChange.IrtLockLiftedBy)
}
