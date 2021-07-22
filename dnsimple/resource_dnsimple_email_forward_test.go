package dnsimple

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCheckDNSimpleEmailForwardConfig_Basic(t *testing.T) {
	var emailForward dnsimple.EmailForward
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleEmailForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleEmailForwardConfigBasic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleEmailForwardExists("dnsimple_email_forward.hello", &emailForward),
					testAccCheckDNSimpleEmailForwardAttributes(&emailForward),
					resource.TestCheckResourceAttr(
						"dnsimple_email_forward.hello", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_email_forward.hello", "alias_name", "hello"),
					resource.TestCheckResourceAttr(
						"dnsimple_email_forward.hello", "destination_email", "hi@example.com"),
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

		if want, got := "hi@example.com", emailForward.To; want != got {
			return fmt.Errorf("Forward.To expected to be %v, got %v", want, got)
		}

		return nil
	}
}

func testAccCheckDNSimpleEmailForwardAttributesUpdated(emailForward *dnsimple.EmailForward) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if emailForward.To != "contacts@example.org" {
			return fmt.Errorf("bad content: %s", emailForward.To)
		}

		return nil
	}
}

func testAccCheckDNSimpleEmailForwardExists(n string, emailForward *dnsimple.EmailForward) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Email Forward ID is set")
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

const testAccCheckDnsimpleEmailForwardConfigBasic = `
resource "dnsimple_email_forward" "hello" {
	domain = "%s"

	alias_name 			= "hello"
	destination_email	= "hi@example.com"
}`

const testAccCheckDnsimpleEmailForwardConfigNewValue = `
resource "dnsimple_email_forward" "hello" {
	domain = "%s"

	alias_name	 		= "hello"
	destination_email 	= "changed@example.com"
}`
