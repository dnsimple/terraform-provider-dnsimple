package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDnsimpleEmailForward_import(t *testing.T) {
	resourceName := "dnsimple_email_forward.wildcard"
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleEmailForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleEmailForwardConfig_import, domain),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: fmt.Sprintf("%s_", domain),
			},
		},
	})
}

const testAccCheckDNSimpleEmailForwardConfig_import = `
resource "dnsimple_email_forward" "wildcard" {
	domain = "%s"

	alias_name = "(.*)"
	destination_email = "contacts@example.org"
}
`
