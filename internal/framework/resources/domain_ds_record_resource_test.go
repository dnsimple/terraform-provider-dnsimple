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
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccDomainDsRecordResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_ds_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainDsRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDsRecordResourceConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "domain"),
					resource.TestCheckResourceAttr(resourceName, "algorithm", "8"),
					resource.TestCheckResourceAttr(resourceName, "digest", "C3D49CB83734B22CF3EF9A193B94302FA3BB68013E3E149786D40CDC1BBACD93"),
					resource.TestCheckResourceAttr(resourceName, "digest_type", "2"),
					resource.TestCheckResourceAttr(resourceName, "key_tag", "51301"),
					resource.TestCheckResourceAttr(resourceName, "public_key", "AwEAAd4gdAYAeCnAsYYStm/eWd6uRn5XvT14D9DDM9TbmCvLKCuRA6WYz7suLAziJ5hvk2I7aTOVK8Wd1fDmVxHXGg0Jd6P2+GQpg7AGghD+oLeg0I7AesSIKO3o1ffr58x6iIsxVZ+fcC7G6vdr/d8oIJ/SZdAvghQnCNmCm49HLoN6bWJWNJIXzmxFrptvfgfB4B+PVzbquZrJ0W10KrD394U="),
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
			// Updates are a no-op
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

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["domain"], rs.Primary.ID), nil
	}
}

func testAccCheckDomainDsRecordResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_ds_record" {
			continue
		}

		dsIdRaw := rs.Primary.Attributes["id"]
		domain := rs.Primary.Attributes["domain"]

		dsId, err := strconv.ParseInt(dsIdRaw, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to convert domain delegation signer record ID to int: %s", err)
		}

		_, err = dnsimpleClient.Domains.GetDelegationSignerRecord(context.Background(), testAccAccount, domain, dsId)

		if err == nil {
			return fmt.Errorf("domain delegation signer record still exists")
		}
	}
	return nil
}

func testAccDomainDsRecordResourceConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_ds_record" "test" {
	domain = %[1]q
	algorithm = "8"
	digest = "C3D49CB83734B22CF3EF9A193B94302FA3BB68013E3E149786D40CDC1BBACD93"
	digest_type = "2"
	key_tag = "51301"
	public_key = "AwEAAd4gdAYAeCnAsYYStm/eWd6uRn5XvT14D9DDM9TbmCvLKCuRA6WYz7suLAziJ5hvk2I7aTOVK8Wd1fDmVxHXGg0Jd6P2+GQpg7AGghD+oLeg0I7AesSIKO3o1ffr58x6iIsxVZ+fcC7G6vdr/d8oIJ/SZdAvghQnCNmCm49HLoN6bWJWNJIXzmxFrptvfgfB4B+PVzbquZrJ0W10KrD394U="
}`, domainName)
}
