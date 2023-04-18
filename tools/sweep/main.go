package main

import (
	"context"
	"os"

	"github.com/dnsimple/dnsimple-go/dnsimple"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
)

func main() {
	token := os.Getenv("DNSIMPLE_TOKEN")
	account := os.Getenv("DNSIMPLE_ACCOUNT")

	dnsimpleClient := dnsimple.NewClient(dnsimple.StaticTokenHTTPClient(context.Background(), token))
	dnsimpleClient.UserAgent = "terraform-provider-dnsimple/test"
	dnsimpleClient.BaseURL = consts.BaseURLSandbox

	domainName := os.Getenv("DNSIMPLE_DOMAIN")
	options := &dnsimple.ZoneRecordListOptions{
		ListOptions: dnsimple.ListOptions{
			PerPage: dnsimple.Int(100),
		},
	}
	records, err := dnsimpleClient.Zones.ListRecords(context.Background(), account, domainName, options)
	if err != nil {
		panic(err)
	}

	for _, record := range records.Data {
		if !record.SystemRecord {
			_, err := dnsimpleClient.Zones.DeleteRecord(context.Background(), account, domainName, record.ID)
			if err != nil {
				panic(err)
			}
		}
	}
}
