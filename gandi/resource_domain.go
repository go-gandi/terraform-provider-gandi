package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FQDN of the domain",
			},
			"nameservers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "A list of nameservers for the domain",
			},
			"autorenew": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Should the domain autorenew",
			},
			"admin":   contactSchema(),
			"billing": contactSchema(),
			"owner":   contactSchema(),
			"tech":    contactSchema(),
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
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
					Description:  "The two letter country code for the contact",
				},
				"email": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Contact email address",
				},
				"family_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Family name of the contact",
				},
				"given_name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Given name of the contact",
				},
				"street_addr": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Street Address of the contact",
				},
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateContactType,
					Description:  "One of 'person', 'company', 'association', 'public body', or 'reseller'",
				},
				"phone": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Phone number for the contact",
				},
				"city": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "City for the contact",
				},
				"organisation": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The legal name of the organisation. Required for types other than person",
				},
				"zip": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Postal Code/Zipcode of the contact",
				},
			},
		},
	}
}

func resourceDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Domain

	fqdn := d.Get("name").(string)
	d.SetId(fqdn)
	request := domain.CreateRequest{FQDN: fqdn,
		Admin:   expandContact(d.Get("admin")),
		Billing: expandContact(d.Get("billing")),
		Owner:   expandContact(d.Get("owner")),
		Tech:    expandContact(d.Get("tech")),
	}

	if nameservers, ok := d.GetOk("nameservers"); ok {
		request.Nameservers = expandNameServers(nameservers.([]interface{}))
	}

	if err := client.CreateDomain(request); err != nil {
		return err
	}

	if autorenew, ok := d.GetOk("autorenew"); ok {
		if err := client.SetAutoRenew(fqdn, autorenew.(bool)); err != nil {
			return err
		}
	}

	return resourceDomainRead(d, meta)
}

func resourceDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Domain
	fqdn := d.Id()
	response, err := client.GetDomain(fqdn)
	if err != nil {
		d.SetId("")
		return err
	}
	d.SetId(response.FQDN)
	d.Set("name", response.FQDN)
	if err = d.Set("nameservers", response.Nameservers); err != nil {
		return fmt.Errorf("Failed to set nameservers for %s: %w", d.Id(), err)
	}
	d.Set("autorenew", response.AutoRenew.Enabled)
	if response.Contacts != nil {
		if response.Contacts.Owner != nil {
			if err = d.Set("owner", flattenContact(response.Contacts.Owner)); err != nil {
				return fmt.Errorf("Failed to set the owner for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Admin != nil {
			if err = d.Set("admin", flattenContact(response.Contacts.Admin)); err != nil {
				return fmt.Errorf("Failed to set the admin for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Billing != nil {
			if err = d.Set("billing", flattenContact(response.Contacts.Billing)); err != nil {
				return fmt.Errorf("Failed to set the billing contact for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Tech != nil {
			if err = d.Set("tech", flattenContact(response.Contacts.Tech)); err != nil {
				return fmt.Errorf("Failed to set the tech contact for %s: %w", d.Id(), err)
			}
		}
	}
	return nil
}

func resourceDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Domain

	if d.HasChanges("admin", "owner", "tech", "billing") {
		if err := client.SetContacts(d.Get("name").(string),
			domain.Contacts{
				Admin:   expandContact(d.Get("admin")),
				Billing: expandContact(d.Get("billing")),
				Owner:   expandContact(d.Get("owner")),
				Tech:    expandContact(d.Get("tech")),
			}); err != nil {
			return err
		}

	}
	if d.HasChange("autorenew") {
		if err := client.SetAutoRenew(d.Get("name").(string), d.Get("autorenew").(bool)); err != nil {
			return err
		}
	}

	if d.HasChange("nameservers") {
		ns := expandNameServers(d.Get("nameservers").([]interface{}))
		if err := client.UpdateNameServers(d.Get("name").(string), ns); err != nil {
			return err
		}
	}
	return resourceDomainRead(d, meta)
}

// The Gandi API doesn't presently support deleting domains
func resourceDomainDelete(d *schema.ResourceData, _ interface{}) error {
	d.SetId("")
	return nil
}
