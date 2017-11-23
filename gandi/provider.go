package gandi

import (
	"github.com/hashicorp/terraform/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
)

// Provider is the provider itself
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_KEY", nil),
				Description: "A Gandi LiveDNS API key",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"gandi_zone":             resourceZone(),
			"gandi_zonerecord":       resourceZonerecord(),
			"gandi_domainattachment": resourceDomainattachment(),
		},
		ConfigureFunc: getGandiClient,
	}
}

func getGandiClient(d *schema.ResourceData) (interface{}, error) {
	gandiClient := g.New(d.Get("key").(string))
	return gandiClient, nil
}
