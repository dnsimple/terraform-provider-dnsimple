package resources_test

import (
	"context"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
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
	if os.Getenv("TF_ACC") != "" {
		token := os.Getenv("DNSIMPLE_TOKEN")
		account := os.Getenv("DNSIMPLE_ACCOUNT")

		dnsimpleClient = dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))
		dnsimpleClient.UserAgent = "terraform-provider-dnsimple/test"
		dnsimpleClient.BaseURL = consts.BaseURLSandbox
		testAccAccount = account
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("DNSIMPLE_TOKEN"); v == "" {
		t.Fatal("DNSIMPLE_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("DNSIMPLE_ACCOUNT"); v == "" {
		t.Fatal("DNSIMPLE_ACCOUNT must be set for acceptance tests")
	}

	if v := os.Getenv("DNSIMPLE_DOMAIN"); v == "" {
		t.Fatal("DNSIMPLE_DOMAIN must be set for acceptance tests. The domain is used to create and destroy record against.")
	}
}
