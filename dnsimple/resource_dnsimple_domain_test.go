package dnsimple

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDNSimpleDomainCreate(t *testing.T) {
	domainAttributes := dnsimple.Domain{Name: "testing" + os.Getenv("DNSIMPLE_DOMAIN")}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleDomainConfig, domainAttributes.Name),
				Check: resource.ComposeTestCheckFunc(
					testAccDNSimpleCreateDomain("dnsimple_domain.foobar", "name")),
			},
		},
	})
}

func testAccCheckDNSimpleDomainDestroy(state *terraform.State) error {
	provider := testAccProvider.Meta().(*Client)

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_domain" {
			continue
		}

		domainName := rs.Primary.Attributes["name"]
		_, err := provider.client.Domains.GetDomain(context.Background(), provider.config.Account, domainName)

		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}
	return nil
}

func testAccDNSimpleCreateDomain(n string, key string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		configuration := state.RootModule().Resources[n]
		domainName := configuration.Primary.Attributes[key]
		provider := testAccProvider.Meta().(*Client)

		_, err := provider.client.Domains.GetDomain(context.Background(), provider.config.Account, domainName)

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccCheckDnsimpleDomainConfig = `
resource "dnsimple_domain" "foobar" {
	name = "%s"
}`
