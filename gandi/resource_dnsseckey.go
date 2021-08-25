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
	res_domain := d.Get("domain").(string)
	public_key := d.Get("public_key").(string)

	request := domain.DNSSECKeyCreateRequest{
		Algorithm: d.Get("algorithm").(int),
		Type:      d.Get("type").(string),
		PublicKey: public_key,
	}

	err = client.CreateDNSSECKey(res_domain, request)
	if err != nil {
		return
	}

	// Sent, got 202 response. What's now?
	time.Sleep(2 * time.Second)

	keys, err := client.ListDNSSECKeys(res_domain)
	if err != nil {
		return
	}

	for _, k := range keys {
		if k.PublicKey == public_key {
			d.SetId(strconv.Itoa(k.ID))
			break
		}
	}

	return resourceDNSSECKeyRead(d, meta)
}

func resourceDNSSECKeyRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Domain
	res_domain := d.Get("domain").(string)
	id := d.Id()
	if strings.Contains(id, "/") {
		parts := strings.SplitN(id, "/", 2)
		res_domain = parts[0]
		id = parts[1]
		d.Set("id", id)
	}

	keys, err := client.ListDNSSECKeys(res_domain)
	if err != nil {
		return
	}

	var found domain.DNSSECKey
	var matched_key bool = false
	for _, k := range keys {
		if strconv.Itoa(k.ID) == id {
			found = k
			matched_key = true
			break
		}
	}
	if !matched_key {
		err = fmt.Errorf("Cannot find DNSSEC key %s for domain %s", id, res_domain)
		return
	}

	d.Set("algorithm", found.Algorithm)
	d.Set("type", found.Type)
	d.Set("public_key", found.PublicKey)
	d.Set("domain", res_domain)
	d.Set("digest", found.Digest)
	d.Set("digest_type", found.DigestType)
	d.Set("keytag", found.KeyTag)
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
		d.Set("id", id)
	}

	err = client.DeleteDNSSECKey(domain, id)
	return
}
