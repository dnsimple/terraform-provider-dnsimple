package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSimpleCertificateBasic(t *testing.T) {
	sandbox := os.Getenv("DNSIMPLE_SANDBOX")
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	if sandbox == "false" {
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(testAccCheckDNSimpleCertificateConfigBasic, domain, os.Getenv("DNSIMPLE_CERTIFICATE_ID")),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(
							"data.dnsimple_certificate.foobar", "domain", domain),
						resource.TestCheckResourceAttr(
							"data.dnsimple_certificate.foobar", "certificate_id", os.Getenv("DNSIMPLE_CERTIFICATE_ID")),
					),
				},
			},
		})
	} else {
		t.Skipf("DNSIMPLE_SANDBOX set to: %s (read in CONTRIBUTING.md how to run this test)", sandbox)
	}
}

const testAccCheckDNSimpleCertificateConfigBasic = `
data "dnsimple_certificate" "foobar" {
	domain         = "%s"
	certificate_id = %s
}`
