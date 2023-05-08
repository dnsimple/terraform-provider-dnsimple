package registered_domain_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

var (
	// dnsimpleClient is the DNSimple client used to make API calls during
	// acceptance testing.
	dnsimpleClient *dnsimple.Client
	// testAccAccount is the DNSimple account used to make API calls during
	// acceptance testing.
	testAccAccount string
)

func init() {
	// If we are running acceptance tests TC_ACC then we initialize the DNSimple client
	// with the credentials provided in the environment variables.
	dnsimpleClient, testAccAccount = test_utils.LoadDNSimpleTestClient()
}

func TestAccRegisteredDomainResource(t *testing.T) {
	// Get convert the contact id to int
	contactID, err := strconv.Atoi(os.Getenv("DNSIMPLE_CONTACT_ID"))
	if err != nil {
		t.Fatal(err)
	}
	domainName := utils.RandomString(62) + ".com"
	resourceName := "dnsimple_registered_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRegisteredDomain(t) },
		ProtoV6ProviderFactories: test_utils.TestAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckRegisteredDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainResourceConfig(domainName, contactID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "registered"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_registration.id"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "whois_privacy_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dnssec_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccRegisteredDomainImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRegisteredDomainResource_WithExtendedAttrs(t *testing.T) {
	// var resourceID string
	// Get convert the contact id to int
	contactID, err := strconv.Atoi(os.Getenv("DNSIMPLE_CONTACT_ID"))
	if err != nil {
		t.Fatal(err)
	}
	domainName := utils.RandomString(62) + ".eu"
	resourceName := "dnsimple_registered_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRegisteredDomain(t) },
		ProtoV6ProviderFactories: test_utils.TestAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckRegisteredDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainResourceConfig_WithExtendedAttrs(domainName, contactID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "registered"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_registration.id"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "whois_privacy_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dnssec_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
					resource.TestCheckResourceAttr(resourceName, "extended_attributes.x-eu-registrant-citizenship", "bg"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccRegisteredDomainImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRegisteredDomainResource_WithOptions(t *testing.T) {
	// Get convert the contact id to int
	contactID, err := strconv.Atoi(os.Getenv("DNSIMPLE_CONTACT_ID"))
	if err != nil {
		t.Fatal(err)
	}
	domainName := utils.RandomString(62) + ".com"
	resourceName := "dnsimple_registered_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRegisteredDomain(t) },
		ProtoV6ProviderFactories: test_utils.TestAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckRegisteredDomainResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainResourceConfig_WithOptions(domainName, contactID, false, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "registered"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_registration.id"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "whois_privacy_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "dnssec_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				Config: testAccRegisteredDomainResourceConfig_WithOptions(domainName, contactID, true, false, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "registered"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_registration.id"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "whois_privacy_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "dnssec_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "expires_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccRegisteredDomainImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccRegisteredDomainResource_ImportedWithDomainOnly(t *testing.T) {
	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_registered_domain.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRegisteredDomain(t) },
		ProtoV6ProviderFactories: test_utils.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				ResourceName:      resourceName,
				Config:            testAccRegisteredDomainResourceConfig(domainName, 1234),
				ImportStateId:     domainName,
				ImportState:       true,
				ImportStateVerify: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", domainName),
					resource.TestCheckResourceAttr(resourceName, "state", "registered"),
					resource.TestCheckNoResourceAttr(resourceName, "domain_registration"),
				),
			},
		},
	})
}

func testAccPreCheckRegisteredDomain(t *testing.T) {
	test_utils.TestAccPreCheck(t)
	if os.Getenv("DNSIMPLE_CONTACT_ID") == "" {
		t.Fatal("DNSIMPLE_CONTACT_ID must be set for acceptance tests")
	}
}

func testAccRegisteredDomainImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.Attributes["name"] == "" {
			return "", errors.New("No resource Name set")
		}

		if rs.Primary.Attributes["domain_registration.id"] == "" {
			return "", errors.New("No domain registration ID set")
		}

		return fmt.Sprintf("%s_%s", rs.Primary.Attributes["name"], rs.Primary.Attributes["domain_registration.id"]), nil
	}
}

func testAccCheckRegisteredDomainResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_domain" {
			continue
		}

		domainName := rs.Primary.Attributes["name"]
		_, err := dnsimpleClient.Domains.GetDomain(context.Background(), testAccAccount, domainName)

		if err == nil {
			return fmt.Errorf("record still exists")
		}
	}
	return nil
}

func testAccRegisteredDomainResourceConfig(domainName string, contactId int) string {
	return fmt.Sprintf(`
resource "dnsimple_registered_domain" "test" {
	name = %[1]q

	contact_id = %[2]d
}`, domainName, contactId)
}

func testAccRegisteredDomainResourceConfig_WithExtendedAttrs(domainName string, contactId int) string {
	return fmt.Sprintf(`
resource "dnsimple_registered_domain" "test" {
	name = %[1]q

	contact_id = %[2]d
	extended_attributes = {
		"x-eu-registrant-citizenship" = "bg"
	}
}`, domainName, contactId)
}

func testAccRegisteredDomainResourceConfig_WithOptions(domainName string, contactId int, withAutoRenew, withWhoisPrivacy, withDNSSEC bool) string {
	return fmt.Sprintf(`
resource "dnsimple_registered_domain" "test" {
	name = %[1]q
	contact_id = %[2]d

	auto_renew_enabled = %[3]t
	whois_privacy_enabled = %[4]t
	dnssec_enabled = %[5]t
}`, domainName, contactId, withAutoRenew, withWhoisPrivacy, withDNSSEC)
}
