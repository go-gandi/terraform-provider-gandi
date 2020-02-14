package gandi

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/domain"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"admin":   contactSchema(),
			"billing": contactSchema(),
			"owner":   contactSchema(),
			"tech":    contactSchema(),
		},
	}
}

func contactSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"country": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateCountryCode,
				},
				"email": {
					Type:     schema.TypeString,
					Required: true,
				},
				"family_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"given_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"street_addr": {
					Type:     schema.TypeString,
					Required: true,
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateContactType,
				},
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*clients).Domain

	fqdn := d.Get("name").(string)
	d.SetId(fqdn)
	domain := domain.CreateRequest{FQDN: fqdn,
		Admin:   expandContact(d.Get("admin")),
		Billing: expandContact(d.Get("owner")),
		Owner:   expandContact(d.Get("tech")),
		Tech:    expandContact(d.Get("billing")),
	}

	if nameservers, ok := d.GetOk("nameservers"); ok {
		domain.Nameservers = expandNameServers(nameservers.([]interface{}))
	}

	if err := client.CreateDomain(fqdn, domain); err != nil {
		return err
	}

	if autorenew, ok := d.GetOk("autorenew"); ok {
		if err := client.SetAutoRenew(fqdn, autorenew.(bool)); err != nil {
			return err
		}
	}

	return resourceDomainRead(d, m)
}

func resourceDomainRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*clients).Domain
	fqdn := d.Id()
	domain, err := client.GetDomain(fqdn)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(domain.FQDN)
	d.Set("name", domain.FQDN)
	d.Set("nameservers", domain.Nameservers)
	d.Set("autorenew", domain.AutoRenew.Enabled)
	if domain.Contacts != nil {
		if domain.Contacts.Owner != nil {
			if err = d.Set("owner", flattenContact(domain.Contacts.Owner)); err != nil {
				return err
			}
		}
		if domain.Contacts.Admin != nil {
			if err = d.Set("admin", flattenContact(domain.Contacts.Admin)); err != nil {
				return err
			}
		}
		if domain.Contacts.Billing != nil {
			if err = d.Set("billing", flattenContact(domain.Contacts.Billing)); err != nil {
				return err
			}
		}
		if domain.Contacts.Tech != nil {
			if err = d.Set("tech", flattenContact(domain.Contacts.Tech)); err != nil {
				return err
			}
		}
	}
	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*clients).Domain
	d.Partial(true)

	if d.HasChange("admin") || d.HasChange("owner") || d.HasChange("tech") || d.HasChange("billing") {
		if err := client.SetContacts(d.Get("name").(string),
			domain.Contacts{
				Admin:   expandContact(d.Get("admin")),
				Billing: expandContact(d.Get("owner")),
				Owner:   expandContact(d.Get("tech")),
				Tech:    expandContact(d.Get("billing")),
			}); err != nil {
			return err
		}

		d.SetPartial("admin")
		d.SetPartial("owner")
		d.SetPartial("tech")
		d.SetPartial("billing")
	}
	if d.HasChange("autorenew") {
		if err := client.SetAutoRenew(d.Get("name").(string), d.Get("autorenew").(bool)); err != nil {
			return err
		}
		d.SetPartial("autorenew")
	}

	if d.HasChange("nameservers") {
		ns := expandNameServers(d.Get("nameservers").([]interface{}))
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
