package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestRecordType_ValidateString(t *testing.T) {
	t.Parallel()

	type testCase struct {
		value       types.String
		expectError bool
		errorCount  int
	}

	tests := map[string]testCase{
		"valid uppercase A": {
			value:       types.StringValue("A"),
			expectError: false,
		},
		"valid uppercase AAAA": {
			value:       types.StringValue("AAAA"),
			expectError: false,
		},
		"valid uppercase CNAME": {
			value:       types.StringValue("CNAME"),
			expectError: false,
		},
		"valid uppercase MX": {
			value:       types.StringValue("MX"),
			expectError: false,
		},
		"valid uppercase TXT": {
			value:       types.StringValue("TXT"),
			expectError: false,
		},
		"valid uppercase SRV": {
			value:       types.StringValue("SRV"),
			expectError: false,
		},
		"valid uppercase NS": {
			value:       types.StringValue("NS"),
			expectError: false,
		},
		"valid uppercase PTR": {
			value:       types.StringValue("PTR"),
			expectError: false,
		},
		"valid uppercase SOA": {
			value:       types.StringValue("SOA"),
			expectError: false,
		},
		"invalid lowercase a": {
			value:       types.StringValue("a"),
			expectError: true,
			errorCount:  1,
		},
		"invalid lowercase mx": {
			value:       types.StringValue("mx"),
			expectError: true,
			errorCount:  1,
		},
		"invalid mixed case Aa": {
			value:       types.StringValue("Aa"),
			expectError: true,
			errorCount:  1,
		},
		"invalid mixed case Mx": {
			value:       types.StringValue("Mx"),
			expectError: true,
			errorCount:  1,
		},
		"invalid mixed case cNAME": {
			value:       types.StringValue("cNAME"),
			expectError: true,
			errorCount:  1,
		},
		"invalid mixed case TXt": {
			value:       types.StringValue("TXt"),
			expectError: true,
			errorCount:  1,
		},
		"null value": {
			value:       types.StringNull(),
			expectError: false,
		},
		"unknown value": {
			value:       types.StringUnknown(),
			expectError: false,
		},
		"empty string": {
			value:       types.StringValue(""),
			expectError: false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			request := validator.StringRequest{
				Path:        path.Root("type"),
				ConfigValue: test.value,
			}
			response := validator.StringResponse{}

			RecordType{}.ValidateString(ctx, request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}

			if test.expectError && test.errorCount > 0 {
				assert.Equal(t, test.errorCount, response.Diagnostics.ErrorsCount(), "expected %d error(s), got %d", test.errorCount, response.Diagnostics.ErrorsCount())

				// Verify the error message contains the expected text
				errorSummary := response.Diagnostics.Errors()[0].Summary()
				assert.Equal(t, "Record type must be uppercase", errorSummary)

				// Verify the error detail contains the suggestion
				errorDetail := response.Diagnostics.Errors()[0].Detail()
				assert.Contains(t, errorDetail, "UPPERCASE")
				assert.Contains(t, errorDetail, test.value.ValueString())
			}
		})
	}
}

func TestRecordType_Description(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	validator := RecordType{}

	description := validator.Description(ctx)
	assert.Equal(t, "record type must be specified in UPPERCASE", description)

	markdownDescription := validator.MarkdownDescription(ctx)
	assert.Equal(t, description, markdownDescription)
}
