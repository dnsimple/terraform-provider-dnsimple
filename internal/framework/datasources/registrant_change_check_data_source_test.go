package datasources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccRegistrantChangeCheckDataSource(t *testing.T) {
	// Get convert the contact id to int
	contactId := os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_CONTACT_ID")
	domainName := os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_DOMAIN")
	resourceName := "data.dnsimple_registrant_change_check.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRegistrantChangeCheckDataSourceConfig(contactId, domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "contact_id", contactId),
					resource.TestCheckResourceAttr(resourceName, "domain_id", domainName),
					resource.TestCheckResourceAttrSet(resourceName, "extended_attributes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "registry_owner_change"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccRegistrantChangeCheckDataSourceConfig(contactId, domainName string) string {
	return fmt.Sprintf(`
data "dnsimple_registrant_change_check" "test" {
	contact_id = %[1]q
	domain_id = %[2]q
}`, contactId, domainName)
}
