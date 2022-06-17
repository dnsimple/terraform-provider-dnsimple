package dnsimple

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
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
	err := Provider().InternalValidate()
	assert.NoError(t, err)
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

	suffix := "Consul/0.81"
	got := client.client.UserAgent
	assert.True(t, strings.HasSuffix(got, suffix), "UserAgent should contain extra suffix:\n%s", got)
}

func TestProvider_AccSandbox(t *testing.T) {
	if v := os.Getenv("DNSIMPLE_SANDBOX"); v != "" {
		provider := testAccProvider.Meta().(*Client)
		assert.True(t, provider.config.Sandbox, "config.Sandbox should be true")
		assert.Equal(t, "https://api.sandbox.dnsimple.com", provider.client.BaseURL)
	}
}

func TestProvider_AccPrefetch(t *testing.T) {
	if v := os.Getenv("PREFETCH"); v != "" {
		provider := testAccProvider.Meta().(*Client)
		assert.True(t, provider.config.Prefetch, "config.Prefetch should be true")
	}
}
