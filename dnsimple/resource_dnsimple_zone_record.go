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

func resourceDNSimpleZoneRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSimpleRecordCreate,
		ReadContext:   resourceDNSimpleRecordRead,
		UpdateContext: resourceDNSimpleRecordUpdate,
		DeleteContext: resourceDNSimpleRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSimpleRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"zone_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"qualified_name": {
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

func resourceDNSimpleRecordCreate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	resp, err := provider.client.Zones.CreateRecord(context.Background(), provider.config.Account, data.Get("zone_name").(string), recordAttributes)
	if err != nil {
		return diag.Errorf("Failed to create DNSimple Record: %s", err)
	}

	data.SetId(strconv.FormatInt(resp.Data.ID, 10))
	log.Printf("[INFO] DNSimple Record ID: %s", data.Id())

	return resourceDNSimpleRecordRead(nil, data, meta)
}

func resourceDNSimpleRecordRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	recordID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Record ID: %s", err)
	}

	resp, err := provider.client.Zones.GetRecord(context.Background(), provider.config.Account, data.Get("zone_name").(string), recordID)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Record Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Couldn't find DNSimple Record: %s", err)
	}

	record := resp.Data
	data.Set("zone_id", record.ZoneID)
	data.Set("name", record.Name)
	data.Set("type", record.Type)
	data.Set("value", record.Content)
	data.Set("ttl", strconv.Itoa(record.TTL))
	data.Set("priority", strconv.Itoa(record.Priority))

	if record.Name == "" {
		data.Set("qualified_name", data.Get("zone_name").(string))
	} else {
		data.Set("qualified_name", fmt.Sprintf("%s.%s", record.Name, data.Get("zone_name").(string)))
	}

	return nil
}

func resourceDNSimpleRecordUpdate(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, err = provider.client.Zones.UpdateRecord(context.Background(), provider.config.Account, data.Get("zone_name").(string), recordID, recordAttributes)
	if err != nil {
		return diag.Errorf("Failed to update DNSimple Record: %s", err)
	}

	return resourceDNSimpleRecordRead(nil, data, meta)
}

func resourceDNSimpleRecordDelete(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	log.Printf("[INFO] Deleting DNSimple Record: %s, %s", data.Get("zone_name").(string), data.Id())

	recordID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Error converting Record ID: %s", err)
	}

	_, err = provider.client.Zones.DeleteRecord(context.Background(), provider.config.Account, data.Get("zone_name").(string), recordID)
	if err != nil {
		return diag.Errorf("Error deleting DNSimple Record: %s", err)
	}

	return nil
}

func resourceDNSimpleRecordImport(_ context.Context, data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(data.Id(), "_")

	if len(parts) != 2 {
		return nil, fmt.Errorf("error Importing dnsimple_zone_record. Please make sure the record ID is in the form DOMAIN_RECORDID (i.e. example.com_1234)")
	}

	data.SetId(parts[1])
	data.Set("zone_name", parts[0])

	if err := resourceDNSimpleRecordRead(nil, data, meta); err != nil {
		return nil, fmt.Errorf(err[0].Summary)
	}
	return []*schema.ResourceData{data}, nil
}
