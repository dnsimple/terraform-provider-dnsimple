package common_test

import (
	"testing"

	"github.com/dnsimple/dnsimple-go/v5/dnsimple"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/common"
)

func TestZoneRecordCache(t *testing.T) {
	t.Parallel()

	// Set up a test cache.
	cache := common.NewZoneRecordCache()

	// Prepare dnsimple.ZoneRecord set to add to cache.
	records := []dnsimple.ZoneRecord{
		{
			ID:      1,
			Name:    "a",
			Content: "1.2.3.4",
			Type:    "A",
		},
		{
			ID:      2,
			Name:    "a",
			Content: "1.2.3.5",
			Type:    "A",
		},
		{
			ID:      3,
			Name:    "a",
			Content: "mail.example.com",
			Type:    "MX",
		},
		{
			ID:      4,
			Name:    "b",
			Content: "1.2.3.4",
			Type:    "A",
		},
	}

	zoneName := "example.com"

	// Add zone to cache
	cache.Set(zoneName, records)

	// Get zone from cache
	records, ok := cache.Get(zoneName)
	assert.True(t, ok)

	// Find record in zone
	record, ok := cache.Find(zoneName, "a", "A", "1.2.3.4")
	assert.True(t, ok)
	assert.Equal(t, records[0], record)

	// Find does not find record in zone
	_, ok = cache.Find(zoneName, "b", "A", "1.2.3.5")
	assert.False(t, ok)
}
