package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceDomainCreate,
		Read:   resourceDomainRead,
		Update: resourceDomainUpdate,
		Delete: resourceDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Deprecated:  "This nameservers attribute will be removed on next major release: the nameservers resource has to be used instead.\nSee https://github.com/go-gandi/terraform-provider-gandi/issues/88 for details.",
			},
			"autorenew": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Should the domain autorenew",
			},
			"owner":   contactSchema(true),
			"admin":   contactSchema(false),
			"billing": contactSchema(false),
			"tech":    contactSchema(false),
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func contactSchema(required bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Required: required,
		Optional: !required,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"country": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validateCountryCode,
					Description:  "The two letter country code for the contact",
				},
				"state": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The state code for the contact",
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
				"data_obfuscated": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Whether or not to obfuscate contact data in WHOIS",
				},
				"mail_obfuscated": {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Whether or not to obfuscate contact email in WHOIS",
				},
				"extra_parameters": {
					Type:        schema.TypeMap,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Extra parameters, needed for some jurisdictions",
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
		Owner: expandContact(d.Get("owner")),
	}

	if billing, ok := d.GetOk("billing"); ok {
		request.Billing = expandContact(billing)
	}
	if tech, ok := d.GetOk("tech"); ok {
		request.Tech = expandContact(tech)
	}
	if admin, ok := d.GetOk("admin"); ok {
		request.Admin = expandContact(admin)
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
	if err = d.Set("name", response.FQDN); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}

	// Nameservers are only set when livedns is not used. When
	// livedns is used, this nameservers list is managed by Gandi:
	// the user should not have to care about them.
	livedns, err := client.GetLiveDNS(fqdn)
	if err != nil {
		d.SetId("")
		return err
	}
	if livedns.Current != "livedns" {
		if err = d.Set("nameservers", response.Nameservers); err != nil {
			return fmt.Errorf("failed to set nameservers for %s: %w", d.Id(), err)
		}
	}
	if err = d.Set("autorenew", response.AutoRenew.Enabled); err != nil {
		return fmt.Errorf("failed to set autorenew for %s: %w", d.Id(), err)
	}
	if response.Contacts != nil {
		if response.Contacts.Owner != nil {
			if err = d.Set("owner", flattenContact(response.Contacts.Owner)); err != nil {
				return fmt.Errorf("failed to set the owner for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Admin != nil {
			if err = d.Set("admin", flattenContact(response.Contacts.Admin)); err != nil {
				return fmt.Errorf("failed to set the admin for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Billing != nil {
			if err = d.Set("billing", flattenContact(response.Contacts.Billing)); err != nil {
				return fmt.Errorf("failed to set the billing contact for %s: %w", d.Id(), err)
			}
		}
		if response.Contacts.Tech != nil {
			if err = d.Set("tech", flattenContact(response.Contacts.Tech)); err != nil {
				return fmt.Errorf("failed to set the tech contact for %s: %w", d.Id(), err)
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
