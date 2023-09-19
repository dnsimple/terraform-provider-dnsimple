package resources_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccZoneRecordResource(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_zone_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckZoneRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneRecordResourceStandardConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", "terraform."+domainName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "2800"),
					resource.TestCheckResourceAttr(resourceName, "value", "192.168.0.10"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccZoneRecordResourceStandardWithRegionsConfig(domainName, []string{"IAD", "SYD"}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", "terraform."+domainName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "4000"),
					resource.TestCheckResourceAttr(resourceName, "value", "192.168.0.11"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "regions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "regions.0", "IAD"),
					resource.TestCheckResourceAttr(resourceName, "regions.1", "SYD"),
				),
			},
			{
				Config: testAccZoneRecordResourceStandardWithDefaultsConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", "terraform."+domainName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3600"),
					resource.TestCheckResourceAttr(resourceName, "value", "192.168.0.12"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccZoneRecordImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccZoneRecordResourceWithPriority(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_zone_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckZoneRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneRecordResourcePriorityConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", "terraform."+domainName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3600"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccZoneRecordResourceWithPrefetch(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_zone_record.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckZoneRecordResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneRecordResourcePriorityConfig(domainName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "qualified_name", "terraform."+domainName),
					resource.TestCheckResourceAttr(resourceName, "ttl", "3600"),
					resource.TestCheckResourceAttr(resourceName, "priority", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccZoneRecordResource_Prefetch_ForEach(t *testing.T) {
	// Issue: https://github.com/dnsimple/terraform-provider-dnsimple/issues/80
	// This test is to ensure that the prefetch behaviour is deterministic
	var recordSet1 dnsimple.ZoneRecord
	var recordSet2 dnsimple.ZoneRecord
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName1 := "dnsimple_zone_record.for_each_test_1"
	resourceName2 := "dnsimple_zone_record.for_each_test_2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			test_utils.TestAccPreCheck(t)
			os.Setenv("PREFETCH", "1")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneRecordResourceMultipleConfig(domainName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneRecordExists(resourceName1, &recordSet1),
					testAccCheckZoneRecordExists(resourceName2, &recordSet2),
				),
			},
			{
				Config: testAccZoneRecordResourceMultipleConfig(domainName),
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						resourceName := resourceName1
						expectedID := fmt.Sprintf("%d", recordSet1.ID)
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}

						if rs.Primary.ID != expectedID {
							return fmt.Errorf("Expected the ID of the resource (%s) to be the same as the record ID, but got %s and %s", resourceName, rs.Primary.ID, expectedID)
						}
						return nil
					},
					func(s *terraform.State) error {
						resourceName := resourceName2
						expectedID := fmt.Sprintf("%d", recordSet2.ID)
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}

						if rs.Primary.ID != expectedID {
							return fmt.Errorf("Expected the ID of the resource (%s) to be the same as the record ID, but got %s and %s", resourceName, rs.Primary.ID, expectedID)
						}
						return nil
					},
				),
			},
		},
	})
}

func testAccCheckZoneRecordResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_zone_record" {
			continue
		}

		zoneName := rs.Primary.Attributes["zone_name"]
		// Convert the ID to an int64
		recordID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("error converting record ID to int64: %s", err)
		}

		_, err = dnsimpleClient.Zones.GetRecord(context.Background(), testAccAccount, zoneName, recordID)

		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}
	return nil
}

func testAccZoneRecordImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["zone_name"], rs.Primary.ID), nil
	}
}

func testAccCheckZoneRecordExists(n string, record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		recordID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		resp, err := dnsimpleClient.Zones.GetRecord(context.Background(), testAccAccount, rs.Primary.Attributes["zone_name"], recordID)
		if err != nil {
			return err
		}

		foundRecord := resp.Data
		if foundRecord.ID != recordID {
			return fmt.Errorf("record not found")
		}

		*record = *foundRecord

		return nil
	}
}

func testAccZoneRecordResourcePriorityConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_zone_record" "test" {
	zone_name = %[1]q

	name = "terraform"
	value = "mail.example.com"
	type = "MX"
	priority = 10
}`, domainName)
}

func testAccZoneRecordResourceStandardConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_zone_record" "test" {
	zone_name = %[1]q

	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 2800
}`, domainName)
}

func testAccZoneRecordResourceStandardWithRegionsConfig(domainName string, regions []string) string {
	regionsRaw, err := json.Marshal(regions)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`
resource "dnsimple_zone_record" "test" {
	zone_name = %[1]q

	name = "terraform"
	value = "192.168.0.11"
	type = "A"
	ttl = 4000
	regions = %[2]s
}`, domainName, regionsRaw)
}

func testAccZoneRecordResourceStandardWithDefaultsConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_zone_record" "test" {
	zone_name = %[1]q

	name = "terraform"
	value = "192.168.0.12"
	type = "A"
}`, domainName)
}

func testAccZoneRecordResourceMultipleConfig(domainName string) string {
	return fmt.Sprintf(`
resource "dnsimple_zone_record" "for_each_test_1" {
	zone_name = %[1]q

	name      = ""
	value     = "1.1.1.1"
	type      = "A"
	ttl       = 600
}

resource "dnsimple_zone_record" "for_each_test_2" {
	zone_name = %[1]q

	name      = "1a"
	value     = "1.1.1.1"
	type      = "A"
	ttl       = 600
}`, domainName)
}
