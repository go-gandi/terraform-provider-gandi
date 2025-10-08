package gandi

import (
	"github.com/go-gandi/go-gandi"
	"github.com/go-gandi/go-gandi/certificate"
	"github.com/go-gandi/go-gandi/config"
	"github.com/go-gandi/go-gandi/domain"
	"github.com/go-gandi/go-gandi/email"
	"github.com/go-gandi/go-gandi/livedns"
	"github.com/go-gandi/go-gandi/simplehosting"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider is the provider itself
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"personal_access_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_PERSONAL_ACCESS_TOKEN", nil),
				Description: "A Gandi API Personal Access Token",
				Sensitive:   true,
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_KEY", nil),
				Description: "(DEPRECATED) A Gandi API key",
				Deprecated:  "use personal_access_token instead",
				Sensitive:   true,
			},
			"sharing_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_SHARING_ID", nil),
				Description: "(DEPRECATED) A Gandi Sharing ID",
				Deprecated:  "use personal_access_token instead",
			},
			"dry_run": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Prevent the Domain provider from taking certain actions",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GANDI_URL", "https://api.gandi.net"),
				Description: "The Gandi API URL",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain":    dataSourceLiveDNSDomain(),
			"gandi_livedns_domain_ns": dataSourceLiveDNSDomainNS(),
			"gandi_livedns_record":    dataSourceLiveDNSRecord(),
			"gandi_domain":            dataSourceDomain(),
			"gandi_mailbox":           dataSourceMailbox(),
			"gandi_glue_record":       dataSourceGlueRecord(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gandi_livedns_domain":         resourceLiveDNSDomain(),
			"gandi_livedns_record":         resourceLiveDNSRecord(),
			"gandi_domain":                 resourceDomain(),
			"gandi_mailbox":                resourceMailbox(),
			"gandi_email_forwarding":       resourceEmailForwarding(),
			"gandi_dnssec_key":             resourceDNSSECKey(),
			"gandi_simplehosting_instance": resourceSimpleHostingInstance(),
			"gandi_glue_record":            resourceGlueRecord(),
			"gandi_simplehosting_vhost":    resourceSimpleHostingVhost(),
			"gandi_nameservers":            resourceNameservers(),
		},
		ConfigureFunc: getGandiClients,
	}
}

type clients struct {
	Domain        *domain.Domain
	Email         *email.Email
	LiveDNS       *livedns.LiveDNS
	SimpleHosting *simplehosting.SimpleHosting
	Certificate   *certificate.Certificate
}

func getGandiClients(d *schema.ResourceData) (interface{}, error) {
	config := config.Config{
		APIURL:              d.Get("url").(string),
		APIKey:              d.Get("key").(string),
		PersonalAccessToken: d.Get("personal_access_token").(string),
		SharingID:           d.Get("sharing_id").(string),
		DryRun:              d.Get("dry_run").(bool),
		Debug:               logging.IsDebugOrHigher(),
	}
	liveDNS := gandi.NewLiveDNSClient(config)
	email := gandi.NewEmailClient(config)
	domainClient := gandi.NewDomainClient(config)
	simpleHostingClient := gandi.NewSimpleHostingClient(config)
	certificateClient := gandi.NewCertificateClient(config)

	return &clients{
		Domain:        domainClient,
		Email:         email,
		LiveDNS:       liveDNS,
		SimpleHosting: simpleHostingClient,
		Certificate:   certificateClient,
	}, nil
}
