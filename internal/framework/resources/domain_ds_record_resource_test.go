package resources_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
)

func TestAccDomainDsRecordResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_domain_ds_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainDsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDsRecordResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "13"),
					resource.TestCheckResourceAttr(resourceName, "digest", "2FD4E1C67A2D28FCED849EE1BB76E7391B93EB12"),
					resource.TestCheckResourceAttr(resourceName, "digest_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "key_tag", "12345"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDomainDsRecordImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// TODO: Add test for update
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDomainDsRecordImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["domain_id"], rs.Primary.ID), nil
	}
}

func testAccCheckDomainDsRecordResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_domain_ds_record" {
			continue
		}

		dsIdRaw := rs.Primary.Attributes["id"]
		domainId := rs.Primary.Attributes["domain_id"]

		dsId, err := strconv.ParseInt(dsIdRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert domain delegation signer record ID to int: %s", err)
		}

		_, err = dnsimpleClient.Domains.GetDelegationSignerRecord(context.Background(), testAccAccount, domainId, dsId)

		if err == nil {
			return fmt.Errorf("domain delegation signer record still exists")
		}
	}
	return nil
}

func testAccDomainDsRecordResourceConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_domain_ds_record" "test" {
	domain_id = %[1]q
	algorithm = "13"
	digest = "2FD4E1C67A2D28FCED849EE1BB76E7391B93EB12"
	digest_type = "2"
	keytag = "12345"
}`, domainName)
}
