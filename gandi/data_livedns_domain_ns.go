package gandi

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLiveDNSDomainNS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLiveDNSDomainNSRead,
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
	}
}

func dataSourceLiveDNSDomainNSRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	name := d.Get("name").(string)
	ns, err := client.GetDomainNS(name)
	if err != nil {
		ns = []string{}
	}

	d.SetId(name)
	if err = d.Set("nameservers", ns); err != nil {
		return fmt.Errorf("failed to set nameservers for %s: %w", d.Id(), err)
	}
	return nil
}
