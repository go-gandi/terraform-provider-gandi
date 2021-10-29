package gandi

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDNSSECKey() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name",
			},
			"algorithm": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "DNSSEC algorithm type",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNSSEC key type",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "DNSSEC public key",
			},
		},
		Create: resourceDNSSECKeyCreate,
		Delete: resourceDNSSECKeyDelete,
		Read:   resourceDNSSECKeyRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDNSSECKeyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("domain").(string)
	publicKey := d.Get("public_key").(string)

	request := domain.DNSSECKeyCreateRequest{
		Algorithm: d.Get("algorithm").(int),
		Type:      d.Get("type").(string),
		PublicKey: publicKey,
	}

	err = client.CreateDNSSECKey(resDomain, request)
	if err != nil {
		return
	}

	// Sent, got 202 response. What's now?
	time.Sleep(2 * time.Second)

	keys, err := client.ListDNSSECKeys(resDomain)
	if err != nil {
		return
	}

	for _, k := range keys {
		if k.PublicKey == publicKey {
			d.SetId(strconv.Itoa(k.ID))
			break
		}
	}

	return resourceDNSSECKeyRead(d, meta)
}

func resourceDNSSECKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	resDomain := d.Get("domain").(string)
	id := d.Id()
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		resDomain = parts[0]
		id = parts[1]
		if err = d.Set("id", id); err != nil {
			return fmt.Errorf("failed to set id for %s: %w", d.Id(), err)
		}
	}

	keys, err := client.ListDNSSECKeys(resDomain)
	if err != nil {
		return
	}

	var found domain.DNSSECKey
	var matchedKey bool = false
	for _, k := range keys {
		if strconv.Itoa(k.ID) == id {
			found = k
			matchedKey = true
			break
		}
	}
	if !matchedKey {
		err = fmt.Errorf("Cannot find DNSSEC key %s for domain %s", id, resDomain)
		return
	}

	if err = d.Set("algorithm", found.Algorithm); err != nil {
		return fmt.Errorf("failed to set algorithm for %s: %w", d.Id(), err)
	}
	if err = d.Set("type", found.Type); err != nil {
		return fmt.Errorf("failed to set type for %s: %w", d.Id(), err)
	}
	if err = d.Set("public_key", found.PublicKey); err != nil {
		return fmt.Errorf("failed to set public key for %s: %w", d.Id(), err)
	}
	if err = d.Set("domain", resDomain); err != nil {
		return fmt.Errorf("failed to set domain for %s: %w", d.Id(), err)
	}
	if err = d.Set("digest", found.Digest); err != nil {
		return fmt.Errorf("failed to set digest for %s: %w", d.Id(), err)
	}
	if err = d.Set("digest_type", found.DigestType); err != nil {
		return fmt.Errorf("failed to set digest_type for %s: %w", d.Id(), err)
	}
	if err = d.Set("keytag", found.KeyTag); err != nil {
		return fmt.Errorf("failed to set keytag for %s: %w", d.Id(), err)
	}
	return
}

func resourceDNSSECKeyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	domain := d.Get("domain").(string)
	id := d.Id()
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		domain = parts[0]
		id = parts[1]
		if err = d.Set("id", id); err != nil {
			return fmt.Errorf("failed to set id for %s: %w", d.Id(), err)
		}
	}

	err = client.DeleteDNSSECKey(domain, id)
	return
}
