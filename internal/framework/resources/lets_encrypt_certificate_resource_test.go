package resources_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccLetsEncryptCertificateResource(t *testing.T) {
	if os.Getenv("DNSIMPLE_SANDBOX") != "false" {
		t.Skip("DNSIMPLE_SANDBOX is not set to `false` (read in CONTRIBUTING.md how to run this test)")
		return
	}
	resourceName := "dnsimple_lets_encrypt_certificate.test"

	domainId := os.Getenv("DNSIMPLE_DOMAIN")
	certName := os.Getenv("DNSIMPLE_CERTIFICATE_NAME")
	certAltNamesRaw := os.Getenv("DNSIMPLE_CERTIFICATE_ALTERNATE_NAMES")
	certAutoRenew := os.Getenv("DNSIMPLE_CERTIFICATE_AUTO_RENEW") == "1"
	certSigAlg := os.Getenv("DNSIMPLE_CERTIFICATE_SIGNATURE_ALGORITHM")

	certAltNames := strings.Split(certAltNamesRaw, ",")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckLetsEncryptCertificateResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLetsEncryptCertificateResourceConfig(domainId, certAutoRenew, certName, certAltNames, certSigAlg),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "domain_id", domainId),
					resource.TestCheckResourceAttr(resourceName, "name", certName),
					resource.TestCheckResourceAttr(resourceName, "alternate_names.#", fmt.Sprintf("%d", len(certAltNames))),
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

func testAccLetsEncryptCertificateResourceConfig(domainId string, autoRenew bool, name string, alternateNames []string, signatureAlgorithm string) string {
	alternateNamesRaw, err := json.Marshal(alternateNames)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`
resource "dnsimple_lets_encrypt_certificate" "test" {
	domain_id = %[1]q
	auto_renew = %[2]t
	name = %[3]q
	alternate_names = %[4]s
	signature_algorithm = %[5]q
}`, domainId, autoRenew, name, string(alternateNamesRaw), signatureAlgorithm)
}
