package gandi

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLiveDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDNSRecordCreate,
		Read:   resourceLiveDNSRecordRead,
		Update: resourceLiveDNSRecordUpdate,
		Delete: resourceLiveDNSRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The FQDN of the domain",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the record",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the record",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The TTL of the record",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href of the record",
			},
			"values": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "A list of values of the record",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func expandRecordID(id string) (zone, name, recordType string, err error) {
	splitID := strings.Split(id, "/")

	if len(splitID) != 3 {
		err = errors.New("Id format should be '{zone_id}/{record_name}/{record_type}'")
		return
	}

	zone = splitID[0]
	name = splitID[1]
	recordType = splitID[2]
	return
}

func resourceLiveDNSRecordCreate(d *schema.ResourceData, meta interface{}) error {
	zoneUUID := d.Get("zone").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	ttl := d.Get("ttl").(int)
	valuesList := d.Get("values").(*schema.Set).List()
	var values []string
	for _, v := range valuesList {
		values = append(values, v.(string))
	}
	client := meta.(*clients).LiveDNS
	_, err := client.CreateDomainRecord(zoneUUID, name, recordType, ttl, values)
	if err != nil {
		return err
	}
	calculatedID := fmt.Sprintf("%s/%s/%s", zoneUUID, name, recordType)
	d.SetId(calculatedID)
	return resourceLiveDNSRecordRead(d, meta)
}

func resourceLiveDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, name, recordType, err := expandRecordID(d.Id())
	if err != nil {
		return err
	}
	record, err := client.GetDomainRecordByNameAndType(zone, name, recordType)
	if err != nil {
		return err
	}
	if err = d.Set("zone", zone); err != nil {
		return fmt.Errorf("Failed to set zone for %s: %w", d.Id(), err)
	}
	if err = d.Set("name", record.RrsetName); err != nil {
		return fmt.Errorf("Failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("type", record.RrsetType); err != nil {
		return fmt.Errorf("Failed to set type for %s: %w", d.Id(), err)
	}
	if err = d.Set("ttl", record.RrsetTTL); err != nil {
		return fmt.Errorf("Failed to set ttl for %s: %w", d.Id(), err)
	}
	if err = d.Set("href", record.RrsetHref); err != nil {
		return fmt.Errorf("Failed to set href for %s: %w", d.Id(), err)
	}
	if err = d.Set("values", record.RrsetValues); err != nil {
		return fmt.Errorf("Failed to set the values for %s: %w", d.Id(), err)
	}
	return nil
}

func resourceLiveDNSRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, name, recordType, err := expandRecordID(d.Id())

	if err != nil {
		return err
	}

	ttl := d.Get("ttl").(int)
	valuesList := d.Get("values").(*schema.Set).List()
	var values []string
	for _, v := range valuesList {
		values = append(values, v.(string))
	}
	_, err = client.UpdateDomainRecordByNameAndType(zone, name, recordType, ttl, values)
	if err != nil {
		return err
	}
	return resourceLiveDNSRecordRead(d, meta)
}

func resourceLiveDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, name, recordType, err := expandRecordID(d.Id())

	if err != nil {
		return err
	}

	if err = client.DeleteDomainRecord(zone, name, recordType); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
