package utils

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"
)

func GetDefaultFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func TestAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("DNSIMPLE_TOKEN"); v == "" {
		t.Fatal("DNSIMPLE_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("DNSIMPLE_ACCOUNT"); v == "" {
		t.Fatal("DNSIMPLE_ACCOUNT must be set for acceptance tests")
	}

	if v := os.Getenv("DNSIMPLE_DOMAIN"); v == "" {
		t.Fatal("DNSIMPLE_DOMAIN must be set for acceptance tests. The domain is used to create and destroy record against.")
	}
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
