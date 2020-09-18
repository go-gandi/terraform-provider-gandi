package main

import (
	"github.com/go-gandi/terraform-provider-gandi/gandi"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return gandi.Provider()
		},
	})
}
