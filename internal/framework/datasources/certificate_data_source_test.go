package datasources_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

func TestAccCertificateDataSource(t *testing.T) {
	if os.Getenv("DNSIMPLE_SANDBOX") != "true" {
		t.Skip("DNSIMPLE_SANDBOX is not set to `true` (read in CONTRIBUTING.md how to run this test)")
		return
	}
	domain := os.Getenv("DNSIMPLE_DOMAIN")
	certificateId := os.Getenv("DNSIMPLE_CERTIFICATE_ID")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { utils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.NewProto6ProviderFactory(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateDataSourceConfig(domain, certificateId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.dnsimple_certificate.test", "domain", domain),
					resource.TestCheckResourceAttr("data.dnsimple_certificate.test", "certificate_id", certificateId),
				),
			},
		},
	})
}

func testAccCertificateDataSourceConfig(domainName string, certificateId string) string {
	return fmt.Sprintf(`
data "dnsimple_certificate" "test" {
	domain = %[1]q
	certificate_id = %[2]q
}`, domainName, certificateId)
}
