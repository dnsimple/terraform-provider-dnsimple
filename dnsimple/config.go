package dnsimple

import (
	"context"
	"log"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
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
func (c *Config) Client() (*Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: c.Token})
	tc := oauth2.NewClient(context.Background(), ts)

	client := dnsimple.NewClient(tc)
	client.SetUserAgent(httpclient.TerraformUserAgent(c.terraformVersion))
	if c.Sandbox {
		client.BaseURL = baseURLSandbox
	}

	provider := &Client{
		client: client,
		config: c,
	}

	log.Printf("[INFO] DNSimple Client configured for account: %s, sandbox: %v", c.Account, c.Sandbox)

	return provider, nil
}
