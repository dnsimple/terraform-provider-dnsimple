package dnsimple

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSimpleEmailForward() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleEmailForwardCreate,
		ReadContext:   resourceDNSimpleEmailForwardRead,
		UpdateContext: resourceDNSimpleEmailForwardUpdate,
		DeleteContext: resourceDNSimpleEmailForwardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSimpleEmailForwardImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alias_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alias_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"destination_email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDNSimpleEmailForwardCreate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	emailForwardAttributes := dnsimple.EmailForward{
		From: data.Get("alias_name").(string),
		To:   data.Get("destination_email").(string),
	}

	log.Printf("[DEBUG] DNSimple Email Forward create forwardAttributes: %#v", emailForwardAttributes)

	resp, err := provider.client.Domains.CreateEmailForward(context.Background(), provider.config.Account, data.Get("domain").(string), emailForwardAttributes)
	if err != nil {
		return diag.Errorf("Failed to create DNSimple EmailForward: %s", err)
	}

	data.SetId(strconv.FormatInt(resp.Data.ID, 10))
	log.Printf("[INFO] DNSimple EmailForward ID: %s", data.Id())

	return resourceDNSimpleEmailForwardRead(nil, data, meta)
}

func resourceDNSimpleEmailForwardRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	emailForwardID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Email Forward ID: %s", err)
	}

	resp, err := provider.client.Domains.GetEmailForward(context.Background(), provider.config.Account, data.Get("domain").(string), emailForwardID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Email Forward Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Couldn't find DNSimple Email Forward: %s", err)
	}

	emailForward := resp.Data
	aliasParts := strings.Split(emailForward.From, "@")

	data.Set("alias_name", aliasParts[0])
	data.Set("alias_email", emailForward.From)
	data.Set("destination_email", emailForward.To)

	return nil
}

func resourceDNSimpleEmailForwardUpdate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] DNSimple doesn't support updating email forwards")

	return resourceDNSimpleEmailForwardRead(nil, data, meta)
}

func resourceDNSimpleEmailForwardDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple EmailForward: %s, %s", data.Get("domain").(string), data.Id())

	emailForwardID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting EmailForward ID: %s", err)
	}

	_, err = provider.client.Domains.DeleteEmailForward(context.Background(), provider.config.Account, data.Get("domain").(string), emailForwardID)
	if err != nil {
		return diag.Errorf("Error deleting DNSimple EmailForward: %s", err)
	}

	return nil
}

func resourceDNSimpleEmailForwardImport(_ context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("error Importing dnsimple_email_forward. Please make sure the email forward ID is in the form DOMAIN_EMAILFORWARDID (i.e. example.com_1234)")
	}

	data.SetId(parts[1])
	data.Set("domain", parts[0])

	if err := resourceDNSimpleEmailForwardRead(nil, data, meta); err != nil {
		return nil, fmt.Errorf(err[0].Summary)
	}
	return []*schema.ResourceData{data}, nil
}
