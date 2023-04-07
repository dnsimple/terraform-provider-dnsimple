package utils

import (
	"os"
	"testing"
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
