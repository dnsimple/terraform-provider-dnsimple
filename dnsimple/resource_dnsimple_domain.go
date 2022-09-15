package dnsimple

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSimpleDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleDomainCreate,
		ReadContext:   resourceDNSimpleDomainRead,
		DeleteContext: resourceDNSimpleDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"registrant_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unicode_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_renew": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"private_whois": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceDNSimpleDomainCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	domainAttributes := dnsimple.Domain{
		Name: data.Get("name").(string),
	}

	response, err := provider.client.Domains.CreateDomain(ctx, provider.config.Account, domainAttributes)

	if err != nil {
		return diag.Errorf("Failed to create DNSimple Domain: %s", err)
	}

	data.SetId(strconv.FormatInt(response.Data.ID, 10))

	return resourceDNSimpleDomainRead(ctx, data, meta)
}

func resourceDNSimpleDomainRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	domainID := data.Id()

	response, err := provider.client.Domains.GetDomain(ctx, provider.config.Account, domainID)

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Domain Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Failed to get DNSimple Domain: %s", err)
	}

	domain := response.Data
	data.Set("account_id", domain.AccountID)
	data.Set("registrant_id", domain.RegistrantID)
	data.Set("name", domain.Name)
	data.Set("unicode_name", domain.UnicodeName)
	data.Set("state", domain.State)
	data.Set("auto_renew", domain.AutoRenew)
	data.Set("private_whois", domain.PrivateWhois)

	return nil

}

func resourceDNSimpleDomainDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple Record: %s, %s", data.Get("name").(string), data.Id())

	domainID := data.Id()

	_, err := provider.client.Domains.DeleteDomain(ctx, provider.config.Account, domainID)
	if err != nil {
		return diag.Errorf("Error deleting DNSimple Record: %s", err)
	}

	return nil
}
