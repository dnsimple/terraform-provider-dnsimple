package dnsimple

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSimpleZoneRead(t *testing.T) {
	zoneName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "data.dnsimple_zone.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDnsimpleZoneDatasource, zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", zoneName),
				),
			},
		},
	})
}

const testAccDnsimpleZoneDatasource = `
data "dnsimple_zone" "foobar" {
	name = "%s"
}`
