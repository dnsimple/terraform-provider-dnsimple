package utils

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/v9/dnsimple"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func GetDefaultFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// RandomName generates a random domain name using a UUID v7 with an optional suffix.
//
// It returns a string in the format "uuid.extension" or "uuid-suffix.extension" if suffix is provided.
// Falls back to UUID v4 if v7 generation fails.
func RandomName(extension string, suffix string) string {
	u, err := uuid.NewV7()
	if err != nil {
		// Fallback to v4 if v7 generation fails
		u = uuid.New()
	}
	name := u.String()
	if suffix != "" {
		name = name + "-" + suffix
	}
	return name + "." + extension
}

func HasUnicodeChars(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

func RetryWithTimeout(ctx context.Context, fn func() (error, bool), timeout time.Duration, delay time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		err, suspend := fn()
		if err == nil {
			return nil
		}

		if suspend {
			return err
		}

		if time.Now().After(deadline) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			continue
		}
	}
}

func AttributeErrorsToDiagnostics(err *dnsimple.ErrorResponse) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	diagnostics.AddError(
		"DNSimple API returned an error",
		err.Message,
	)

	for field, errors := range err.AttributeErrors {
		terraformField := TranslateFieldFromAPIToTerraform(field)

		diagnostics.AddAttributeError(
			path.Root(terraformField),
			fmt.Sprintf("DNSimple API validation error for field %s", terraformField),
			strings.Join(errors, ", "),
		)
	}

	return diagnostics
}

// IsDomainNotRegisteredOrExpiredError returns true if err is a DNSimple API
// error indicating that a registrar-level operation was rejected because the
// domain is no longer registered at the registry (e.g. it lapsed and moved
// past its renewal/redemption grace period). The DNSimple API surfaces this
// as an HTTP 400 rather than a 404, so it cannot be treated as a generic
// not-found response.
func IsDomainNotRegisteredOrExpiredError(err error) bool {
	var errorResponse *dnsimple.ErrorResponse
	if !errors.As(err, &errorResponse) {
		return false
	}

	if errorResponse.HTTPResponse == nil || errorResponse.HTTPResponse.StatusCode != http.StatusBadRequest {
		return false
	}

	return strings.Contains(strings.ToLower(errorResponse.Message), "not registered or expired")
}

func TranslateFieldFromAPIToTerraform(field string) string {
	switch field {
	case "record_type":
		return "type"
	case "content":
		return "value"
	default:
		return field
	}
}
