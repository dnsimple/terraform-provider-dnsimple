package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
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
		"API returned an error",
		err.Message,
	)

	for field, errors := range err.AttributeErrors {
		terraformField := TranslateFieldFromAPIToTerraform(field)

		diagnostics.AddAttributeError(
			path.Root(terraformField),
			fmt.Sprintf("API returned a Validation Error for: %s", terraformField),
			strings.Join(errors, ", "),
		)
	}

	return diagnostics
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
