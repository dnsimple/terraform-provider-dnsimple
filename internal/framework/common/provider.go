package common

import (
	"github.com/dnsimple/dnsimple-go/v5/dnsimple"
)

type DnsimpleProviderConfig struct {
	Client          *dnsimple.Client
	AccountID       string
	Prefetch        bool
	ZoneRecordCache ZoneRecordCache
}
