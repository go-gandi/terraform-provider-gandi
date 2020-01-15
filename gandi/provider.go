package gandi

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
	"github.com/tiramiseb/go-gandi-livedns/gandi_config"
	"github.com/tiramiseb/go-gandi-livedns/gandi_domain"
	"github.com/tiramiseb/go-gandi-livedns/gandi_livedns"
)

// Provider is the provider itself
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_KEY", nil),
				Description: "A Gandi API key",
			},
			"sharing_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_SHARING_ID", nil),
				Description: "A Gandi Sharing ID",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain": dataSourceLiveDNSDomain(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain":             resourceLiveDNSDomain(),
			"gandi_livedns_record":       resourceLiveDNSRecord(),
		},
		ConfigureFunc: getGandiClient,
	}
}

type GandiClients struct {
	Domain  *gandi_domain.Domain
	LiveDNS *gandi_livedns.LiveDNS
}

func getGandiClient(d *schema.ResourceData) (interface{}, error) {
	config := gandi_config.Config{SharingID: d.Get("sharing_id").(string)}
	liveDNS := g.NewLiveDNSClient(d.Get("key").(string), config)
	domain := g.NewDomainClient(d.Get("key").(string), config)

	return &GandiClients{
		Domain:  domain,
		LiveDNS: liveDNS,
	}, nil
}
