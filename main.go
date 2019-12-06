package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-dnsimple/dnsimple"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: dnsimple.Provider})
}
