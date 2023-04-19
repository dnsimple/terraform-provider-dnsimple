package resources_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccDomainResource(t *testing.T) {
	domainName := "test-" + os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "hosted"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDomainImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update is a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDomainImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return rs.Primary.ID, nil
	}
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
