package gandi

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLiveDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceLiveDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Gandi zone '%s'", name)
	found, err := client.GetDomain(name)
	if err != nil {
		return fmt.Errorf("Unknown domain '%s': %w", name, err)
	}
	d.SetId(found.ZoneUUID)
	d.Set("name", found.FQDN)
	return nil
}
