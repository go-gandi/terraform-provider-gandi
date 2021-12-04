package gandi

import (
	"fmt"
	"sort"
	"time"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGlueRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlueRecordCreate,
		Read:   resourceGlueRecordRead,
		Update: resourceGlueRecordUpdate,
		Delete: resourceGlueRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Host name of the glue record.",
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "List of IP addresses.",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func resourceGlueRecordCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)
	name := d.Get("name").(string)

	var values []string
	for _, i := range d.Get("values").([]interface{}) {
		values = append(values, i.(string))
	}
	sort.Strings(values)

	request := domain.GlueRecordCreateRequest{
		Name: name,
		IPs:  values,
	}

	err = client.CreateGlueRecord(resDomain, request)
	if err != nil {
		return
	}

	// Sent, got 202 response. Delay read.
	time.Sleep(2 * time.Second)

	d.SetId(name)

	return resourceGlueRecordRead(d, meta)
}

func resourceGlueRecordRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)

	id :=d.Id()
	records, err := client.ListGlueRecords(resDomain)
	if err != nil {
		return
	}

	var found domain.GlueRecord
	var matchedRecord bool = false
	for _, r := range records {
		if r.Name == id {
			found = r
			matchedRecord = true
			break
		}
	}

	if !matchedRecord {
		err = fmt.Errorf("Cannot find Glue Record %s for zone %s", id, resDomain)
		return
	}

	if err = d.Set("name", found.Name); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("href", found.Href); err != nil {
		return fmt.Errorf("failed to set href for %s: %w", d.Id(), err)
	}
	if err = d.Set("values", found.IPs); err != nil {
		return fmt.Errorf("failed to set values for %s: %w", d.Id(), err)
	}
	if err = d.Set("fqdn", found.FQDN); err != nil {
		return fmt.Errorf("failed to set fqdn for %s: %w", d.Id(), err)
	}
	if err = d.Set("fqdnUnicode", found.FQDNUnicode); err != nil {
		return fmt.Errorf("failed to set fqdn unicode for %s: %w", d.Id(), err)
	}
	return
}

func resourceGlueRecordUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)
	id :=d.Id()

	if d.HasChanges("values") {
		values := d.Get("values").([] string)

		if err := client.UpdateGlueRecord(resDomain, id, values); err != nil {
			return fmt.Errorf("failed to update values for glue record at %s: %w", id, err)
		}
	}

	return resourceGlueRecordRead(d, meta)
}

func resourceGlueRecordDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)

	id := d.Id()

	if err != nil {
		return err
	}

	if err = client.DeleteGlueRecord(resDomain, id); err != nil {
		return err
	}

	d.SetId("")
	return nil
}