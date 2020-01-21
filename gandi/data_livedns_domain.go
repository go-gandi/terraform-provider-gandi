package gandi

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/livedns"
)

func dataSourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGandiLiveDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceGandiLiveDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*livedns.LiveDNS)
	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Gandi zone '%s'", name)
	found, err := client.GetDomain(name)
	if err != nil {
		return fmt.Errorf("Unknown domain with name : '%s'", name)
	}
	d.SetId(found.ZoneUUID)
	d.Set("name", found.FQDN)
	return nil
}
