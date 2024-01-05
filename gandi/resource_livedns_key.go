package gandi

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLiveDNSKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"flags": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DNSSEC key flags",
			},
			"algorithm": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "DNSSEC algorithm type",
			},
			"algorithm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNSSEC algorithm name",
			},
			"deleted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is the key deleted?",
			},
			"ds": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DS record as RFC1035 line",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the key",
			},
			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public key",
			},
			"tag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Tag",
			},
		},
		CreateContext: resourceLiveDNSKeyCreate,
		Delete:        resourceLiveDNSKeyDelete,
		Read:          resourceLiveDNSKeyRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceLiveDNSKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients).LiveDNS
	domain := d.Get("domain").(string)

	// The API does not return the key UUID. Not very convenient.
	// We will assume we will use the last key.
	last := ""
	keys, err := client.GetDomainKeys(domain)
	if err != nil {
		return diag.FromErr(err)
	}
	if len(keys) > 0 {
		last = keys[len(keys)-1].UUID
	}

	_, err = client.SignDomain(domain)
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		keys, err := client.GetDomainKeys(domain)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error getting keys: %s", err))
		}
		if len(keys) == 0 || keys[len(keys)-1].UUID == last {
			return resource.RetryableError(fmt.Errorf("expected domain key not found"))
		}
		key := keys[len(keys)-1]
		d.SetId(key.UUID)
		if err := d.Set("domain", domain); err != nil {
			return resource.NonRetryableError(
				fmt.Errorf("failed to set domain for %s/%s: %s", domain, key.UUID, err))
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(resourceLiveDNSKeyRead(d, meta))
}

func resourceLiveDNSKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).LiveDNS
	domain := d.Get("domain").(string)
	id := d.Id()
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		domain = parts[0]
		id = parts[1]
		d.SetId(id)
	}

	key, err := client.GetDomainKey(domain, id)
	if err != nil {
		return
	}

	if err := d.Set("domain", domain); err != nil {
		return fmt.Errorf("failed to set domain for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("flags", key.Flags); err != nil {
		return fmt.Errorf("failed to set flags for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("algorithm", key.Algorithm); err != nil {
		return fmt.Errorf("failed to set algorithm for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("algorithm_name", key.AlgorithmName); err != nil {
		return fmt.Errorf("failed to set algorithm name for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("deleted", key.Deleted); err != nil {
		return fmt.Errorf("failed to set deleted flag for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("ds", key.DS); err != nil {
		return fmt.Errorf("failed to set DS line for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("status", key.Status); err != nil {
		return fmt.Errorf("failed to set status for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("public_key", key.PublicKey); err != nil {
		return fmt.Errorf("failed to set status for %s/%s: %w", domain, id, err)
	}
	if err := d.Set("tag", key.Tag); err != nil {
		return fmt.Errorf("failed to set status for %s/%s: %w", domain, id, err)
	}
	return
}

func resourceLiveDNSKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).LiveDNS
	domain := d.Get("domain").(string)
	id := d.Id()
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		domain = parts[0]
		id = parts[1]
		d.SetId(id)
	}

	if err := client.DeleteDomainKey(domain, id); err != nil {
		return fmt.Errorf("failed to delete key %s/%s: %s", domain, id, err)
	}
	return nil
}
