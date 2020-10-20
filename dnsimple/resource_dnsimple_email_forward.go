package dnsimple

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDNSimpleEmailForward() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSimpleEmailForwardCreate,
		Read:   resourceDNSimpleEmailForwardRead,
		Update: resourceDNSimpleEmailForwardUpdate,
		Delete: resourceDNSimpleEmailForwardDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDNSimpleEmailForwardImport,
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

			"destination_email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDNSimpleEmailForwardCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	emailForwardAttributes := dnsimple.EmailForward{
		From: d.Get("alias_name").(string),
		To:   d.Get("destination_email").(string),
	}

	log.Printf("[DEBUG] DNSimple Email Forward create forwardAttributes: %#v", emailForwardAttributes)

	resp, err := provider.client.Domains.CreateEmailForward(context.Background(), provider.config.Account, d.Get("domain").(string), emailForwardAttributes)
	if err != nil {
		return fmt.Errorf("Failed to create DNSimple EmailForward: %s", err)
	}

	d.SetId(strconv.FormatInt(resp.Data.ID, 10))
	log.Printf("[INFO] DNSimple EmailForward ID: %s", d.Id())

	return resourceDNSimpleEmailForwardRead(d, meta)
}

func resourceDNSimpleEmailForwardRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	emailForwardID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Email Forward ID: %s", err)
	}

	resp, err := provider.client.Domains.GetEmailForward(context.Background(), provider.config.Account, d.Get("domain").(string), emailForwardID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Email Forward Not Found - Refreshing from State")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find DNSimple Email Forward: %s", err)
	}

	emailForward := resp.Data
	aliasParts := strings.Split(emailForward.From, "@")
	d.Set("alias_name", aliasParts[0])
	d.Set("alias_email", emailForward.From)
	d.Set("destination_email", emailForward.To)

	return nil
}

func resourceDNSimpleEmailForwardUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] DNSimple doesn't support updating email forwards")

	return resourceDNSimpleEmailForwardRead(d, meta)
}

func resourceDNSimpleEmailForwardDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple EmailForward: %s, %s", d.Get("domain").(string), d.Id())

	emailForwardID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting EmailForward ID: %s", err)
	}

	_, err = provider.client.Domains.DeleteEmailForward(context.Background(), provider.config.Account, d.Get("domain").(string), emailForwardID)
	if err != nil {
		return fmt.Errorf("Error deleting DNSimple EmailForward: %s", err)
	}

	return nil
}

func resourceDNSimpleEmailForwardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Error Importing dnsimple_email_forward. Please make sure the email forward ID is in the form DOMAIN_EMAILFORWARDID (i.e. example.com_1234)")
	}

	d.SetId(parts[1])
	d.Set("domain", parts[0])

	if err := resourceDNSimpleEmailForwardRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
