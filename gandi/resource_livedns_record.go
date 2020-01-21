package gandi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/livedns"
)

func resourceLiveDNSRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceLiveDNSRecordCreate,
		Read:   resourceLiveDNSRecordRead,
		Update: resourceLiveDNSRecordUpdate,
		Delete: resourceLiveDNSRecordDelete,
		Exists: resourceLiveDNSRecordExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"values": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func makeRecordID(zone, name, recordType string) string {
	return fmt.Sprintf("%s/%s/%s", zone, name, recordType)
}

func explodeRecordID(id string) (zone, name, recordType string, err error) {
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

func resourceLiveDNSRecordCreate(d *schema.ResourceData, m interface{}) error {
	zoneUUID := d.Get("zone").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	ttl := d.Get("ttl").(int)
	valuesList := d.Get("values").(*schema.Set).List()
	var values []string
	for _, v := range valuesList {
		values = append(values, v.(string))
	}
	client := m.(*livedns.LiveDNS)
	_, err := client.CreateDomainRecord(zoneUUID, name, recordType, ttl, values)
	if err != nil {
		return err
	}
	calculatedID := fmt.Sprintf("%s/%s/%s", zoneUUID, name, recordType)
	d.SetId(calculatedID)
	return nil
}

func resourceLiveDNSRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*livedns.LiveDNS)
	zone, name, recordType, err := explodeRecordID(d.Id())
	record, err := client.GetDomainRecordWithNameAndType(zone, name, recordType)
	if err != nil {
		return err
	}
	d.Set("zone", zone)
	d.Set("name", record.RrsetName)
	d.Set("type", record.RrsetType)
	d.Set("ttl", record.RrsetTTL)
	d.Set("href", record.RrsetHref)
	d.Set("values", record.RrsetValues)
	return nil
}

func resourceLiveDNSRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*livedns.LiveDNS)
	zone, name, recordType, err := explodeRecordID(d.Id())

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

func resourceLiveDNSRecordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*livedns.LiveDNS)
	zone, name, recordType, err := explodeRecordID(d.Id())

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

func resourceLiveDNSRecordExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*livedns.LiveDNS)
	zone, name, recordType, err := explodeRecordID(d.Id())

	if err != nil {
		return false, err
	}

	_, err = client.GetDomainRecordWithNameAndType(zone, name, recordType)
	if err != nil {
		if strings.Index(err.Error(), "404") == 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
