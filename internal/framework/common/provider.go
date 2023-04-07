package common

import (
	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DnsimpleProviderConfig struct {
	Client          *dnsimple.Client
	AccountID       string
	Prefetch        bool
	ZoneRecordCache ZoneRecordCache
}
