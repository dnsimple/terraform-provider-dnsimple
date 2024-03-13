package resources_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	_ "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/resources"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

func TestAccDomainDelegationResource(t *testing.T) {
	domainId := os.Getenv("DNSIMPLE_DOMAIN")
	resourceName := "dnsimple_domain_delegation.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { test_utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckDomainDelegationResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainDelegationResourceConfig(domainId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainId),
					resource.TestCheckResourceAttr(resourceName, "name_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-998.awsdns-60.net"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-1556.awsdns-02.co.uk"),
				),
			},
			{
				Config:             testAccDomainDelegationResourceConfigReversed(domainId),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainId),
					resource.TestCheckResourceAttr(resourceName, "name_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-998.awsdns-60.net"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-1556.awsdns-02.co.uk"),
				),
			},
			{
				Config:             testAccDomainDelegationResourceConfigWithSuffix(domainId),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "domain", domainId),
					resource.TestCheckResourceAttr(resourceName, "name_servers.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-998.awsdns-60.net"),
					resource.TestCheckTypeSetElemAttr(resourceName, "name_servers.*", "ns-1556.awsdns-02.co.uk"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccDomainDelegationImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Deleting simply reliquishes control from Terraform and leaves server state intact.
func testAccCheckDomainDelegationResourceDestroy(state *terraform.State) error {
	return nil
}

func testAccDomainDelegationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return "", errors.New("No resource ID set")
		}

		return rs.Primary.ID, nil
	}
}

func testAccDomainDelegationResourceConfig(domainId string) string {
	return fmt.Sprintf(`
resource "dnsimple_domain_delegation" "test" {
	domain = %[1]q
	name_servers = ["ns-998.awsdns-60.net", "ns-1556.awsdns-02.co.uk"]
}`, domainId)
}

func testAccDomainDelegationResourceConfigReversed(domainId string) string {
	return fmt.Sprintf(`
resource "dnsimple_domain_delegation" "test" {
	domain = %[1]q
	name_servers = ["ns-1556.awsdns-02.co.uk", "ns-998.awsdns-60.net"]
}`, domainId)
}

func testAccDomainDelegationResourceConfigWithSuffix(domainId string) string {
	return fmt.Sprintf(`
resource "dnsimple_domain_delegation" "test" {
	domain = %[1]q
	name_servers = ["ns-1556.awsdns-02.co.uk.", "ns-998.awsdns-60.net"]
}`, domainId)
}
