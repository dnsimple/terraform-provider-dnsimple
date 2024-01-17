package resources_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccZoneResource(t *testing.T) {
	zoneName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_zone.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneResourceConfig(zoneName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "reverse", "false"),
					resource.TestCheckResourceAttr(resourceName, "secondary", "false"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
				),
			},
			{
				Config: testAccZoneResourceConfigWithActive(zoneName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "reverse", "false"),
					resource.TestCheckResourceAttr(resourceName, "secondary", "false"),
					resource.TestCheckResourceAttr(resourceName, "active", "false"),
				),
			},
			{
				Config: testAccZoneResourceConfigWithActive(zoneName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
					resource.TestCheckResourceAttr(resourceName, "reverse", "false"),
					resource.TestCheckResourceAttr(resourceName, "secondary", "false"),
					resource.TestCheckResourceAttr(resourceName, "active", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccZoneImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccZoneImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccZoneResourceConfig(zoneName string) string {
	return fmt.Sprintf(`
resource "dnsimple_zone" "test" {
	name = %[1]q
}`, zoneName)
}

func testAccZoneResourceConfigWithActive(zoneName string, active bool) string {
	return fmt.Sprintf(`
resource "dnsimple_zone" "test" {
	name = %[1]q
	active = %[2]t
}`, zoneName, active)
}
