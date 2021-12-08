package gandi

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGlueRecord() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host name of the glue record",
			},
			"ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "A list of the ip addresses provided for the glue record",
			},
		},
		Read: dataSourceGlueRecordRead,
	}
}

func dataSourceGlueRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Domain
	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	found, err := client.GetGlueRecord(zone, name)
	if err != nil {
		return fmt.Errorf("unknown domain '%s': %w", name, err)
	}
	d.SetId(found.Name)
	if err = d.Set("ips", found.IPs); err != nil {
		return fmt.Errorf("failed to set ips for %s: %w", d.Id(), err)
	}
	return nil
}
