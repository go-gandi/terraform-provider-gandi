package gandi

import (
	"sort"
	"strings"

	"github.com/go-gandi/go-gandi/email"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEmailForwarding() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name",
			},
			"source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account alias name",
			},
			"destination": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: "Forwards to email address",
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
	domain := d.Get("domain").(string)
	source := d.Get("source").(string)

	var destinations []string
	for _, i := range d.Get("destination").([]interface{}) {
		destinations = append(destinations, i.(string))
	}
	sort.Strings(destinations)

	request := email.CreateForwardRequest{
		Source:       source,
		Destinations: destinations,
	}

	err = client.CreateForward(domain, request)
	if err != nil {
		return
	}

	d.SetId(source + "@" + domain)

	return resourceEmailForwardingRead(d, meta)
}

func resourceEmailForwardingRead(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	forwards, err := client.GetForwards(domain)
	if err != nil {
		return
	}

	var response email.GetForwardRequest

	for _, found := range forwards {
		if found.Source == source {
			response = found
			break
		}
	}

	d.Set("href", response.Href)
	//sort.Strings(response.Destinations)
	d.Set("destination", response.Destinations)
	return
}

func resourceEmailForwardingUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	var destinations []string
	for _, i := range d.Get("destination").([]interface{}) {
		destinations = append(destinations, i.(string))
	}
	//sort.Strings(destinations)

	request := email.UpdateForwardRequest{
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

	err = client.DeleteForward(domain, source)
	return
}

func resourceEmailForwardingImport(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	client := meta.(*clients).Email
	source, domain := splitID(d.Id())

	forwards, err := client.GetForwards(domain)
	if err != nil {
		return
	}

	var response email.GetForwardRequest

	for _, found := range forwards {
		if found.Source == source {
			response = found
			break
		}
	}

	d.Set("href", response.Href)
	sort.Strings(response.Destinations)
	d.Set("destination", response.Destinations)

	data = []*schema.ResourceData{d}
	return
}
