package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/dnsimple/dnsimple-go/v7/dnsimple"
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

	cancelContactChanges(context.Background(), dnsimpleClient, account)
	cleanupDomains(context.Background(), dnsimpleClient, account)
}

// RegistrantChangeCancelStates is a list of states that can be cancelled
var RegistrantChangeCancelStates = []string{
	consts.RegistrantChangeStateNew,
	consts.RegistrantChangeStatePending,
}

func cancelContactChanges(ctx context.Context, dnsimpleClient *dnsimple.Client, account string) {
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
		fmt.Println("Skipping domain cleanup as DNSIMPLE_CLEANUP_DOMAINS is not set")
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
