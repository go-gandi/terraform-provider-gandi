package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLiveDNSRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLiveDNSRecordRead,
		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FQDN of the domain",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the record",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of the record",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
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
				Computed:    true,
				Description: "A list of values of the record",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func dataSourceLiveDNSRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS

	zone := d.Get("zone").(string)
	name := d.Get("name").(string)
	recordType := d.Get("type").(string)

	record, err := client.GetDomainRecordByNameAndType(zone, name, recordType)

	if err != nil {
		requestError, ok := err.(*types.RequestError)
		if ok && requestError.StatusCode == 404 {
			return fmt.Errorf(
				"unknown '%s' '%s' record on domain '%s': %w",
				recordType, name, zone, err)
		}

		return fmt.Errorf(
			"failed to fetch '%s' '%s' record on domain '%s': %w",
			recordType, name, zone, err)
	}

	calculatedID := fmt.Sprintf("%s/%s/%s", zone, name, recordType)
	d.SetId(calculatedID)

	if err = d.Set("name", record.RrsetName); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("type", record.RrsetType); err != nil {
		return fmt.Errorf("failed to set type for %s: %w", d.Id(), err)
	}
	if err = d.Set("ttl", record.RrsetTTL); err != nil {
		return fmt.Errorf("failed to set ttl for %s: %w", d.Id(), err)
	}
	if err = d.Set("href", record.RrsetHref); err != nil {
		return fmt.Errorf("failed to set href for %s: %w", d.Id(), err)
	}
	if err = d.Set("values", record.RrsetValues); err != nil {
		return fmt.Errorf("failed to set the values for %s: %w", d.Id(), err)
	}

	return nil
}
