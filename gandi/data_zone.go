package gandi

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGandiZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceGandiZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*g.Gandi)
	name := d.Get("name").(string)
	log.Printf("[INFO] Reading Gandi zone '%s'", name)
	zones, err := client.ListZones()
	if err != nil {
		return err
	}
	var found *g.Zone
	for _, zone := range zones {
		if zone.Name == name {
			found = &zone
		}
	}
	if found == nil {
		return fmt.Errorf("Unknown zone with name : '%s'", name)
	}
	d.SetId(found.UUID)
	d.Set("name", found.Name)
	return nil
}
