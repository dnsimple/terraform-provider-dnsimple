package resources_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
)

func TestAccDomainResource(t *testing.T) {
	domainName := "test-" + os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("dnsimple_domain.test", "name", domainName),
					resource.TestCheckResourceAttr("dnsimple_domain.test", "state", "hosted"),
				),
			},
			// Update is a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckDomainResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_domain" {
			continue
		}

		domainName := rs.Primary.Attributes["name"]
		_, err := dnsimpleClient.Domains.GetDomain(context.Background(), testAccAccount, domainName)

		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}
	return nil
}

func testAccDomainResourceConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_domain" "test" {
	name = %[1]q
}`, domainName)
}
