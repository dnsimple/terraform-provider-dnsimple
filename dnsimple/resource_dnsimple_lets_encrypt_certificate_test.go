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

func TestAccDNSimpleLetsEncryptCertificateCreate(t *testing.T) {
	sandbox := os.Getenv("DNSIMPLE_SANDBOX")

	if sandbox == "false" {
		var certificate dnsimple.Certificate
		domain := os.Getenv("DNSIMPLE_DOMAIN")
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: testAccProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(testAccLetsEncryptConfig, domain, os.Getenv("DNSIMPLE_CERTIFICATE_NAME")),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckLetsEncryptCertificate("dnsimple_lets_encrypt_certificate.foobar", &certificate)),
				},
			},
		})
	} else {
		t.Skipf("DNSIMPLE_SANDBOX set to: %s (read in CONTRIBUTING.md how to run this test)", sandbox)
	}
}

func testAccCheckLetsEncryptCertificate(resourceName string, certificate *dnsimple.Certificate) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("could not find resource")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Certificate ID is set")
		}

		provider := testAccProvider.Meta().(*Client)

		certificateID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)

		response, err := provider.client.Certificates.GetCertificate(context.Background(), provider.config.Account, rs.Primary.Attributes["domain_id"], certificateID)

		if err != nil {
			return err
		}

		foundCertificate := response.Data

		if response.Data.ID != certificateID {
			return fmt.Errorf("not the same certificate")
		}

		*certificate = *foundCertificate

		return nil
	}
}

const testAccLetsEncryptConfig = `
resource "dnsimple_lets_encrypt_certificate" "foobar" {
	domain_id = "%s"
	contact_id = 1234
	auto_renew = false
	name = "%s"
	signature_algorithm = "RSA"
}`
