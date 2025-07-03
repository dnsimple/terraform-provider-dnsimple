package test_utils

import (
	"context"
	"os"
	"testing"

	"github.com/dnsimple/dnsimple-go/v5/dnsimple"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
)

func TestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"dnsimple": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

func TestAccPreCheck(t *testing.T) {
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

func LoadDNSimpleTestClient() (*dnsimple.Client, string) {
	if os.Getenv("TF_ACC") != "" {
		token := os.Getenv("DNSIMPLE_TOKEN")
		account := os.Getenv("DNSIMPLE_ACCOUNT")

		dnsimpleClient := dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))
		dnsimpleClient.UserAgent = "terraform-provider-dnsimple/test"
		dnsimpleClient.BaseURL = consts.BaseURLSandbox

		return dnsimpleClient, account
	}

	return nil, ""
}
