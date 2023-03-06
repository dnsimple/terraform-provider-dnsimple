package dnsimple

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAccDNSimpleZoneRecord_Basic(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigBasic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					testAccCheckDNSimpleZoneRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "192.168.0.10"),
				),
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_CreateMxWithPriority(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigMx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "mx.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "priority", "5"),
				),
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_Updated(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigBasic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					testAccCheckDNSimpleZoneRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "192.168.0.10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigNewValue, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					testAccCheckDNSimpleZoneRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "192.168.0.11"),
				),
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_disappears(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigBasic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					testAccCheckDNSimpleZoneRecordDisappears(&record, domain),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_UpdatedMx(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigMx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "mx.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "priority", "5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigMxNewValue, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "zone_name", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "value", "mx2.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_zone_record.foobar", "priority", "10"),
				),
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_Prefetch(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			os.Setenv("PREFETCH", "1")
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDnsimpleZoneRecordConfigBasic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists("dnsimple_zone_record.foobar", &record),
					testAccCheckDnsimpleZoneRecordPrefetch("dnsimple_zone_record.foobar"),
				),
			},
		},
	})
}

func TestAccDNSimpleZoneRecord_Prefetch_ForEach(t *testing.T) {
	// Issue: https://github.com/dnsimple/terraform-provider-dnsimple/issues/80
	// This test is to ensure that the prefetch behaviour is deterministic
	var recordSet1 dnsimple.ZoneRecord
	var recordSet2 dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName1 := "dnsimple_zone_record.for_each_example_1"
	resourceName2 := "dnsimple_zone_record.for_each_example_2"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			os.Setenv("PREFETCH", "1")
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDNSimpleZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccExampleDnsimpleZoneRecordsForEach, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleZoneRecordExists(resourceName1, &recordSet1),
					testAccCheckDNSimpleZoneRecordExists(resourceName2, &recordSet2),
				),
			},
			{
				Config: fmt.Sprintf(testAccExampleDnsimpleZoneRecordsForEach, domain),
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

func testAccCheckDnsimpleZoneRecordPrefetch(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		provider := testAccProvider.Meta().(*Client)

		if len(provider.cache) == 0 {
			return errors.New("cache wasn't populated")
		}

		return nil
	}
}

func testAccCheckDNSimpleZoneRecordDisappears(record *dnsimple.ZoneRecord, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		provider := testAccProvider.Meta().(*Client)

		_, err := provider.client.Zones.DeleteRecord(context.Background(), provider.config.Account, domain, record.ID)
		if err != nil {
			return err
		}

		return nil
	}

}

func testAccCheckDNSimpleZoneRecordDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dnsimple_zone_record" {
			continue
		}

		recordID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		_, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, rs.Primary.Attributes["zone_name"], recordID)
		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}

	return nil
}

func testAccCheckDNSimpleZoneRecordAttributes(record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if want, got := "192.168.0.10", record.Content; want != got {
			return fmt.Errorf("Record.Content expected to be %v, got %v", want, got)
		}

		return nil
	}
}

func testAccCheckDNSimpleZoneRecordAttributesUpdated(record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if want, got := "192.168.0.11", record.Content; want != got {
			return fmt.Errorf("Record.Content expected to be %v, got %v", want, got)
		}

		return nil
	}
}

func testAccCheckDNSimpleZoneRecordExists(n string, record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		provider := testAccProvider.Meta().(*Client)

		recordID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		resp, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, rs.Primary.Attributes["zone_name"], recordID)
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

func testZoneRecordInstanceStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"domain": "example.com",
	}
}

func testZoneRecordInstanceStateDataV1() map[string]interface{} {
	v0 := testZoneRecordInstanceStateDataV0()
	return map[string]interface{}{
		"zone_name": v0["domain"],
	}
}

func TestResourceExampleInstanceStateUpgradeV0(t *testing.T) {
	expected := testZoneRecordInstanceStateDataV1()
	actual, err := resourceDNSimpleZoneRecordInstanceStateUpgradeV0(context.Background(), testZoneRecordInstanceStateDataV0(), nil)
	assert.NoError(t, err, "error migrating state")
	assert.Equal(t, expected, actual)
}

const testAccCheckDnsimpleZoneRecordConfigBasic = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
}`

const testAccExampleDnsimpleZoneRecordsForEach = `
resource "dnsimple_zone_record" "for_each_example_1" {
  zone_name = %[1]q

  name      = ""
  value     = "1.1.1.1"
  type      = "A"
  ttl       = 600
}

resource "dnsimple_zone_record" "for_each_example_2" {
  zone_name = %[1]q

  name      = "1a"
  value     = "1.1.1.1"
  type      = "A"
  ttl       = 600
}`

const testAccCheckDnsimpleZoneRecordConfigNewValue = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = "terraform"
	value = "192.168.0.11"
	type = "A"
	ttl = 3600
}`

const testAccCheckDnsimpleZoneRecordConfigMx = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = ""
	value = "mx.example.com"
	type = "MX"
	ttl = 3600
	priority = 5
}`

const testAccCheckDnsimpleZoneRecordConfigMxNewValue = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = ""
	value = "mx2.example.com"
	type = "MX"
	ttl = 3600
	priority = 10
}`

func TestImporterSplitId(t *testing.T) {
	assert.Equal(t, []string{"example.com", "1234"}, importerSplitId("example.com_1234"))
	assert.Equal(t, []string{"record.example.com", "1234"}, importerSplitId("record.example.com_1234"))
	assert.Equal(t, []string{"_record.example.com", "1234"}, importerSplitId("_record.example.com_1234"))
	assert.Nil(t, importerSplitId("_my_domain.example.com"))
	assert.Nil(t, importerSplitId("_my_domain.example.com 1234"))
}
