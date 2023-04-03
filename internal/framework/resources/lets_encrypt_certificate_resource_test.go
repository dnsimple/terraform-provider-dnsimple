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

func TestAccLetsEncryptCertificateResource(t *testing.T) {
	if os.Getenv("DNSIMPLE_SANDBOX") != "true" {
		t.Skip("DNSIMPLE_SANDBOX is not set to `true` (read in CONTRIBUTING.md how to run this test)")
		return
	}

	domainId := os.Getenv("DNSIMPLE_DOMAIN")
	certName := os.Getenv("DNSIMPLE_CERTIFICATE_NAME")
	certAutoRenew := os.Getenv("DNSIMPLE_CERTIFICATE_AUTO_RENEW") == "1"
	certSigAlg := os.Getenv("DNSIMPLE_CERTIFICATE_SIGNATURE_ALGORITHM")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLetsEncryptCertificateResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLetsEncryptCertificateResourceConfig(domainId, certName, certAutoRenew, certSigAlg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "id"),
					resource.TestCheckResourceAttr("dnsimple_lets_encrypt_certificate.test", "domain_id", domainId),
					resource.TestCheckResourceAttr("dnsimple_lets_encrypt_certificate.test", "name", certName),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "years"),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "state"),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "authority_identifier"),
					resource.TestCheckResourceAttr("dnsimple_lets_encrypt_certificate.test", "auto_renew", strconv.FormatBool(certAutoRenew)),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "created_at"),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "updated_at"),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "expires_on"),
					resource.TestCheckResourceAttrSet("dnsimple_lets_encrypt_certificate.test", "csr"),
					resource.TestCheckResourceAttr("dnsimple_lets_encrypt_certificate.test", "signature_algorithm", certSigAlg),
				),
			},
			// Update is a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckLetsEncryptCertificateResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_lets_encrypt_certificate" {
			continue
		}

		domainId := rs.Primary.Attributes["domain_id"]
		certId, err := strconv.ParseInt(rs.Primary.Attributes["id"], 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing certificate id: %s", err)
		}
		_, err = dnsimpleClient.Certificates.GetCertificate(context.Background(), testAccAccount, domainId, certId)

		if err != nil {
			return fmt.Errorf("record no longer exists (certificates cannot be deleted)")
		}
	}
	return nil
}

func testAccLetsEncryptCertificateResourceConfig(domainId string, name string, autoRenew bool, signatureAlgorithm string) string {
	return fmt.Sprintf(`
resource "dnsimple_lets_encrypt_certificate" "test" {
	domain_id = %[1]q
	auto_renew = %[3]t
	name = %[2]q
	signature_algorithm = %[4]q
}`, domainId, name, autoRenew, signatureAlgorithm)
}
