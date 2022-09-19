package dnsimple

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
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
				Sensitive:   true,
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
			"prefetch": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("PREFETCH", nil),
				Description: "Flag to enable the prefetch of zone records.",
			},
			"user_agent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom string to append to the user agent used for sending HTTP requests to the API.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"dnsimple_domain":                   resourceDNSimpleDomain(),
			"dnsimple_email_forward":            resourceDNSimpleEmailForward(),
			"dnsimple_lets_encrypt_certificate": resourceDNSimpleLetsEncryptCertificate(),
			"dnsimple_zone_record":              resourceDNSimpleZoneRecord(),
			"dnsimple_record":                   resourceDNSimpleRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dnsimple_certificate": dataSourceDNSimpleCertificate(),
			"dnsimple_zone":        datasourceDNSimpleZone(),
		},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			terraformVersion := schema.Provider{}.TerraformVersion
			if terraformVersion == "" {
				// Terraform 0.12 introduced this field to the protocol
				// We can therefore assume that if it's missing is 0.10 or 0.11
				terraformVersion = "0.11+compatible"
			}
			config := Config{
				Token:    data.Get("token").(string),
				Account:  data.Get("account").(string),
				Sandbox:  data.Get("sandbox").(bool),
				Prefetch: data.Get("prefetch").(bool),

				userAgentExtra:   data.Get("user_agent").(string),
				terraformVersion: terraformVersion,
			}

			return config.Client(ctx)
		},
	}
	return provider
}

func attributeErrorsToDiagnostics(attributeErrors map[string][]string) diag.Diagnostics {
	result := make([]diag.Diagnostic, 0, len(attributeErrors))

	for field, errors := range attributeErrors {
		terraformField := translateFieldFromAPIToTerraform(field)
		result = append(result, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("API returned a Validation Error for: %s", terraformField),
			Detail:        strings.Join(errors, ", "),
			AttributePath: cty.Path{cty.GetAttrStep{Name: terraformField}},
		})
	}

	return result
}

func translateFieldFromAPIToTerraform(field string) string {
	switch field {
	case "record_type":
		return "type"
	case "content":
		return "value"
	default:
		return field
	}
}
