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

func resourceDNSimpleRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleRecordCreate,
		ReadContext:   resourceDNSimpleRecordRead,
		UpdateContext: resourceDNSimpleRecordUpdate,
		DeleteContext: resourceDNSimpleRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSimpleRecordImport,
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

func deprecationWarning() {
	fmt.Println("WARNING! This resource (dnsimple_record) is deprecated and will be removed in future versions")
	fmt.Println("Please consider changing your configuration to use dnsimple_zone_record instead")
}

func resourceDNSimpleRecordCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deprecationWarning()
	provider := meta.(*Client)

	// Create the new record
	recordAttributes := dnsimple.ZoneRecordAttributes{
		Name:    dnsimple.String(data.Get("name").(string)),
		Type:    data.Get("type").(string),
		Content: data.Get("value").(string),
	}
	if attr, ok := data.GetOk("ttl"); ok {
		recordAttributes.TTL, _ = strconv.Atoi(attr.(string))
	}
	if attr, ok := data.GetOk("priority"); ok {
		recordAttributes.Priority, _ = strconv.Atoi(attr.(string))
	}

	log.Printf("[DEBUG] DNSimple Record create recordAttributes: %#v", recordAttributes)

	resp, err := provider.client.Zones.CreateRecord(ctx, provider.config.Account, data.Get("domain").(string), recordAttributes)
	if err != nil {
		return diag.Errorf("Failed to create DNSimple Record: %s", err)
	}

	data.SetId(strconv.FormatInt(resp.Data.ID, 10))
	log.Printf("[INFO] DNSimple Record ID: %s", data.Id())

	return resourceDNSimpleRecordRead(ctx, data, meta)
}

func resourceDNSimpleRecordRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deprecationWarning()
	provider := meta.(*Client)

	recordID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Record ID: %s", err)
	}

	resp, err := provider.client.Zones.GetRecord(ctx, provider.config.Account, data.Get("domain").(string), recordID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Record Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Couldn't find DNSimple Record: %s", err)
	}

	record := resp.Data
	data.Set("domain_id", record.ZoneID)
	data.Set("name", record.Name)
	data.Set("type", record.Type)
	data.Set("value", record.Content)
	data.Set("ttl", strconv.Itoa(record.TTL))
	data.Set("priority", strconv.Itoa(record.Priority))

	if record.Name == "" {
		data.Set("hostname", data.Get("domain").(string))
	} else {
		data.Set("hostname", fmt.Sprintf("%s.%s", record.Name, data.Get("domain").(string)))
	}

	return nil
}

func resourceDNSimpleRecordUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deprecationWarning()
	provider := meta.(*Client)

	recordID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Record ID: %s", err)
	}

	recordAttributes := dnsimple.ZoneRecordAttributes{}
	if attr, ok := data.GetOk("name"); ok {
		recordAttributes.Name = dnsimple.String(attr.(string))
	}
	if attr, ok := data.GetOk("type"); ok {
		recordAttributes.Type = attr.(string)
	}
	if attr, ok := data.GetOk("value"); ok {
		recordAttributes.Content = attr.(string)
	}
	if attr, ok := data.GetOk("ttl"); ok {
		recordAttributes.TTL, _ = strconv.Atoi(attr.(string))
	}
	if attr, ok := data.GetOk("priority"); ok {
		recordAttributes.Priority, _ = strconv.Atoi(attr.(string))
	}

	log.Printf("[DEBUG] DNSimple Record update configuration: %#v", recordAttributes)

	_, err = provider.client.Zones.UpdateRecord(ctx, provider.config.Account, data.Get("domain").(string), recordID, recordAttributes)
	if err != nil {
		return diag.Errorf("Failed to update DNSimple Record: %s", err)
	}

	return resourceDNSimpleRecordRead(ctx, data, meta)
}

func resourceDNSimpleRecordDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	deprecationWarning()
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple Record: %s, %s", data.Get("domain").(string), data.Id())

	recordID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Record ID: %s", err)
	}

	_, err = provider.client.Zones.DeleteRecord(ctx, provider.config.Account, data.Get("domain").(string), recordID)
	if err != nil {
		return diag.Errorf("Error deleting DNSimple Record: %s", err)
	}

	return nil
}

func resourceDNSimpleRecordImport(ctx context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	deprecationWarning()
	parts := strings.Split(data.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("error Importing dnsimple_record. Please make sure the record ID is in the form DOMAIN_RECORDID (i.e. example.com_1234)")
	}

	data.SetId(parts[1])
	data.Set("domain", parts[0])

	if err := resourceDNSimpleRecordRead(ctx, data, meta); err != nil {
		return nil, fmt.Errorf(err[0].Summary)
	}
	return []*schema.ResourceData{data}, nil
}
