package registered_domain_contact_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources/registered_domain_contact"
)

func TestRegistrantChangeStateModifier(t *testing.T) {
	t.Parallel()

	type testCase struct {
		skipDomainRegistration bool
		plannedValue           types.String
		currentValue           types.String
		expectedValue          types.String
		expectError            bool
	}
	tests := map[string]testCase{
		"imported has no plan and state and is unkown": {
			plannedValue:  types.StringUnknown(),
			currentValue:  types.StringNull(),
			expectedValue: types.StringNull(),
		},
		"imported has no plan and is null": {
			plannedValue:  types.StringNull(),
			currentValue:  types.StringNull(),
			expectedValue: types.StringNull(),
		},
		"not completed has plan and state": {
			plannedValue:  types.StringValue(consts.RegistrantChangeStateNew),
			currentValue:  types.StringValue(consts.RegistrantChangeStateNew),
			expectedValue: types.StringValue(consts.RegistrantChangeStateCompleted),
		},
		"state null but plan is known": {
			plannedValue:  types.StringValue(consts.RegistrantChangeStateNew),
			currentValue:  types.StringNull(),
			expectedValue: types.StringValue(consts.RegistrantChangeStateCompleted),
		},
		"null planned": {
			plannedValue:  types.StringNull(),
			currentValue:  types.StringValue(consts.RegistrantChangeStateCancelled),
			expectedValue: types.StringNull(),
		},
		"state cancelled": {
			plannedValue:  types.StringValue(consts.RegistrantChangeStateCancelled),
			currentValue:  types.StringValue(consts.RegistrantChangeStateCancelled),
			expectedValue: types.StringValue(consts.RegistrantChangeStateCancelled),
		},
		"state cancelling": {
			plannedValue:  types.StringValue(consts.RegistrantChangeStateCancelling),
			currentValue:  types.StringValue(consts.RegistrantChangeStateCancelling),
			expectedValue: types.StringValue(consts.RegistrantChangeStateCancelling),
		},
		"on create": {
			plannedValue:  types.StringNull(),
			currentValue:  types.StringNull(),
			expectedValue: types.StringNull(),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			request := planmodifier.StringRequest{
				Path:       path.Root("test"),
				PlanValue:  test.plannedValue,
				StateValue: test.currentValue,
			}

			if test.skipDomainRegistration {
				request.Private.SetKey(ctx, "skipDomainRegistration", []byte("true"))
			}

			response := planmodifier.StringResponse{
				PlanValue: request.PlanValue,
			}
			registered_domain_contact.RegistrantChangeState().PlanModifyString(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			assert.Equal(t, test.expectedValue.ValueString(), response.PlanValue.ValueString())
		})
	}
}
