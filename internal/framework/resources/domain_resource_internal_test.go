package resources

import (
	"testing"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// prevent_delete has no API counterpart and is only ever stored in state, so state
// written by a provider version that predates it carries no value. Refresh settles it to
// the schema default rather than leaving it null, which would otherwise surface as a
// null -> false diff on the first plan after upgrading. An explicit opt-in must survive.
func TestDomainResource_updateModelFromAPIResponse_PreventDelete(t *testing.T) {
	domain := &dnsimple.Domain{ID: 1, Name: "example.com"}

	for _, tt := range []struct {
		name  string
		prior types.Bool
		want  bool
	}{
		{name: "null prior state settles to the unprotected default", prior: types.BoolNull(), want: false},
		{name: "unknown prior state settles to the unprotected default", prior: types.BoolUnknown(), want: false},
		{name: "explicit false is preserved", prior: types.BoolValue(false), want: false},
		{name: "explicit opt-in to protection is preserved", prior: types.BoolValue(true), want: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := &DomainResource{}
			data := &DomainResourceModel{PreventDelete: tt.prior}

			r.updateModelFromAPIResponse(domain, data)

			assert.False(t, data.PreventDelete.IsNull(), "prevent_delete must never stay null")
			assert.False(t, data.PreventDelete.IsUnknown(), "prevent_delete must never stay unknown")
			assert.Equal(t, tt.want, data.PreventDelete.ValueBool())
		})
	}
}
