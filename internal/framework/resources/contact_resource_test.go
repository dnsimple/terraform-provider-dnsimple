package resources_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
)

func TestAccContactResource(t *testing.T) {
	resourceName := "dnsimple_contact.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckContactResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccContactResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "first_name"),
					resource.TestCheckResourceAttr(resourceName, "state", "last_name"),
					resource.TestCheckResourceAttr(resourceName, "state", "address1"),
					resource.TestCheckResourceAttr(resourceName, "state", "city"),
					resource.TestCheckResourceAttr(resourceName, "state", "state_province"),
					resource.TestCheckResourceAttr(resourceName, "state", "postal_code"),
					resource.TestCheckResourceAttr(resourceName, "state", "country"),
					resource.TestCheckResourceAttr(resourceName, "state", "phone"),
					resource.TestCheckResourceAttr(resourceName, "state", "email"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccContactImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO: Add test for update
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccContactImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
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

func testAccCheckContactResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_contact" {
			continue
		}

		contactIdRaw := rs.Primary.Attributes["id"]
		contactId, err := strconv.ParseInt(contactIdRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert contact ID to int: %s", err)
		}

		_, err = dnsimpleClient.Contacts.GetContact(context.Background(), testAccAccount, contactId)

		if err == nil {
			return fmt.Errorf("contact still exists")
		}
	}
	return nil
}

func testAccContactResourceConfig() string {
	// Required attributes only.
	return `
resource "dnsimple_contact" "test" {
	first_name = "Alice"
	last_name = "Appleseed"
	address1 = "123 Main St"
	city = "San Francisco"
	state_province = "CA"
	postal_code = "94105"
	country = "US"
	phone = "+1.5555555555"
	email = "alice.appleseed@example.com"
}`
}
