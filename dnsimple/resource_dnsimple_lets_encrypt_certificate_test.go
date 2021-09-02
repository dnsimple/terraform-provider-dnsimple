package dnsimple

import (
	"context"
	"os"
	"testing"

	_ "github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDNSimpleLetsEncryptCertificateCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLetsEncryptConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLetsEncryptCertificate("dnsimple_lets_encrypt_certificate.foobar", os.Getenv("DNSIMPLE_DOMAIN"))),
			},
		},
	})
}

func testAccCheckLetsEncryptCertificate(n string, domainName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		provider := testAccProvider.Meta().(*Client)

		// This is not going  to  work...
		_, err := provider.client.Certificates.GetCertificate(context.Background(), provider.config.Account, domainName, 123)

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccLetsEncryptConfig = `
resource "dnsimple_lets_encrypt_certificate" "foobar" {
	domain_id = "%s"
	contact_id = "2824"
	auto_renew = false
	name = "acc-test-certs"
	alternate_names = []
}`