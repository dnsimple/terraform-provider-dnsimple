package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	framework "github.com/terraform-providers/terraform-provider-dnsimple/internal/framework/provider"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./example/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//<remove to run>go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// version is the version of the provider.
var version = "dev"

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()
	serveOpts := providerserver.ServeOpts{
		Address: "registry.terraform.io/dnsimple/dnsimple",
		Debug:   debugMode,
	}

	providerserver.Serve(ctx, framework.New(version), serveOpts)

	if err != nil {
		log.Fatal(err)
	}
}
