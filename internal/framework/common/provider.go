package common

import (
	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/utils"
)

type DnsimpleProviderConfig struct {
	Client *dnsimple.Client
	// TODO: Remove this once the official client supports the new endpoints
	TempClient      *utils.DNSimpleClient
	AccountID       string
	Prefetch        bool
	ZoneRecordCache ZoneRecordCache
}
