package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/livedns"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLiveDNSDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDNSDomainCreate,
		Read:   resourceLiveDNSDomainRead,
		Update: resourceLiveDNSDomainUpdate,
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
			"automatic_snapshots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable the automatic creation of snapshots when records are changed",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func resourceLiveDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	ttl := d.Get("ttl").(int)
	client := meta.(*clients).LiveDNS
	_, err := client.CreateDomain(name, ttl)
	if err != nil {
		return err
	}
	d.SetId(name)
	if autosnap, ok := d.GetOk("automatic_snapshots"); ok {
		a := autosnap.(bool)
		if _, err := client.UpdateDomain(name, livedns.UpdateDomainRequest{AutomaticSnapshots: &a}); err != nil {
			return fmt.Errorf("failed to enable automatic snapshots for %s: %w", name, err)
		}
	}
	return resourceLiveDNSDomainRead(d, meta)
}

func resourceLiveDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, err := client.GetDomain(d.Id())
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(zone.FQDN)
	if err = d.Set("name", zone.FQDN); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("automatic_snapshots", zone.AutomaticSnapshots); err != nil {
		return fmt.Errorf("failed to set automatic_snapshots for %s: %w", d.Id(), err)
	}
	return nil
}

func resourceLiveDNSDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	name := d.Get("name").(string)

	if d.HasChange("automatic_snapshots") {
		a := d.Get("automatic_snapshots").(bool)
		if _, err := client.UpdateDomain(name, livedns.UpdateDomainRequest{AutomaticSnapshots: &a}); err != nil {
			return fmt.Errorf("failed to enable automatic snapshots for %s: %w", name, err)
		}
	}
	return resourceLiveDNSDomainRead(d, meta)
}

func resourceLiveDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
