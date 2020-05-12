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

func resourceDNSimpleRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSimpleRecordCreate,
		Read:   resourceDNSimpleRecordRead,
		Update: resourceDNSimpleRecordUpdate,
		Delete: resourceDNSimpleRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDNSimpleRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
			},

			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "3600",
			},

			"priority": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceDNSimpleRecordCreate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	// Create the new record
	recordAttributes := dnsimple.ZoneRecordAttributes{
		Name:    dnsimple.String(d.Get("name").(string)),
		Type:    d.Get("type").(string),
		Content: d.Get("value").(string),
	}
	if attr, ok := d.GetOk("ttl"); ok {
		recordAttributes.TTL, _ = strconv.Atoi(attr.(string))
	}
	if attr, ok := d.GetOk("priority"); ok {
		recordAttributes.Priority, _ = strconv.Atoi(attr.(string))
	}

	log.Printf("[DEBUG] DNSimple Record create recordAttributes: %#v", recordAttributes)

	resp, err := provider.client.Zones.CreateRecord(context.Background(), provider.config.Account, d.Get("domain").(string), recordAttributes)
	if err != nil {
		return fmt.Errorf("Failed to create DNSimple Record: %s", err)
	}

	d.SetId(strconv.FormatInt(resp.Data.ID, 10))
	log.Printf("[INFO] DNSimple Record ID: %s", d.Id())

	return resourceDNSimpleRecordRead(d, meta)
}

func resourceDNSimpleRecordRead(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	recordID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	resp, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, d.Get("domain").(string), recordID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Record Not Found - Refreshing from State")
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Couldn't find DNSimple Record: %s", err)
	}

	record := resp.Data
	d.Set("domain_id", record.ZoneID)
	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("value", record.Content)
	d.Set("ttl", strconv.Itoa(record.TTL))
	d.Set("priority", strconv.Itoa(record.Priority))

	if record.Name == "" {
		d.Set("hostname", d.Get("domain").(string))
	} else {
		d.Set("hostname", fmt.Sprintf("%s.%s", record.Name, d.Get("domain").(string)))
	}

	return nil
}

func resourceDNSimpleRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	recordID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	recordAttributes := dnsimple.ZoneRecordAttributes{}
	if attr, ok := d.GetOk("name"); ok {
		recordAttributes.Name = dnsimple.String(attr.(string))
	}
	if attr, ok := d.GetOk("type"); ok {
		recordAttributes.Type = attr.(string)
	}
	if attr, ok := d.GetOk("value"); ok {
		recordAttributes.Content = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		recordAttributes.TTL, _ = strconv.Atoi(attr.(string))
	}
	if attr, ok := d.GetOk("priority"); ok {
		recordAttributes.Priority, _ = strconv.Atoi(attr.(string))
	}

	log.Printf("[DEBUG] DNSimple Record update configuration: %#v", recordAttributes)

	_, err = provider.client.Zones.UpdateRecord(context.Background(), provider.config.Account, d.Get("domain").(string), recordID, recordAttributes)
	if err != nil {
		return fmt.Errorf("Failed to update DNSimple Record: %s", err)
	}

	return resourceDNSimpleRecordRead(d, meta)
}

func resourceDNSimpleRecordDelete(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple Record: %s, %s", d.Get("domain").(string), d.Id())

	recordID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return fmt.Errorf("Error converting Record ID: %s", err)
	}

	_, err = provider.client.Zones.DeleteRecord(context.Background(), provider.config.Account, d.Get("domain").(string), recordID)
	if err != nil {
		return fmt.Errorf("Error deleting DNSimple Record: %s", err)
	}

	return nil
}

func resourceDNSimpleRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Error Importing dnsimple_record. Please make sure the record ID is in the form DOMAIN_RECORDID (i.e. example.com_1234)")
	}

	d.SetId(parts[1])
	d.Set("domain", parts[0])

	if err := resourceDNSimpleRecordRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
