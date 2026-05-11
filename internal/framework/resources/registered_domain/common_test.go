package registered_domain

import (
	"context"
	"testing"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

func TestUpdateModelFromAPIResponsePartialCreate(t *testing.T) {
	t.Parallel()

	domain := &dnsimple.Domain{
		ID:           123,
		AccountID:    456,
		Name:         "example.com.br",
		UnicodeName:  "example.com.br",
		State:        "hosted",
		AutoRenew:    false,
		PrivateWhois: false,
		ExpiresAt:    "2027-05-11T00:00:00Z",
		Trustee:      true,
	}

	domainRegistration := &dnsimple.DomainRegistration{
		ID:     789,
		State:  "registering",
		Period: 1,
	}

	t.Run("sets unknown auto_renew_enabled to known value", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.False(t, data.AutoRenewEnabled.IsUnknown(), "auto_renew_enabled should not be unknown")
		assert.Equal(t, false, data.AutoRenewEnabled.ValueBool())
	})

	t.Run("sets unknown whois_privacy_enabled to known value", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.False(t, data.WhoisPrivacyEnabled.IsUnknown(), "whois_privacy_enabled should not be unknown")
		assert.Equal(t, false, data.WhoisPrivacyEnabled.ValueBool())
	})

	t.Run("sets unknown trustee to known value", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.False(t, data.Trustee.IsUnknown(), "trustee should not be unknown")
		assert.Equal(t, true, data.Trustee.ValueBool())
	})

	t.Run("preserves user-specified values", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolValue(true),
			WhoisPrivacyEnabled: types.BoolValue(true),
			Trustee:             types.BoolValue(false),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.Equal(t, true, data.AutoRenewEnabled.ValueBool())
		assert.Equal(t, true, data.WhoisPrivacyEnabled.ValueBool())
		assert.Equal(t, false, data.Trustee.ValueBool())
	})

	t.Run("sets registrant_change to null", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
			RegistrantChange:    types.ObjectUnknown(common.RegistrantChangeAttrType),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.True(t, data.RegistrantChange.IsNull(), "registrant_change should be null")
		assert.False(t, data.RegistrantChange.IsUnknown(), "registrant_change should not be unknown")
	})

	t.Run("sets dnssec and transfer lock to false", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)
		assert.Equal(t, false, data.DNSSECEnabled.ValueBool())
		assert.Equal(t, false, data.TransferLockEnabled.ValueBool())
	})

	t.Run("all values are known after partial create", func(t *testing.T) {
		t.Parallel()
		r := &RegisteredDomainResource{}
		data := &RegisteredDomainResourceModel{
			AutoRenewEnabled:    types.BoolUnknown(),
			WhoisPrivacyEnabled: types.BoolUnknown(),
			Trustee:             types.BoolUnknown(),
			RegistrantChange:    types.ObjectUnknown(common.RegistrantChangeAttrType),
		}

		diags := r.updateModelFromAPIResponsePartialCreate(context.Background(), data, domainRegistration, domain)
		assert.Nil(t, diags)

		assert.False(t, data.AutoRenewEnabled.IsUnknown(), "auto_renew_enabled should be known")
		assert.False(t, data.WhoisPrivacyEnabled.IsUnknown(), "whois_privacy_enabled should be known")
		assert.False(t, data.Trustee.IsUnknown(), "trustee should be known")
		assert.False(t, data.DNSSECEnabled.IsUnknown(), "dnssec_enabled should be known")
		assert.False(t, data.TransferLockEnabled.IsUnknown(), "transfer_lock_enabled should be known")
		assert.False(t, data.RegistrantChange.IsUnknown(), "registrant_change should be known")
		assert.False(t, data.DomainRegistration.IsUnknown(), "domain_registration should be known")
		assert.False(t, data.State.IsUnknown(), "state should be known")
		assert.False(t, data.Name.IsUnknown(), "name should be known")
		assert.False(t, data.UnicodeName.IsUnknown(), "unicode_name should be known")
		assert.False(t, data.AccountId.IsUnknown(), "account_id should be known")
		assert.False(t, data.ExpiresAt.IsUnknown(), "expires_at should be known")
		assert.False(t, data.Id.IsUnknown(), "id should be known")
	})
}
