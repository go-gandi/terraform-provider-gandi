package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/tiramiseb/terraform-provider-gandi/gandi"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return gandi.Provider()
		},
	})
}
