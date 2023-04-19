package utils

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

func GetDefaultFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(b)
}

func HasUnicodeChars(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

func RetryWithTimeout(ctx context.Context, fn func() error, timeout time.Duration, delay time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		err := fn()
		if err == nil {
			return nil
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
