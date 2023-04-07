package datasources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func TestAccZoneDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccZoneDataSourceConfig(os.Getenv("DNSIMPLE_DOMAIN")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dnsimple_zone.test", "reverse", "false"),
					resource.TestCheckResourceAttrSet("data.dnsimple_zone.test", "id"),
				),
			},
		},
	})
}

func testAccZoneDataSourceConfig(domainName string) string {
	return fmt.Sprintf(`
data "dnsimple_zone" "test" {
	name = %[1]q
}`, domainName)
}
