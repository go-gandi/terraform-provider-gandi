package gandi

// A domain is always attached to a zone, so "delete" will not do anything

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
)

func resourceDomainattachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainattachmentCreate,
		Read:   resourceDomainattachmentRead,
		Delete: func(d *schema.ResourceData, m interface{}) error { return nil },
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDomainattachmentCreate(d *schema.ResourceData, m interface{}) error {
	domain := d.Get("domain").(string)
	zone := d.Get("zone").(string)
	client := m.(*g.Gandi)
	_, err := client.AttachDomainToZone(zone, domain)
	if err != nil {
		return err
	}
	d.SetId(domain)
	return nil
}

func resourceDomainattachmentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	domain, err := client.GetDomain(d.Id())
	if err != nil {
		return err
	}
	d.Set("domain", domain.FQDN)
	d.Set("zone", domain.ZoneUUID)
	return nil
}

func resourceDomainattachmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	err := client.DetachDomain(d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
