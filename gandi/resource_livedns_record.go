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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The FQDN of the domain",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The name of the record",
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The type of the record",
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "The TTL of the record",
			},
			"values": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
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
	return nil
}

func resourceLiveDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, name, recordType, err := expandRecordID(d.Id())
	record, err := client.GetDomainRecordWithNameAndType(zone, name, recordType)
	if err != nil {
		return err
	}
	d.Set("zone", zone)
	d.Set("name", record.RrsetName)
	d.Set("type", record.RrsetType)
	d.Set("ttl", record.RrsetTTL)
	d.Set("href", record.RrsetHref)
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
	_, err = client.ChangeDomainRecordWithNameAndType(zone, name, recordType, ttl, values)
	return err
}

func resourceLiveDNSRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	zone, name, recordType, err := expandRecordID(d.Id())

	if err != nil {
		return err
	}

	err = client.DeleteDomainRecord(zone, name, recordType)

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
