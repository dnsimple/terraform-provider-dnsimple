package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDnsimpleRecord_import(t *testing.T) {
	resourceName := "dnsimple_zone_record.foobar"
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleRecordConfigImport, domain),
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

const testAccCheckDnsimpleRecordConfigImport = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
}`
