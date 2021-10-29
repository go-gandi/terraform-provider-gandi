package gandi

import (
	"fmt"
	"sort"
	"strings"

	gandiemail "github.com/go-gandi/go-gandi/email"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEmailForwarding() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account alias name",
			},
			"destinations": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Forwards to email addresses",
			},
		},
		Create: resourceEmailForwardingCreate,
		Delete: resourceEmailForwardingDelete,
		Read:   resourceEmailForwardingRead,
		Update: resourceEmailForwardingUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceEmailForwardingImport,
		},
	}
}

func splitID(id string) (source, domain string) {
	parts := strings.SplitN(id, "@", 2)
	return parts[0], parts[1]
}

func resourceEmailForwardingCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source := d.Get("source").(string)
	email, domain := splitID(source)

	var destinations []string
	for _, i := range d.Get("destinations").([]interface{}) {
		destinations = append(destinations, i.(string))
	}
	sort.Strings(destinations)

	request := gandiemail.CreateForwardRequest{
		Source:       email,
		Destinations: destinations,
	}

	if err = client.CreateForward(domain, request); err != nil {
		return
	}

	d.SetId(email + "@" + domain)

	return resourceEmailForwardingRead(d, meta)
}

func resourceEmailForwardingRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	forwards, err := client.GetForwards(domain)
	if err != nil {
		return
	}

	var response gandiemail.GetForwardRequest

	for _, found := range forwards {
		if found.Source == source {
			response = found
			break
		}
	}

	if err = d.Set("href", response.Href); err != nil {
		return fmt.Errorf("failed to set href for %s: %s", d.Id(), err)
	}
	if err = d.Set("destinations", response.Destinations); err != nil {
		return fmt.Errorf("failed to set destination for %s: %s", d.Id(), err)
	}
	return
}

func resourceEmailForwardingUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	var destinations []string
	for _, i := range d.Get("destinations").([]interface{}) {
		destinations = append(destinations, i.(string))
	}

	request := gandiemail.UpdateForwardRequest{
		Destinations: destinations,
	}

	err = client.UpdateForward(domain, source, request)
	if err != nil {
		return
	}
	return resourceEmailForwardingRead(d, meta)
}

func resourceEmailForwardingDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	if err = client.DeleteForward(domain, source); err != nil {
		return
	}
	return
}

func resourceEmailForwardingImport(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	forwards, err := client.GetForwards(domain)
	if err != nil {
		return
	}

	var response gandiemail.GetForwardRequest

	for _, found := range forwards {
		if found.Source == source {
			response = found
			break
		}
	}

	if err = d.Set("href", response.Href); err != nil {
		return nil, fmt.Errorf("failed to set href for %s: %s", d.Id(), err)
	}
	sort.Strings(response.Destinations)
	if err = d.Set("destinations", response.Destinations); err != nil {
		return nil, fmt.Errorf("failed to set destinations for %s: %s", d.Id(), err)
	}

	data = []*schema.ResourceData{d}
	return
}
