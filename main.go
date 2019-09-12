package main

import (
	kind "gbbirkisson/terraform-provider-kind/kind"

	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: kind.Provider,
	})
}
