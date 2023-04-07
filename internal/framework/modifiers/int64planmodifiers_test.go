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

func TestInt64DefaultValue(t *testing.T) {
	t.Parallel()

	type testCase struct {
		plannedValue  types.Int64
		currentValue  types.Int64
		defaultValue  int64
		expectedValue types.Int64
		expectError   bool
	}
	tests := map[string]testCase{
		"default int64": {
			plannedValue:  types.Int64Null(),
			currentValue:  types.Int64Value(1),
			defaultValue:  1,
			expectedValue: types.Int64Value(1),
		},
		"default int64 on create": {
			plannedValue:  types.Int64Null(),
			currentValue:  types.Int64Null(),
			defaultValue:  1,
			expectedValue: types.Int64Value(1),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			request := planmodifier.Int64Request{
				Path:       path.Root("test"),
				PlanValue:  test.plannedValue,
				StateValue: test.currentValue,
			}
			response := planmodifier.Int64Response{
				PlanValue: request.PlanValue,
			}
			modifiers.Int64DefaultValue(test.defaultValue).PlanModifyInt64(ctx, request, &response)

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
