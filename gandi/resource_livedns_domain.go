package gandi

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDNSDomainCreate,
		Read:   resourceLiveDNSDomainRead,
		Delete: resourceLiveDNSDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The FQDN of the domain",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The default TTL of the domain",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func resourceLiveDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)
	client := meta.(*clients).LiveDNS
	response, err := client.CreateDomain(name, ttl)
	if err != nil {
		return err
	}
	d.SetId(response.UUID)
	return nil
}

func resourceLiveDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, err := client.GetDomain(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", zone.FQDN)
	return nil
}

func resourceLiveDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
