package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/simplehosting"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSimpleHostingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSimpleHostingInstanceCreate,
		Read:   resourceSimpleHostingInstanceRead,
		Delete: resourceSimpleHostingInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the SimpleHosting instance",
			},
			"size": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The size of the SimpleHosting instance ('s+', 'm', 'l' or 'xxl')",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the database type ('mysql' or 'pgsql')",
			},
			"language_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the language ('php', 'python', 'nodejs' or 'ruby')",
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The datacenter location of the instance ('FR' or 'LU')",
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(5 * time.Minute)},
	}
}

func resourceSimpleHostingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	id := d.Id()
	found, err := client.GetInstance(id)

	if err != nil {
		return fmt.Errorf("unknown simplehosting instance '%s': %w", id, err)
	}
	d.SetId(found.ID)
	if err = d.Set("name", found.Name); err != nil {
		return fmt.Errorf("failed to set name for %s: %w", d.Id(), err)
	}
	if err = d.Set("size", found.Size); err != nil {
		return fmt.Errorf("failed to set size for %s: %w", d.Id(), err)
	}
	if err = d.Set("location", found.Datacenter.Region); err != nil {
		return fmt.Errorf("failed to set location for %s: %w", d.Id(), err)
	}
	if err = d.Set("database_name", found.Database.Name); err != nil {
		return fmt.Errorf("failed to set database_name for %s: %w", d.Id(), err)
	}
	if err = d.Set("language_name", found.Language.Name); err != nil {
		return fmt.Errorf("failed to set language_name for %s: %w", d.Id(), err)
	}
	return nil
}

func resourceSimpleHostingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	instanceId, err := client.CreateInstance(
		simplehosting.CreateInstanceRequest{
			Name:     d.Get("name").(string),
			Location: d.Get("location").(string),
			Size:     d.Get("size").(string),
			Type: &simplehosting.InstanceType{
				Database: &simplehosting.Database{
					Name: d.Get("database_name").(string),
				},
				Language: &simplehosting.Language{
					Name: d.Get("language_name").(string),
				},
			},
		},
	)
	if err != nil {
		return err
	}
	d.SetId(instanceId)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		client := meta.(*clients).SimpleHosting
		instance, err := client.GetInstance(instanceId)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error getting instance %s: %s", instanceId, err))
		}
		if instance.Status != "running" {
			return resource.RetryableError(fmt.Errorf("Expected instance %s to be running but was in state %s", instanceId, instance.Status))
		}
		return nil
	})
	if err != nil {
		return err
	}
	return resourceSimpleHostingInstanceRead(d, meta)
}

func resourceSimpleHostingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	instanceId := d.Id()
	_, err := client.DeleteInstance(instanceId)
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.GetInstance(instanceId)
		if err != nil {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("The instance %s have not been deleted yet", instanceId))
	})
}
