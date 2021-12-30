package gandi

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-gandi/go-gandi/email"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMailbox() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"login": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Login",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Password",
			},
			"mailbox_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "standard",
				Description: "Mailbox type",
			},
			"aliases": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Aliases for email",
			},
		},
		Create: resourceMailboxCreate,
		Delete: resourceMailboxDelete,
		Read:   resourceMailboxRead,
		Update: resourceMailboxUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMailboxCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)
	login := d.Get("login").(string)

	var aliases []string
	for _, i := range d.Get("aliases").([]interface{}) {
		aliases = append(aliases, i.(string))
	}
	sort.Strings(aliases)

	request := email.CreateEmailRequest{
		Aliases:     aliases,
		Login:       login,
		MailboxType: d.Get("mailbox_type").(string),
		Password:    d.Get("password").(string),
	}

	err = client.CreateEmail(domain, request)
	if err != nil {
		return
	}

	// Sent, got 202 response. What's now?
	time.Sleep(2 * time.Second)

	boxes, err := client.ListMailboxes(domain)
	if err != nil {
		return
	}

	for _, b := range boxes {
		if b.Login == login {
			d.SetId(b.ID)
			break
		}
	}

	return resourceMailboxRead(d, meta)
}

func resourceMailboxRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)

	found, err := client.GetMailbox(domain, d.Id())
	if err != nil {
		return
	}

	if err = d.Set("address", found.Address); err != nil {
		return fmt.Errorf("failed to set address for %s: %s", d.Id(), err)
	}
	if err = d.Set("aliases", found.Aliases); err != nil {
		return fmt.Errorf("failed to set aliases for %s: %s", d.Id(), err)
	}
	if err = d.Set("domain", found.Domain); err != nil {
		return fmt.Errorf("failed to set domain for %s: %s", d.Id(), err)
	}
	if err = d.Set("href", found.Href); err != nil {
		return fmt.Errorf("failed to set href for %s: %s", d.Id(), err)
	}
	if err = d.Set("quota_used", found.QuotaUsed); err != nil {
		return fmt.Errorf("failed to set quota_used for %s: %s", d.Id(), err)
	}
	if err = d.Set("login", found.Login); err != nil {
		return fmt.Errorf("failed to set login for %s: %s", d.Id(), err)
	}
	if err = d.Set("mailbox_type", found.MailboxType); err != nil {
		return fmt.Errorf("failed to set mailbox_type for %s: %s", d.Id(), err)
	}
	return
}

func resourceMailboxUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)

	var aliases []string
	for _, i := range d.Get("aliases").([]interface{}) {
		aliases = append(aliases, i.(string))
	}
	sort.Strings(aliases)

	request := email.UpdateEmailRequest{
		Aliases:  aliases,
		Login:    d.Get("login").(string),
		Password: d.Get("password").(string),
	}

	if err = client.UpdateEmail(domain, d.Id(), request); err != nil {
		return
	}
	return
}

func resourceMailboxDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)

	if err = client.DeleteEmail(domain, d.Id()); err != nil {
		return
	}

	return
}
