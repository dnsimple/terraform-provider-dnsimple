package dnsimple

import (
	"context"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

func fetchZoneRecords(ctx context.Context, provider *Client, accountId string, zoneName string, options *dnsimple.ZoneRecordListOptions) ([]dnsimple.ZoneRecord, error) {
	if provider.cache[zoneName] == nil {
		records, err := provider.client.Zones.ListRecords(ctx, accountId, zoneName, options)
		if err != nil {
			return nil, err
		}
		provider.cache[zoneName] = records.Data
	}

	return provider.cache[zoneName], nil
}
