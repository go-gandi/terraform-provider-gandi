package gandi

import (
	"github.com/go-gandi/go-gandi"
	"github.com/go-gandi/go-gandi/domain"
	"github.com/go-gandi/go-gandi/livedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Provider is the provider itself
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_KEY", nil),
				Description: "A Gandi API key",
			},
			"sharing_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_SHARING_ID", nil),
				Description: "A Gandi Sharing ID",
			},
			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Prevent the Domain provider from taking certain actions",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain":    dataSourceLiveDNSDomain(),
			"gandi_livedns_domain_ns": dataSourceLiveDNSDomainNS(),
			"gandi_domain":            dataSourceDomain(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain": resourceLiveDNSDomain(),
			"gandi_livedns_record": resourceLiveDNSRecord(),
			"gandi_domain":         resourceDomain(),
		},
		ConfigureFunc: getGandiClients,
	}
}

type clients struct {
	Domain  *domain.Domain
	LiveDNS *livedns.LiveDNS
}

func getGandiClients(d *schema.ResourceData) (interface{}, error) {
	logging.SetOutput()

	config := gandi.Config{SharingID: d.Get("sharing_id").(string), DryRun: d.Get("dry_run").(bool)}
	liveDNS := gandi.NewLiveDNSClient(d.Get("key").(string), config)
	domainClient := gandi.NewDomainClient(d.Get("key").(string), config)

	return &clients{
		Domain:  domainClient,
		LiveDNS: liveDNS,
	}, nil
}
