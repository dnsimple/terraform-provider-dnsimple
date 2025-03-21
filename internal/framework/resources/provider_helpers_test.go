package resources_test

import (
	"github.com/dnsimple/dnsimple-go/v4/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/test_utils"
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"dnsimple": providerserver.NewProtocol6WithError(provider.New("test")()),
	}

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
