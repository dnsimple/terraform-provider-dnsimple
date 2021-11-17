package dnsimple

import (
	"context"
	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func fetchZoneRecords(provider *Client, accountId string, zoneName string, options *dnsimple.ZoneRecordListOptions) []dnsimple.ZoneRecord {

	if provider.cache[zoneName] == nil {
		records, err := provider.client.Zones.ListRecords(context.Background(), accountId, zoneName, options)
		if err != nil {
			return nil
		}
		provider.cache[zoneName] = records.Data
	}

	return provider.cache[zoneName]
}
