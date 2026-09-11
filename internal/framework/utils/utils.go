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

// RetryWithBackoff behaves like RetryWithTimeout, but waits progressively longer between
// attempts: the delay starts at initialDelay and doubles after each failure, up to
// maxDelay. Retrying a struggling upstream on a fixed interval adds load at the worst
// moment, and backing off gives it room to recover within the same overall budget.
//
// As with RetryWithTimeout, fn reports an error plus a suspend flag, and a suspended
// error is returned without further attempts.
//
// initialDelay must be positive. A non-positive delay makes a single attempt rather than
// retrying, because retrying with no wait at all would hammer an upstream that is already
// struggling, which is the opposite of what any caller wants here.
func RetryWithBackoff(ctx context.Context, fn func() (error, bool), timeout time.Duration, initialDelay time.Duration, maxDelay time.Duration) error {
	deadline := time.Now().Add(timeout)

	delay := initialDelay
	if delay <= 0 {
		delay = maxDelay
	}

	for {
		err, suspend := fn()
		if err == nil {
			return nil
		}

		if suspend {
			return err
		}

		// A non-positive delay would spin: the wait returns immediately and doubling never
		// grows it, so fn would run flat out until the deadline.
		if delay <= 0 {
			return err
		}

		// Give up when the next wait would overrun the budget, rather than when the
		// deadline has already passed. Checking after the fact lets the last wait run past
		// the deadline by up to maxDelay, which makes the budget mean less than it says.
		if time.Now().Add(delay).After(deadline) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}

		if delay < maxDelay {
			delay = min(delay*2, maxDelay)
		}
	}
}

// IsTransientRegistrarError returns true if err is a DNSimple API error reporting that
// the upstream registrar connection failed, which is a momentary condition the same
// request can recover from on a retry.
//
// The API renders this as HTTP 400 with a "please try again later" message across the
// registrar-backed endpoints. A 4xx is otherwise a hard stop that a well-behaved client
// never retries, so matching the message is the only signal available and this function
// is deliberately a stopgap. dnsimple/dnsimple-app#36024 tracks giving the condition a
// machine-readable status code; once that lands, this collapses to a status-code check
// and the message matching goes away.
//
// The match is kept narrow for that reason: an unrelated 400 stays a hard stop.
func IsTransientRegistrarError(err error) bool {
	var errorResponse *dnsimple.ErrorResponse
	if !errors.As(err, &errorResponse) {
		return false
	}

	if errorResponse.HTTPResponse == nil || errorResponse.HTTPResponse.StatusCode != http.StatusBadRequest {
		return false
	}

	return strings.Contains(strings.ToLower(errorResponse.Message), "registrar connection failed")
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
//
// The message match was taken from a real response. Requesting the delegation of
// a lapsed domain returns:
//
//	HTTP/1.1 400 Bad Request
//	{"message":"Change rejected: domain is not registered or expired"}
//
// Matching on message text is unavoidable here, since the status code alone does
// not distinguish this from any other rejected request. If the wording changes,
// this returns false and callers fall back to reporting a hard error, so the
// failure mode is a worse message rather than a wrong action.
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
