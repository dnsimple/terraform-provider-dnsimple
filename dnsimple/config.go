package dnsimple

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"log"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	_ "github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"golang.org/x/oauth2"
)

const (
	baseURLSandbox = "https://api.sandbox.dnsimple.com"
)

type Config struct {
	Email   string
	Account string
	Token   string
	Sandbox bool

	terraformVersion string
}

// Client represents the DNSimple provider client.
// This is a convenient container for the configuration and the underlying API client.
type Client struct {
	client *dnsimple.Client
	config *Config
}

// Client returns a new client for accessing DNSimple.
func (config *Config) Client() (*Client, diag.Diagnostics) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token})
	tc := oauth2.NewClient(context.Background(), ts)

	client := dnsimple.NewClient(tc)
	client.SetUserAgent(fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", config.terraformVersion, meta.SDKVersionString()))
	if config.Sandbox {
		client.BaseURL = baseURLSandbox
	}

	provider := &Client{
		client: client,
		config: config,
	}

	log.Printf("[INFO] DNSimple Client configured for account: %s, sandbox: %v", config.Account, config.Sandbox)

	return provider, nil
}
