package dnsimple

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDNSimpleRecord_Basic(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					testAccCheckDNSimpleRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "192.168.0.10"),
				),
			},
		},
	})
}

func TestAccDNSimpleRecord_CreateMxWithPriority(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_mx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "mx.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "priority", "5"),
				),
			},
		},
	})
}

func TestAccDNSimpleRecord_Updated(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					testAccCheckDNSimpleRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "192.168.0.10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_new_value, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					testAccCheckDNSimpleRecordAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", "terraform"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "192.168.0.11"),
				),
			},
		},
	})
}

func TestAccDNSimpleRecord_disappears(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_basic, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					testAccCheckDNSimpleRecordDisappears(&record, domain),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDNSimpleRecord_UpdatedMx(t *testing.T) {
	var record dnsimple.ZoneRecord
	domain := os.Getenv("DNSIMPLE_DOMAIN")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDNSimpleRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_mx, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "mx.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "priority", "5"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDNSimpleRecordConfig_mx_new_value, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSimpleRecordExists("dnsimple_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "name", ""),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "domain", domain),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "value", "mx2.example.com"),
					resource.TestCheckResourceAttr(
						"dnsimple_record.foobar", "priority", "10"),
				),
			},
		},
	})
}

func testAccCheckDNSimpleRecordDisappears(record *dnsimple.ZoneRecord, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		provider := testAccProvider.Meta().(*Client)

		_, err := provider.client.Zones.DeleteRecord(context.Background(), provider.config.Account, domain, record.ID)
		if err != nil {
			return err
		}

		return nil
	}

}

func testAccCheckDNSimpleRecordDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "dnsimple_record" {
			continue
		}

		recordID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		_, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, rs.Primary.Attributes["domain"], recordID)
		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckDNSimpleRecordAttributes(record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if want, got := "192.168.0.10", record.Content; want != got {
			return fmt.Errorf("Record.Content expected to be %v, got %v", want, got)
		}

		return nil
	}
}

func testAccCheckDNSimpleRecordAttributesUpdated(record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if want, got := "192.168.0.11", record.Content; want != got {
			return fmt.Errorf("Record.Content expected to be %v, got %v", want, got)
		}

		return nil
	}
}

func testAccCheckDNSimpleRecordExists(n string, record *dnsimple.ZoneRecord) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		provider := testAccProvider.Meta().(*Client)

		recordID, _ := strconv.ParseInt(rs.Primary.ID, 10, 64)
		resp, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, rs.Primary.Attributes["domain"], recordID)
		if err != nil {
			return err
		}

		foundRecord := resp.Data
		if foundRecord.ID != recordID {
			return fmt.Errorf("Record not found")
		}

		*record = *foundRecord

		return nil
	}
}

const testAccCheckDNSimpleRecordConfig_basic = `
resource "dnsimple_record" "foobar" {
	domain = "%s"

	name = "terraform"
	value = "192.168.0.10"
	type = "A"
	ttl = 3600
}`

const testAccCheckDNSimpleRecordConfig_new_value = `
resource "dnsimple_record" "foobar" {
	domain = "%s"

	name = "terraform"
	value = "192.168.0.11"
	type = "A"
	ttl = 3600
}`

const testAccCheckDNSimpleRecordConfig_mx = `
resource "dnsimple_record" "foobar" {
	domain = "%s"

	name = ""
	value = "mx.example.com"
	type = "MX"
	ttl = 3600
	priority = 5
}`

const testAccCheckDNSimpleRecordConfig_mx_new_value = `
resource "dnsimple_record" "foobar" {
	domain = "%s"

	name = ""
	value = "mx2.example.com"
	type = "MX"
	ttl = 3600
	priority = 10
}`
