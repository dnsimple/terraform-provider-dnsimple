package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsimpleDomain_import(t *testing.T) {
	resourceName := "dnsimple_domain.foobar"
	domain := "testing" + os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleDomainConfigImport, domain),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCheckDnsimpleDomainConfigImport = `
resource "dnsimple_domain" "foobar" {
	name = "%s"
}`
