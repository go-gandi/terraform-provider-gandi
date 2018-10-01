package gandi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	g "github.com/tiramiseb/go-gandi-livedns"
)

func resourceZonerecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceZonerecordCreate,
		Read:   resourceZonerecordRead,
		Update: resourceZonerecordUpdate,
		Delete: resourceZonerecordDelete,
		Exists: resourceZonerecordExists,
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

func makeZonerecordID(zone, name, recordType string) string {
	return fmt.Sprintf("%s/%s/%s", zone, name, recordType)
}

func unmakeZonerecordID(id string) (zone, name, recordType string, err error) {
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

func resourceZonerecordCreate(d *schema.ResourceData, m interface{}) error {
	zoneUUID := d.Get("zone").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)
	ttl := d.Get("ttl").(int)
	valuesList := d.Get("values").(*schema.Set).List()
	var values []string
	for _, v := range valuesList {
		values = append(values, v.(string))
	}
	client := m.(*g.Gandi)
	_, err := client.CreateZoneRecord(zoneUUID, name, recordType, ttl, values)
	if err != nil {
		return err
	}
	calculatedID := fmt.Sprintf("%s/%s/%s", zoneUUID, name, recordType)
	d.SetId(calculatedID)
	return nil
}

func resourceZonerecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	zone, name, recordType, err := unmakeZonerecordID(d.Id())
	record, err := client.GetZoneRecordWithNameAndType(zone, name, recordType)
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

func resourceZonerecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	zone, name, recordType, err := unmakeZonerecordID(d.Id())

	if err != nil {
		return err
	}

	ttl := d.Get("ttl").(int)
	valuesList := d.Get("values").(*schema.Set).List()
	var values []string
	for _, v := range valuesList {
		values = append(values, v.(string))
	}
	_, err = client.ChangeZoneRecordWithNameAndType(zone, name, recordType, ttl, values)
	return err
}

func resourceZonerecordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*g.Gandi)
	zone, name, recordType, err := unmakeZonerecordID(d.Id())

	if err != nil {
		return err
	}

	err = client.DeleteZoneRecord(zone, name, recordType)

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceZonerecordExists(d *schema.ResourceData, m interface{}) (bool, error) {
	client := m.(*g.Gandi)
	zone, name, recordType, err := unmakeZonerecordID(d.Id())

	if err != nil {
		return false, err
	}

	_, err = client.GetZoneRecordWithNameAndType(zone, name, recordType)
	if err != nil {
		if strings.Index(err.Error(), "404") == 0 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
