package dnsimple

import (
	"context"
	"log"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"golang.org/x/oauth2"
)

type Config struct {
	Email   string
	Account string
	Token   string

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
	client.UserAgent = httpclient.TerraformUserAgent(c.terraformVersion)

	provider := &Client{
		client: client,
		config: c,
	}

	log.Printf("[INFO] DNSimple Client configured for account: %s", c.Account)

	return provider, nil
}
