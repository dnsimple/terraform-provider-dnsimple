package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dnsimple/dnsimple-go/v4/dnsimple"
	"github.com/terraform-providers/terraform-provider-dnsimple/internal/consts"
)

func main() {
	token := os.Getenv("DNSIMPLE_TOKEN")
	account := os.Getenv("DNSIMPLE_ACCOUNT")
	sandbox := os.Getenv("DNSIMPLE_SANDBOX")

	if sandbox != "true" {
		panic("DNSIMPLE_SANDBOX must be set to true")
	}

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
			if err != nil && strings.Contains(err.Error(), "404") {
				// 404 is expected if the record was already deleted
			} else if err != nil {
				panic(err)
			}
		}
	}

	cancelAllContactChanges(context.Background(), dnsimpleClient, account)
	cleanupDomains(context.Background(), dnsimpleClient, account)
	cleanupEmailForwards(context.Background(), dnsimpleClient, account)
}

// RegistrantChangeCancelStates is a list of states that can be cancelled
var RegistrantChangeCancelStates = []string{
	consts.RegistrantChangeStateNew,
	consts.RegistrantChangeStatePending,
}

func cancelAllContactChanges(ctx context.Context, dnsimpleClient *dnsimple.Client, account string) {
	domainName := os.Getenv("DNSIMPLE_REGISTRANT_CHANGE_DOMAIN")

	if domainName == "" {
		fmt.Println("Skipping registrant change cleanup as DNSIMPLE_REGISTRANT_CHANGE_DOMAIN is not set")
		return
	}

	// Get the domain ID
	domainResponse, err := dnsimpleClient.Domains.GetDomain(ctx, account, domainName)
	if err != nil {
		panic(err)
	}

	listOptions := &dnsimple.RegistrantChangeListOptions{
		State: dnsimple.String(consts.RegistrantChangeStateNew),
		ListOptions: dnsimple.ListOptions{
			PerPage: dnsimple.Int(100),
		},
	}

	contactChanges, err := dnsimpleClient.Registrar.ListRegistrantChange(ctx, account, listOptions)
	if err != nil {
		panic(err)
	}

	for _, contactChange := range contactChanges.Data {
		if !contains(RegistrantChangeCancelStates, contactChange.State) {
			continue
		}

		if contactChange.DomainId != int(domainResponse.Data.ID) {
			continue
		}

		fmt.Printf("Cancelling registrant change for %s id=%d state=%s\n", domainName, contactChange.Id, contactChange.State)
		_, err := dnsimpleClient.Registrar.DeleteRegistrantChange(ctx, account, contactChange.Id)
		if err != nil {
			panic(err)
		}
	}

	if contactChanges.Pagination.TotalPages > 1 {
		for page := 2; page <= contactChanges.Pagination.TotalPages; page++ {
			listOptions := &dnsimple.RegistrantChangeListOptions{
				State: dnsimple.String(consts.RegistrantChangeStateNew),
				ListOptions: dnsimple.ListOptions{
					Page:    dnsimple.Int(page),
					PerPage: dnsimple.Int(100),
				},
			}

			contactChanges, err := dnsimpleClient.Registrar.ListRegistrantChange(ctx, account, listOptions)
			if err != nil {
				panic(err)
			}

			for _, contactChange := range contactChanges.Data {
				if !contains(RegistrantChangeCancelStates, contactChange.State) {
					continue
				}

				if contactChange.DomainId != int(domainResponse.Data.ID) {
					continue
				}

				fmt.Printf("Cancelling registrant change for %s id=%d state=%s\n", domainName, contactChange.Id, contactChange.State)
				_, err := dnsimpleClient.Registrar.DeleteRegistrantChange(ctx, account, contactChange.Id)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func cleanupDomains(ctx context.Context, dnsimpleClient *dnsimple.Client, account string) {
	cleanupEnabled := os.Getenv("DNSIMPLE_CLEANUP_DOMAINS")

	if cleanupEnabled != "true" {
		fmt.Println("Skipping domain cleanup as DNSIMPLE_CLEANUP_DOMAIN is not set")
		return
	}

	domainsToKeep := os.Getenv("DNSIMPLE_DOMAINS_TO_KEEP")
	var domainsToKeepList []string
	if domainsToKeep != "" {
		// Split the comma separated list of domains to keep
		domainsToKeepList = strings.Split(domainsToKeep, ",")
	}

	listOptions := &dnsimple.DomainListOptions{
		ListOptions: dnsimple.ListOptions{
			PerPage: dnsimple.Int(100),
		},
	}

	domains, err := dnsimpleClient.Domains.ListDomains(ctx, account, listOptions)
	if err != nil {
		panic(err)
	}

	for _, domain := range domains.Data {
		if contains(domainsToKeepList, domain.Name) {
			continue
		}

		fmt.Printf("Deleting domain %s\n", domain.Name)
		_, err := dnsimpleClient.Domains.DeleteDomain(ctx, account, domain.Name)
		if err != nil && strings.Contains(err.Error(), "The domain cannot be deleted because it is either being registered or is transferring in") {
			fmt.Printf("Skipping domain %s because it is being registered or is transferring in\n", domain.Name)
			continue
		} else if err != nil {
			panic(err)
		}
	}

	if domains.Pagination.TotalPages > 1 {
		for page := 2; page <= domains.Pagination.TotalPages; page++ {
			listOptions := &dnsimple.DomainListOptions{
				ListOptions: dnsimple.ListOptions{
					Page:    dnsimple.Int(page),
					PerPage: dnsimple.Int(100),
				},
			}

			domains, err := dnsimpleClient.Domains.ListDomains(ctx, account, listOptions)
			if err != nil {
				panic(err)
			}

			for _, domain := range domains.Data {
				if contains(domainsToKeepList, domain.Name) {
					continue
				}

				fmt.Printf("Deleting domain %s\n", domain.Name)
				_, err := dnsimpleClient.Domains.DeleteDomain(ctx, account, domain.Name)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func cleanupEmailForwards(ctx context.Context, dnsimpleClient *dnsimple.Client, account string) {
	emailForwardDomain := os.Getenv("DNSIMPLE_DOMAIN")
	if emailForwardDomain == "" {
		fmt.Println("Skipping email forward cleanup as DNSIMPLE_DOMAIN is not set")
		return
	}

	listOptions := &dnsimple.ListOptions{
		PerPage: dnsimple.Int(100),
	}

	emailForwards, err := dnsimpleClient.Domains.ListEmailForwards(ctx, account, emailForwardDomain, listOptions)
	if err != nil {
		panic(err)
	}

	for _, emailForward := range emailForwards.Data {
		fmt.Printf("Deleting email forward %s\n", emailForward.AliasName)
		_, err := dnsimpleClient.Domains.DeleteEmailForward(ctx, account, emailForwardDomain, emailForward.ID)
		if err != nil {
			panic(err)
		}
	}
}
