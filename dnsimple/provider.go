package dnsimple

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a schema.Provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
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

			"sandbox": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSIMPLE_SANDBOX", nil),
				Description: "Flag to enable the sandbox API.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dnsimple_domain":        resourceDNSimpleDomain(),
			"dnsimple_email_forward": resourceDNSimpleEmailForward(),
			"dnsimple_zone_record":   resourceDNSimpleZoneRecord(),
			"dnsimple_record":        resourceDNSimpleZoneRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dnsimple_zone": datasourceDNSimpleZone(),
		},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := schema.Provider{}.TerraformVersion
			if terraformVersion == "" {
				// Terraform 0.12 introduced this field to the protocol
				// We can therefore assume that if it's missing it's 0.10 or 0.11
				terraformVersion = "0.11+compatible"
			}
			config := Config{
				Token:            data.Get("token").(string),
				Account:          data.Get("account").(string),
				Sandbox:          data.Get("sandbox").(bool),
				terraformVersion: terraformVersion,
			}

			return config.Client()
		},
	}
	return provider
}
