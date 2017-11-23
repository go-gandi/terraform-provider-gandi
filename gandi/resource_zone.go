package gandi

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,
		Exists: resourceZoneExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := m.(*g.Gandi)
	response, err := client.CreateZone(name)
	if err != nil {
		return err
	}
	d.SetId(response.UUID)
	return nil
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	zone, err := client.GetZone(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", zone.Name)
	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	name := d.Get("name").(string)
	_, err := client.UpdateZone(d.Id(), name)
	return err
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	err := client.DeleteZone(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceZoneExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*g.Gandi)
	_, err := client.GetZone(d.Id())
	if err != nil {
		if strings.Index(err.Error(), "404") == 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
