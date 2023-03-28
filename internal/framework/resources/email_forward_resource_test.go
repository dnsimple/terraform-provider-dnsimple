package resources_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
)

func TestAccEmailForwardResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckEmailForwardResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEmailForwardResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dnsimple_email_forward.test", "id"),
					resource.TestCheckResourceAttrSet("dnsimple_email_forward.test", "alias_email"),
					resource.TestCheckResourceAttr("dnsimple_email_forward.test", "domain", domainName),
					resource.TestCheckResourceAttr("dnsimple_email_forward.test", "alias_name", "hello"),
					resource.TestCheckResourceAttr("dnsimple_email_forward.test", "destination_email", "hi@hey.com"),
				),
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

func testAccEmailForwardResourceConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_email_forward" "test" {
	domain = %[1]q
	alias_name = "hello"
	destination_email = "hi@hey.com"
}`, domainName)
}
