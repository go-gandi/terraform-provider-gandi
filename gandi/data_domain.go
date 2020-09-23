package gandi

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FQDN of the domain",
			},
			"nameservers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A list of nameservers for the domain",
			},
		},
		Read: dataSourceDomainRead,
	}
}

func dataSourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Domain
	name := d.Get("name").(string)
	found, err := client.GetDomain(name)
	if err != nil {
		return fmt.Errorf("Unknown domain '%s': %w", d.Id(), err)
	}
	d.SetId(found.FQDN)
	d.Set("name", found.FQDN)
	if err = d.Set("nameservers", found.Nameservers); err != nil {
		return fmt.Errorf("Failed to set nameservers for %s: %w", d.Id(), err)
	}
	return nil
}
