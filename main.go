package main

import (
	"context"
	"github.com/autophagy/terraform-provider-cachix/cachix"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name cachix
func main() {
	providerserver.Serve(context.Background(), cachix.New, providerserver.ServeOpts{
		Address: "autophagy/cachix",
	})
}
