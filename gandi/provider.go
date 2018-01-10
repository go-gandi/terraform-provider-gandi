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
			"sharing_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_SHARING_ID", nil),
				Description: "A Gandi LiveDNS sharing_id",
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
	gandiClient := g.New(d.Get("key").(string), d.Get("sharing_id").(string))
	return gandiClient, nil
}
