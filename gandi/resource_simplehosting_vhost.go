package gandi

import (
	"fmt"
	"time"

	"github.com/go-gandi/go-gandi/certificate"
	"github.com/go-gandi/go-gandi/simplehosting"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSimpleHostingVhost() *schema.Resource {
	return &schema.Resource{
		Create: resourceSimpleHostingVhostCreate,
		Read:   resourceSimpleHostingVhostRead,
		Delete: resourceSimpleHostingVhostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the SimpleHosting instance",
				ForceNew:    true,
			},
			"fqdn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The FQDN of the Vhost",
				ForceNew:    true,
			},
			"linked_dns_zone_alteration": {
				Type:        schema.TypeBool,
				Description: "Whether to alter the linked DNS Zone",
				ForceNew:    true,
				Optional:    true,
				Default:     true,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the created free certificate",
			},
			"application": {
				Type:        schema.TypeString,
				Description: "The name of an application to install ('grav', 'matomo', 'nextcloud', 'prestashop', 'wordpress')",
				ForceNew:    true,
				Optional:    true,
				Default:     true,
			},
		},
		Timeouts: &schema.ResourceTimeout{Default: schema.DefaultTimeout(5 * time.Minute)},
	}
}

// freeCertificateCreate creates a free certificate for the Vhost.
func freeCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).Certificate
	response, err := client.CreateCertificate(
		certificate.CreateCertificateRequest{
			CN:      d.Get("fqdn").(string),
			Package: "cert_free_1_0_0",
		},
	)
	if err != nil {
		return err
	}
	certificateId := response.ID
	err = d.Set("certificate_id", certificateId)
	if err != nil {
		return err
	}
	return nil
}

func freeCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Certificate
	certificateId := d.Get("certificate_id").(string)
	if certificateId == "" {
		return nil
	}
	_, err = client.DeleteCertificate(certificateId)
	return
}

func resourceSimpleHostingVhostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	instanceId := d.Get("instance_id").(string)
	fqdn := d.Get("fqdn").(string)
	request := simplehosting.CreateVhostRequest{
		FQDN: fqdn,
	}
	if d.Get("linked_dns_zone_alteration").(bool) {
		request.LinkedDNSZone = &simplehosting.LinkedDNSZoneRequest{
			AllowAlteration:        true,
			AlowAlterationOverride: true,
		}
	}
	_, err := client.CreateVhost(
		instanceId,
		request,
	)
	if err != nil {
		return err
	}
	d.SetId(fqdn)

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		instance, err := client.GetVhost(instanceId, fqdn)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error getting vhost %s of instance %s: %w", instanceId, fqdn, err))
		}

		if instance.Status != "running" {
			return resource.RetryableError(fmt.Errorf("Expected vhost %s of instance %s to be running but was in state %s", instanceId, fqdn, instance.Status))
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Note it is plan to let the SimpleHosting API manage the
	// certificate. For now, we have to create in the Terraform
	// provider.
	err = freeCertificateCreate(d, meta)
	if err != nil {
		return err
	}

	// Unfortunately, it is not possible to set an Application on
	// the Vhost creation: we then have to update the vhost once
	// created:/
	applicationName := d.Get("application").(string)
	if applicationName != "" {
		_, err = client.UpdateVhost(
			instanceId,
			fqdn,
			simplehosting.PatchVhostRequest{
				Application: &simplehosting.Application{
					Name: applicationName,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return resourceSimpleHostingVhostRead(d, meta)
}

func resourceSimpleHostingVhostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	instanceId := d.Get("instance_id").(string)
	fqdn := d.Get("fqdn").(string)
	found, err := client.GetVhost(instanceId, fqdn)

	if err != nil {
		return fmt.Errorf("unknown simplehosting vhost '%s' of instance '%s': %w", instanceId, fqdn, err)
	}

	d.SetId(found.FQDN)
	if err = d.Set("linked_dns_zone_alteration", found.LinkedDNSZone.AllowAlteration); err != nil {
		return fmt.Errorf("failed to set linked_dns_zone_alteration for %s: %w", d.Id(), err)
	}
	if found.Application != nil {
		if err = d.Set("application", found.Application.Name); err != nil {
			return fmt.Errorf("failed to set application for %s: %w", d.Id(), err)
		}
	}
	return nil
}

func resourceSimpleHostingVhostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients).SimpleHosting
	instanceId := d.Get("instance_id").(string)
	fqdn := d.Get("fqdn").(string)
	_, err := client.DeleteVhost(instanceId, fqdn)
	if err != nil {
		return err
	}

	// Note it is plan to let the SimpleHosting API manage the
	// certificate. For now, we try to delete it but don't care
	// about potential error.
	_ = freeCertificateDelete(d, meta)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := client.GetVhost(instanceId, fqdn)
		if err != nil {
			return nil
		}
		// We should check the return code is 404 but this is
		// currently not provided by the go-gandi client
		// library
		return resource.RetryableError(fmt.Errorf("The vhost %s of instance %s have not been deleted yet", fqdn, instanceId))
	})
}
