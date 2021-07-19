package dnsimple

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	testing2 "testing"
)

func TestAccDnsimpleEmailForward_import(t *testing2.T) {

	resourceName := "dnsimple_email_forward.wildcard"
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleEmailForwardDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckdnsimpleemailforwardconfigImport, domain),
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

const testacccheckdnsimpleemailforwardconfigImport = `
resource "dnsimple_email_forward" "wildcard" {
	domain = "%s"

	alias_name = "(.*)"
	destination_email = "contacts@example.org"
}
`
