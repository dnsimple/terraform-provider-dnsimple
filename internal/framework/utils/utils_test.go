package utils_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func TestHasUnicodeChars(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{
			name: "empty string",
			s:    "",
			want: false,
		},
		{
			name: "ascii string",
			s:    "hello-world",
			want: false,
		},
		{
			name: "unicode string",
			s:    "hello-世界",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.HasUnicodeChars(tt.s))
		})
	}
}

func TestIsDomainNotRegisteredOrExpiredError(t *testing.T) {
	newErrorResponse := func(statusCode int, message string) error {
		return &dnsimple.ErrorResponse{
			Response: dnsimple.Response{
				HTTPResponse: &http.Response{StatusCode: statusCode},
			},
			Message: message,
		}
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matching 400 error",
			err:  newErrorResponse(http.StatusBadRequest, "Change rejected: domain is not registered or expired"),
			want: true,
		},
		{
			name: "matching 400 error with different casing",
			err:  newErrorResponse(http.StatusBadRequest, "Change rejected: Domain Is Not Registered Or Expired"),
			want: true,
		},
		{
			name: "unrelated 400 error",
			err:  newErrorResponse(http.StatusBadRequest, "Validation failed"),
			want: false,
		},
		{
			name: "matching message but not a 400",
			err:  newErrorResponse(http.StatusNotFound, "domain is not registered or expired"),
			want: false,
		},
		{
			name: "non-dnsimple error",
			err:  errors.New("boom"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, utils.IsDomainNotRegisteredOrExpiredError(tt.err))
		})
	}
}
