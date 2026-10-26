package datasources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccDsRecordsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDsRecordsDataSourceConfig(os.Getenv("DNSIMPLE_DOMAIN")),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.dnsimple_ds_records.test", "id"),
					resource.TestCheckResourceAttrSet("data.dnsimple_ds_records.test", "records.#"),
				),
			},
		},
	})
}

func testAccDsRecordsDataSourceConfig(domainName string) string {
	return fmt.Sprintf(`
data "dnsimple_ds_records" "test" {
	domain = %[1]q
}`, domainName)
}
