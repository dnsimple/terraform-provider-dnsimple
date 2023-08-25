package registered_domain_contact_test

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
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
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

func TestAccRegisteredDomainContactResource_WithExtendedAttrs(t *testing.T) {
	contactID, err := strconv.Atoi(os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_CONTACT_ID"))
	if err != nil {
		t.Fatal(err)
	}
	domainName := os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_DOMAIN")
	resourceName := "dnsimple_registered_domain_contact.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckRegisteredDomainContact(t) },
		ProtoV6ProviderFactories: test_utils.TestAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckRegisteredDomainContactResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRegisteredDomainContactResourceConfig_WithExtendedAttrs(domainName, contactID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "registry_owner_change"),
					resource.TestCheckResourceAttr(resourceName, "domain_id", domainName),
					resource.TestCheckResourceAttr(resourceName, "contact_id", fmt.Sprintf("%d", contactID)),
					resource.TestCheckResourceAttr(resourceName, "extended_attributes.x-eu-registrant-citizenship", "bg"),
				),
				// We expect the plan to be non-empty because we are creating a registrant change that will not be completed
				// and we will attempt to converge it by setting the state to completed
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccRegisteredDomainImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "registry_owner_change"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPreCheckRegisteredDomainContact(t *testing.T) {
	test_utils.TestAccPreCheck(t)
	if os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_CONTACT_ID") == "" {
		t.Fatal("DNSIMPLE_REGISTRANT_CHANGE_CONTACT_ID must be set for acceptance tests")
	}
	if os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_DOMAIN") == "" {
		t.Fatal("DNSIMPLE_REGISTRANT_CHANGE_DOMAIN must be set for acceptance tests")
	}
}

func testAccRegisteredDomainImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.Attributes["id"] == "" {
			return "", errors.New("No resource Id set")
		}

		return rs.Primary.Attributes["id"], nil
	}
}

func testAccCheckRegisteredDomainContactResourceDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "dnsimple_registered_domain_contact" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}

		if rs.Primary.Attributes["state"] == consts.RegistrantChangeStateCompleted ||
			rs.Primary.Attributes["state"] == consts.RegistrantChangeStateCancelled ||
			rs.Primary.Attributes["state"] == consts.RegistrantChangeStateCancelling {
			continue
		}

		_, err = dnsimpleClient.Registrar.DeleteRegistrantChange(context.Background(), testAccAccount, id)

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccRegisteredDomainContactResourceConfig_WithExtendedAttrs(domainName string, contactId int) string {
	return fmt.Sprintf(`
resource "dnsimple_registered_domain_contact" "test" {
	domain_id = %[1]q
	contact_id = %[2]d
	extended_attributes = {
		"x-eu-registrant-citizenship" = "bg"
	}
	timeouts = {
		create = "40s"
	}
}`, domainName, contactId)
}
