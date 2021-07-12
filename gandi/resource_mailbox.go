package gandi

import (
	"sort"
	"time"

	"github.com/go-gandi/go-gandi/email"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			State: schema.ImportStatePassthrough,
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

	d.Set("address", found.Address)
	d.Set("aliases", found.Aliases)
	d.Set("domain", found.Domain)
	d.Set("href", found.Href)
	d.Set("quota_used", found.QuotaUsed)
	d.Set("login", found.Login)
	d.Set("mailbox_type", found.MailboxType)
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

	err = client.UpdateEmail(domain, d.Id(), request)
	return
}

func resourceMailboxDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)

	err = client.DeleteEmail(domain, d.Id())
	return
}
