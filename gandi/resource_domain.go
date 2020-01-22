package gandi

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nameservers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"autorenew": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	fqdn := d.Get("name").(string)
	d.SetId(fqdn)
	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*GandiClients).Domain
	fqdn := d.Get("name").(string)
	domain, err := client.GetDomain(fqdn)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(domain.ID)
	d.Set("name", domain.FQDN)
	d.Set("nameservers", domain.Nameservers)
	d.Set("autorenew", domain.AutoRenew)
	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*GandiClients).Domain
	d.Partial(true)
	if d.HasChange("autorenew") {
		if err := client.SetAutoRenew(d.Get("name").(string), d.Get("autorenew").(bool)); err != nil {
			return err
		}
		d.SetPartial("autorenew")
	}

	if d.HasChange("nameservers") {
		ns := ensureNameServers(d.Get("nameservers").([]interface{}))
		if err := client.UpdateNameServers(d.Get("name").(string), ns); err != nil {
			return err
		}
		d.SetPartial("nameservers")
	}
	d.Partial(false)
	return resourceDomainRead(d, m)
}

// The Gandi API doesn't presently support deleting domains
func resourceDomainDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}

func ensureNameServers(ns []interface{}) (ret []string) {
	for _, v := range ns {
		ret = append(ret, v.(string))
	}
	return
}
