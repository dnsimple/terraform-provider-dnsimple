package dnsimple

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

const ProviderNameDNSimple = "dnsimple"

func init() {
	testAccProvider = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameDNSimple: func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
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

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_Impl(t *testing.T) {
	var _ = Provider()
}

func TestProvider_UserAgentExtra(t *testing.T) {
	raw := map[string]interface{}{
		"token":      "a1b2c3d4f5",
		"user_agent": "Consul/0.81",
	}

	provider := Provider()
	provider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	client := provider.Meta().(*Client)

	want := "Consul/0.81"
	if got := client.client.UserAgent; !strings.HasSuffix(got, want) {
		t.Fatalf("Config UserAgent expected to end with `%v`, got `%v`", want, got)
	}
}

func TestProvider_AccSandbox(t *testing.T) {
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

func TestProvider_AccPrefetch(t *testing.T) {
	if v := os.Getenv("PREFETCH"); v != "" {
		provider := testAccProvider.Meta().(*Client)
		if provider.config.Prefetch != true {
			t.Fatal("Config Prefetch Flag does not equal True!")
		}
	}
}
