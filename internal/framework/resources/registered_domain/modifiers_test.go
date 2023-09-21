package registered_domain_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources/registered_domain"
)

func TestDomainRegistrationStateModifier(t *testing.T) {
	t.Parallel()

	getDomainRegistration := func(ctx context.Context, value types.Object) *common.DomainRegistration {
		domainRegistration := &common.DomainRegistration{}
		if diags := value.As(ctx, domainRegistration, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags.HasError() {
			t.Fatal(diags)
			return nil
		}

		return domainRegistration
	}

	domainRegistrationToObject := func(ctx context.Context, domainRegistration *common.DomainRegistration) types.Object {
		obj, diags := types.ObjectValueFrom(ctx, common.DomainRegistrationAttrType, domainRegistration)
		if diags.HasError() {
			t.Fatal(diags)
			return obj
		}

		return obj
	}

	ctx := context.Background()

	type testCase struct {
		skipDomainRegistration bool
		plannedValue           types.Object
		currentValue           types.Object
		expectedValue          types.Object
		expectError            bool
	}
	tests := map[string]testCase{
		"imported has no plan and state and is unkown": {
			skipDomainRegistration: true,
			plannedValue:           types.ObjectUnknown(common.DomainRegistrationAttrType),
			currentValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			expectedValue:          types.ObjectNull(common.DomainRegistrationAttrType),
		},
		"imported has no plan and is null": {
			skipDomainRegistration: true,
			plannedValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			currentValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			expectedValue:          types.ObjectNull(common.DomainRegistrationAttrType),
		},
		"not registered has plan and state": {
			skipDomainRegistration: false,
			plannedValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateNew)}),
			currentValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateNew)}),
			expectedValue:          domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateRegistered)}),
		},
		"state null but plan is known": {
			skipDomainRegistration: false,
			plannedValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateNew)}),
			currentValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			expectedValue:          domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateRegistered)}),
		},
		"null planned": {
			skipDomainRegistration: false,
			plannedValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			currentValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateFailed)}),
			expectedValue:          types.ObjectNull(common.DomainRegistrationAttrType),
		},
		"state failing": {
			skipDomainRegistration: false,
			plannedValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateFailed)}),
			currentValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateFailed)}),
			expectedValue:          domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateFailed)}),
		},
		"state cancelling": {
			skipDomainRegistration: false,
			plannedValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelling)}),
			currentValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelling)}),
			expectedValue:          domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelling)}),
		},
		"state cancelled": {
			skipDomainRegistration: false,
			plannedValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelled)}),
			currentValue:           domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelled)}),
			expectedValue:          domainRegistrationToObject(ctx, &common.DomainRegistration{State: types.StringValue(consts.DomainStateCancelled)}),
		},
		"on create": {
			skipDomainRegistration: false,
			plannedValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			currentValue:           types.ObjectNull(common.DomainRegistrationAttrType),
			expectedValue:          types.ObjectNull(common.DomainRegistrationAttrType),
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			request := planmodifier.ObjectRequest{
				Path:       path.Root("test"),
				PlanValue:  test.plannedValue,
				StateValue: test.currentValue,
			}

			if test.skipDomainRegistration {
				request.Private.SetKey(ctx, "skipDomainRegistration", []byte("true"))
			}

			response := planmodifier.ObjectResponse{
				PlanValue: request.PlanValue,
			}
			registered_domain.DomainRegistrationState().PlanModifyObject(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			got := getDomainRegistration(ctx, response.PlanValue)
			want := getDomainRegistration(ctx, test.expectedValue)

			assert.Equal(t, want, got)
		})
	}
}

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
			registered_domain.RegistrantChangeState().PlanModifyString(ctx, request, &response)

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
