package dnsimple

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSIMPLE_EMAIL", ""),
				Description: "The DNSimple account email address.",
			},

			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSIMPLE_TOKEN", nil),
				Description: "The API v2 token for API operations.",
			},

			"account": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSIMPLE_ACCOUNT", nil),
				Description: "The account for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dnsimple_record": resourceDNSimpleRecord(),
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := p.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return p
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	// DNSimple API v1 requires email+token to authenticate.
	// DNSimple API v2 requires only an OAuth token and in this particular case
	// the reference of the account for API operations (to avoid fetching it in real time).
	//
	// v2 is not backward compatible with v1, therefore return an error in case email is set,
	// to inform the user to upgrade to v2. Also, v1 token is not the same of v2.
	if email := d.Get("email").(string); email != "" {
		return nil, errors.New(
			"DNSimple API v2 requires an account identifier and the new OAuth token. " +
				"Please upgrade your configuration.")
	}

	config := Config{
		Token:            d.Get("token").(string),
		Account:          d.Get("account").(string),
		terraformVersion: terraformVersion,
	}

	return config.Client()
}
