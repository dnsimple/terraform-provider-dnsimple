package resources_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccEmailForwardResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_email_forward.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEmailForwardResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailForwardResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "alias_email"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "alias_name", "hello"),
					resource.TestCheckResourceAttr(resourceName, "destination_email", "hi@hey.com"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccEmailForwardImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update is a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckEmailForwardResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_email_forward" {
			continue
		}

		domainName := rs.Primary.Attributes["name"]
		forwardId, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing email forward id: %s", err)
		}
		_, err = dnsimpleClient.Domains.GetEmailForward(context.Background(), testAccAccount, domainName, forwardId)

		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}
	return nil
}

func testAccEmailForwardImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["domain"], rs.Primary.ID), nil
	}
}

func testAccEmailForwardResourceConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_email_forward" "test" {
	domain = %[1]q
	alias_name = "hello"
	destination_email = "hi@hey.com"
}`, domainName)
}
