package utils_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

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

func newErrorResponse(statusCode int, message string) error {
	// ErrorResponse.Error() reads HTTPResponse.Request, so a request has to be present for
	// the error to be printable — which testify does when an assertion about it fails.
	request := &http.Request{Method: http.MethodGet, URL: &url.URL{Scheme: "https", Host: "api.dnsimple.com"}}

	return &dnsimple.ErrorResponse{
		Response: dnsimple.Response{
			HTTPResponse: &http.Response{StatusCode: statusCode, Request: request},
		},
		Message: message,
	}
}

func TestIsDomainNotRegisteredOrExpiredError(t *testing.T) {
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

func TestIsTransientRegistrarError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "the message the registrar endpoints actually return",
			err:  newErrorResponse(http.StatusBadRequest, "Registrar connection failed, please try again later"),
			want: true,
		},
		{
			name: "matching 400 error with different casing",
			err:  newErrorResponse(http.StatusBadRequest, "REGISTRAR CONNECTION FAILED, please try again later"),
			want: true,
		},
		{
			// An ordinary 400 is a hard stop and must not be retried, which is what keeps
			// the stopgap match narrow.
			name: "unrelated 400 error",
			err:  newErrorResponse(http.StatusBadRequest, "Validation failed"),
			want: false,
		},
		{
			// The adjacent 400 the provider already classifies. The two must not overlap:
			// a lapsed domain is permanent and retrying it would only delay the warning.
			name: "domain not registered or expired",
			err:  newErrorResponse(http.StatusBadRequest, "Change rejected: domain is not registered or expired"),
			want: false,
		},
		{
			name: "matching message but not a 400",
			err:  newErrorResponse(http.StatusInternalServerError, "Registrar connection failed, please try again later"),
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
			assert.Equal(t, tt.want, utils.IsTransientRegistrarError(tt.err))
		})
	}
}

func TestRetryWithBackoff(t *testing.T) {
	transient := newErrorResponse(http.StatusBadRequest, "Registrar connection failed, please try again later")

	// The suspend flag the registered_domain read path passes: retry a transient registrar
	// failure, stop on anything else.
	retryTransientOnly := func(err error) (error, bool) {
		return err, err != nil && !utils.IsTransientRegistrarError(err)
	}

	t.Run("returns immediately when the first attempt succeeds", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(nil)
		}, time.Second, time.Millisecond, 10*time.Millisecond)

		assert.NoError(t, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("retries a transient failure until it clears", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			if attempts < 3 {
				return retryTransientOnly(transient)
			}
			return retryTransientOnly(nil)
		}, time.Second, time.Millisecond, 10*time.Millisecond)

		assert.NoError(t, err)
		assert.Equal(t, 3, attempts)
	})

	t.Run("does not retry a suspended error", func(t *testing.T) {
		attempts := 0
		permanent := newErrorResponse(http.StatusBadRequest, "Validation failed")

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(permanent)
		}, time.Second, time.Millisecond, 10*time.Millisecond)

		assert.Equal(t, permanent, err)
		assert.Equal(t, 1, attempts, "a hard stop must surface on the first attempt")
	})

	t.Run("gives up at the budget and returns the last error", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(transient)
		}, 20*time.Millisecond, time.Millisecond, 4*time.Millisecond)

		assert.Equal(t, transient, err)
		assert.Greater(t, attempts, 1, "the budget should allow more than one attempt")
	})

	t.Run("backs off, so a fixed budget spends fewer attempts than a fixed delay would", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(transient)
		}, 100*time.Millisecond, 10*time.Millisecond, 40*time.Millisecond)

		// Waits of 10ms + 20ms + 40ms put the fourth attempt at 70ms, and a further 40ms
		// would overrun the budget, so it stops there. A fixed 10ms delay would have made
		// roughly eleven attempts in the same budget.
		assert.Equal(t, transient, err)
		assert.LessOrEqual(t, attempts, 6)
	})

	t.Run("never starts a wait that would cross the budget", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(transient)
		}, 50*time.Millisecond, 10*time.Millisecond, 40*time.Millisecond)

		// Waits of 10ms and 20ms put the third attempt at 30ms, and a further 40ms would
		// cross the 50ms budget, so it stops there. Checking the deadline only after
		// waiting would instead sleep that 40ms and make a fourth attempt at 70ms.
		//
		// This asserts the attempt count rather than elapsed time on purpose: time.After
		// only guarantees a lower bound, so a loaded runner can overshoot the budget
		// without the implementation being wrong. Overshoot can only lower the attempt
		// count, never raise it, so the bound below cannot flake upward.
		assert.Equal(t, transient, err)
		assert.LessOrEqual(t, attempts, 3)
	})

	t.Run("makes a single attempt rather than spinning when given no delay", func(t *testing.T) {
		attempts := 0

		err := utils.RetryWithBackoff(context.Background(), func() (error, bool) {
			attempts++
			return retryTransientOnly(transient)
		}, time.Minute, 0, 0)

		// Doubling never grows a zero delay, so retrying here would run fn flat out for a
		// full minute against an upstream that is already failing.
		assert.Equal(t, transient, err)
		assert.Equal(t, 1, attempts)
	})

	t.Run("stops when the context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		attempts := 0

		err := utils.RetryWithBackoff(ctx, func() (error, bool) {
			attempts++
			cancel()
			return retryTransientOnly(transient)
		}, time.Minute, time.Millisecond, time.Millisecond)

		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, attempts)
	})
}
