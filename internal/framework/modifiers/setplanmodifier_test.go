package modifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/modifiers"
)

func TestSetTrimSuffixValue(t *testing.T) {
	t.Parallel()

	serializeList := func(elements []string) types.Set {
		listValue, _ := types.SetValueFrom(context.Background(), types.StringType, elements)
		return listValue
	}

	type testCase struct {
		plannedValue  types.Set
		stateValue    types.Set
		configValue   types.Set
		expectedValue types.Set
		expectError   bool
	}
	tests := map[string]testCase{
		"trim `.` in string value ": {
			plannedValue:  serializeList([]string{"beta.alpha.", "gamma.omega"}),
			stateValue:    types.SetNull(types.StringType),
			configValue:   serializeList([]string{"beta.alpha.", "gamma.omega"}),
			expectedValue: serializeList([]string{"beta.alpha", "gamma.omega"}),
		},
		"when state is Null": {
			plannedValue:  serializeList([]string{"beta.alpha.", "gamma.omega."}),
			stateValue:    types.SetNull(types.StringType),
			configValue:   serializeList([]string{"beta.alpha.", "gamma.omega."}),
			expectedValue: serializeList([]string{"beta.alpha", "gamma.omega"}),
		},
		"when plan is Null": {
			plannedValue:  types.SetNull(types.StringType),
			stateValue:    types.SetNull(types.StringType),
			configValue:   serializeList([]string{"beta.alpha.", "gamma.omega."}),
			expectedValue: types.SetNull(types.StringType),
		},
		"when plan is Unknown": {
			plannedValue:  types.SetUnknown(types.StringType),
			stateValue:    types.SetNull(types.StringType),
			configValue:   serializeList([]string{"beta.alpha.", "gamma.omega."}),
			expectedValue: types.SetUnknown(types.StringType),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			request := planmodifier.SetRequest{
				Path:        path.Root("test"),
				PlanValue:   test.plannedValue,
				StateValue:  test.stateValue,
				ConfigValue: test.configValue,
			}
			response := planmodifier.SetResponse{
				PlanValue: request.PlanValue,
			}
			modifiers.SetTrimSuffixValue().PlanModifySet(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			assert.Equal(t, test.expectedValue, response.PlanValue)
		})
	}
}
