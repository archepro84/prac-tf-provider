package main

import (
	"context"
	"flag"
	"log"

	"main/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/archepro84/null",
		Debug:           debug,
		ProtocolVersion: 5,
	})

	if err != nil {
		log.Fatal(err)
	}
}
