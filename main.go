package main

import (
	"context"
	"flag"
	"log"

	"github.com/FireWeave-HQ/terraform-provider-fireweave/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the Terraform provider documentation via `go generate`.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name fireweave

var (
	// version is set by the goreleaser ldflags at build time.
	version = "dev"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/FireWeave-HQ/fireweave",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
