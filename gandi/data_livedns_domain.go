package gandi

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLiveDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FQDN of the domain",
			},
		},
	}
}

func dataSourceLiveDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	name := d.Get("name").(string)
	found, err := client.GetDomain(name)
	if err != nil {
		return fmt.Errorf("unknown domain '%s': %w", name, err)
	}
	d.SetId(found.FQDN)
	if err = d.Set("name", found.FQDN); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	return nil
}
