package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/shaharby7/showbiz/terraform/internal/provider"
)

func main() {
	if err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/showbiz-io/showbiz",
	}); err != nil {
		log.Fatal(err)
	}
}
