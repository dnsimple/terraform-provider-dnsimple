package dnsimple

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCheckDNSimpleEmailForwardConfig_Basic(t *testing.T) {
	var emailForward dnsimple.EmailForward
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleEmailForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleEmailForwardConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleEmailForwardExists("dnsimple_email_forward.hello", &emailForward),
					testAccCheckDNSimpleEmailForwardAttributes(&emailForward),
					resource.TestCheckResourceAttr(
						"dnsimple_email_forward.hello", "domain", domain),
					// We can't check "from" value here because the API returns `addr@domain` not just `addr`.
					resource.TestCheckResourceAttr(
						"dnsimple_email_forward.hello", "to", "hi@example.org"),
				),
			},
		},
	})
}

func testAccCheckDNSimpleEmailForwardDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dnsimple_email_forward" {
			continue
		}

		emailForwardID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		_, err := provider.client.Domains.GetEmailForward(context.Background(), provider.config.Account, rs.Primary.Attributes["domain"], emailForwardID)
		if err == nil {
			return fmt.Errorf("EmailForward still exists")
		}
	}

	return nil
}

func testAccCheckDNSimpleEmailForwardAttributes(emailForward *dnsimple.EmailForward) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if emailForward.To != "hi@example.org" {
			return fmt.Errorf("Bad content: %s", emailForward.To)
		}

		return nil
	}
}

func testAccCheckDNSimpleEmailForwardAttributesUpdated(emailForward *dnsimple.EmailForward) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if emailForward.To != "contacts@example.org" {
			return fmt.Errorf("Bad content: %s", emailForward.To)
		}

		return nil
	}
}

func testAccCheckDNSimpleEmailForwardExists(n string, emailForward *dnsimple.EmailForward) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Email Forward ID is set")
		}

		provider := testAccProvider.Meta().(*Client)

		emailForwardID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		resp, err := provider.client.Domains.GetEmailForward(context.Background(), provider.config.Account, rs.Primary.Attributes["domain"], emailForwardID)
		if err != nil {
			return err
		}

		foundEmailForward := resp.Data
		if foundEmailForward.ID != emailForwardID {
			return fmt.Errorf("EmailForward not found")
		}

		*emailForward = *foundEmailForward

		return nil
	}
}

const testAccCheckDNSimpleEmailForwardConfig_basic = `
resource "dnsimple_email_forward" "hello" {
	domain = "%s"

	from = "hello"
	to   = "hi@example.org"
}`

const testAccCheckDNSimpleEmailForwardConfig_new_value = `
resource "dnsimple_email_forward" "hello" {
	domain = "%s"

	from = "hello"
	to   = "contacts@example.org"
}`

const testAccCheckDNSimpleEmailForwardConfig_wildcard = `
resource "dnsimple_email_forward" "wildcard" {
	domain = "%s"

	from = "(.*)"
	to = "contacts@example.org"
}
`
