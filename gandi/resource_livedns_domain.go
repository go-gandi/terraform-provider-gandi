package gandi

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/livedns"
)

func resourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDNSDomainCreate,
		Read:   resourceLiveDNSDomainRead,
		Delete: resourceLiveDNSDomainDelete,
		Exists: resourceLiveDNSDomainExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceLiveDNSDomainCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)
	client := m.(*livedns.LiveDNS)
	response, err := client.CreateDomain(name, ttl)
	if err != nil {
		return err
	}
	d.SetId(response.UUID)
	return nil
}

func resourceLiveDNSDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*livedns.LiveDNS)
	zone, err := client.GetDomain(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", zone.FQDN)
	return nil
}

func resourceLiveDNSDomainDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func resourceLiveDNSDomainExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*livedns.LiveDNS)
	_, err := client.GetDomain(d.Id())
	if err != nil {
		if strings.Index(err.Error(), "404") == 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
