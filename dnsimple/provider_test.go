package dnsimple

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"dnsimple": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_sandbox(t *testing.T) {
	if v := os.Getenv("DNSIMPLE_SANDBOX"); v != "" {
		provider := testAccProvider.Meta().(*Client)
		if provider.config.Sandbox != true {
			t.Fatal("Config Sandbox Flag does not equal True!")
		}

		if provider.client.BaseURL != "https://api.sandbox.dnsimple.com" {
			t.Fatalf("Client.BaseURL is not the expected sandbox URL! Currently set to: %s", provider.client.BaseURL)
		}
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DNSIMPLE_EMAIL"); v != "" {
		t.Fatal("DNSIMPLE_EMAIL is no longer required for DNSimple API v2")
	}

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
