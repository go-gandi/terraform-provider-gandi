package gandi

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlueRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlueRecordCreate,
		Read:          resourceGlueRecordRead,
		UpdateContext: resourceGlueRecordUpdate,
		DeleteContext: resourceGlueRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "Host name of the glue record",
			},
			"ips": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "List of IP addresses",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href of the record",
			},
			"fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fqdn of the record",
			},
			"fqdn_unicode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fqdn unicode of the record",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},
	}
}

func resourceGlueRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)
	name := d.Get("name").(string)

	var ips []string
	for _, i := range d.Get("ips").([]interface{}) {
		ips = append(ips, i.(string))
	}
	sort.Strings(ips)

	request := domain.GlueRecordCreateRequest{
		Name: name,
		IPs:  ips,
	}

	err := client.CreateGlueRecord(resDomain, request)
	if err != nil {
		return diag.Errorf("error creating instance: %s", err)
	}

	d.SetId(name)

	return diag.FromErr(resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		return resourceGlueRecordReadWithRetry(d, meta)
	}))
}

func resourceGlueRecordRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)

	id := d.Id()
	var found domain.GlueRecord
	found, err = client.GetGlueRecord(resDomain, id)
	if err != nil {
		return
	}

	if found.Name == "" {
		err = fmt.Errorf("cannot find Glue Record %s for zone %s", id, resDomain)
		return
	}

	if err = d.Set("name", found.Name); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("href", found.Href); err != nil {
		return fmt.Errorf("failed to set href for %s: %w", d.Id(), err)
	}
	if err = d.Set("ips", found.IPs); err != nil {
		return fmt.Errorf("failed to set ips for %s: %w", d.Id(), err)
	}
	if err = d.Set("fqdn", found.FQDN); err != nil {
		return fmt.Errorf("failed to set fqdn for %s: %w", d.Id(), err)
	}
	if err = d.Set("fqdn_unicode", found.FQDNUnicode); err != nil {
		return fmt.Errorf("failed to set fqdn unicode for %s: %w", d.Id(), err)
	}
	return
}

func resourceGlueRecordReadWithRetry(d *schema.ResourceData, meta interface{}) *resource.RetryError {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)
	id := d.Id()

	gluerecord, err := client.GetGlueRecord(resDomain, id)
	if err != nil {
		return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))
	}

	if gluerecord.Name == "" {
		return resource.RetryableError(fmt.Errorf("expected glue record to be created but was not found"))
	}

	err = resourceGlueRecordRead(d, meta)
	if err != nil {
		return resource.NonRetryableError(err)
	}
	return nil
}

func resourceGlueRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)
	id := d.Id()

	if d.HasChanges("ips") {
		var ips []string
		for _, i := range d.Get("ips").([]interface{}) {
			ips = append(ips, i.(string))
		}
		sort.Strings(ips)

		if err := client.UpdateGlueRecord(resDomain, id, ips); err != nil {
			return diag.FromErr(fmt.Errorf("failed to update ips for glue record at %s: %w", id, err))
		}
	}
	return diag.FromErr(resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		return resourceGlueRecordReadWithRetry(d, meta)
	}))
}

func resourceGlueRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients).Domain
	resDomain := d.Get("zone").(string)

	id := d.Id()

	return diag.FromErr(client.DeleteGlueRecord(resDomain, id))
}
