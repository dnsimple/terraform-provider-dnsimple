package dnsimple

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					testAccCheckDnsimpleZoneRecordPrefetch("dnsimple_zone_record.foobar", &record),
				),
			},
		},
	})
}

func testAccCheckDnsimpleZoneRecordPrefetch(n string, record *dnsimple.ZoneRecord) resource.TestCheckFunc {
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
	actual, err := resourceDNSimpleZoneRecordInstanceStateUpgradeV0(nil, testZoneRecordInstanceStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}

const testAccCheckDnsimpleZoneRecordConfigBasic = `
resource "dnsimple_zone_record" "foobar" {
	zone_name = "%s"

	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
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
