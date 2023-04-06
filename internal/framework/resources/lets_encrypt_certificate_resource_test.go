package resources_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
)

func TestAccLetsEncryptCertificateResource(t *testing.T) {
	if os.Getenv("DNSIMPLE_SANDBOX") != "false" {
		t.Skip("DNSIMPLE_SANDBOX is not set to `false` (read in CONTRIBUTING.md how to run this test)")
		return
	}
	resourceName := "dnsimple_lets_encrypt_certificate.test"

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
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "domain_id", domainId),
					resource.TestCheckResourceAttr(resourceName, "name", certName),
					resource.TestCheckResourceAttrSet(resourceName, "years"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "authority_identifier"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", strconv.FormatBool(certAutoRenew)),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "csr"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", certSigAlg),
				),
			},
			// Update is a no-op
			// Delete testing automatically occurs in TestCase
		},
	})
}

// We cannot delete certificates from the server.
func testAccCheckLetsEncryptCertificateResourceDestroy(state *terraform.State) error {
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
