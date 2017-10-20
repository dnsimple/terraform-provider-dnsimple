package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDNSimpleCertificate_Basic(t *testing.T) {
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleCertificateConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"dnsimple_certificate.foobar", "domain", "hashicorp.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_certificate.foobar", "certificate_id", "7"),
				),
			},
		},
	})
}

const testAccCheckDNSimpleCertificateConfig_basic = `
data "dnsimple_certificate" "foobar" {
	domain         = "%s"
	certificate_id = "7"
}`
