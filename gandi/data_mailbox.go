package gandi

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMailbox() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"mailbox_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Mailbox ID",
			},
		},
		Read: dataSourceMailboxRead,
	}
}

func dataSourceMailboxRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	domain := d.Get("domain").(string)
	id := d.Get("mailbox_id").(string)

	found, err := client.GetMailbox(domain, id)
	if err != nil {
		return
	}

	d.SetId(id)
	d.Set("address", found.Address)
	d.Set("aliases", found.Aliases)
	d.Set("domain", found.Domain)
	d.Set("href", found.Href)
	d.Set("quota_used", found.QuotaUsed)
	d.Set("login", found.Login)
	d.Set("mailbox_type", found.MailboxType)
	return
}
