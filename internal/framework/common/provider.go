package common

import (
	"context"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DnsimpleProviderConfig struct {
	Client          *dnsimple.Client
	AccountID       string
	Prefetch        bool
	ZoneRecordCache ZoneRecordCache
}

type ZoneRecordCache map[string][]dnsimple.ZoneRecord

func (c ZoneRecordCache) Get(zoneName string) ([]dnsimple.ZoneRecord, bool) {
	records, ok := c[zoneName]
	return records, ok
}

func (c ZoneRecordCache) Set(zoneName string, records []dnsimple.ZoneRecord) {
	c[zoneName] = records
}

func (c ZoneRecordCache) Find(zoneName, recordName, recordType string) (dnsimple.ZoneRecord, bool) {
	records, ok := c.Get(zoneName)
	if !ok {
		return dnsimple.ZoneRecord{}, false
	}

	for _, record := range records {
		if record.Name == recordName && record.Type == recordType {
			return record, true
		}
	}

	return dnsimple.ZoneRecord{}, false
}

func (c ZoneRecordCache) Hydrate(ctx context.Context, client *dnsimple.Client, accountId string, zoneName string, options *dnsimple.ZoneRecordListOptions) error {
	if _, ok := c.Get(zoneName); !ok {
		var records []dnsimple.ZoneRecord

		if options == nil {
			options = &dnsimple.ZoneRecordListOptions{}
		}

		// Always use max page size
		options.PerPage = dnsimple.Int(100)
		// Fetch all records for the zone
		for {
			response, err := client.Zones.ListRecords(ctx, accountId, zoneName, options)
			if err != nil {
				return err
			}

			records = append(records, response.Data...)

			if response.Pagination.CurrentPage >= response.Pagination.TotalPages {
				break
			}

			options.Page = dnsimple.Int(response.Pagination.CurrentPage + 1)
		}

		c.Set(zoneName, records)
	}

	return nil
}
