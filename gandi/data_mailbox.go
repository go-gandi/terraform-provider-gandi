package gandi

import (
	"fmt"

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
	return nil
}
