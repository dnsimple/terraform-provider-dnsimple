package dnsimple

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceDNSimpleZone() *schema.Resource {
	return &schema.Resource{
		ReadContext:   datasourceDNSimpleZoneRead,
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reverse": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

}

func datasourceDNSimpleZoneRead(_ context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	provider := meta.(*Client)

	resp, err := provider.client.Zones.GetZone(context.Background(), provider.config.Account, data.Get("name").(string))

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("DNSimple Zone Not Found - Refreshing from State")
			data.SetId("")
			return nil
		}
		return diag.Errorf("Couldn't find DNSimple Zone: %s", err)
	}

	zone := resp.Data

	data.SetId(strconv.FormatInt(zone.ID, 10))
	data.Set("account_id", zone.AccountID)
	data.Set("name", zone.Name)
	data.Set("reverse", zone.Reverse)

	return nil
}
