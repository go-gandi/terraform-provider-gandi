package gandi

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/domain"
)

func flattenContact(in *domain.Contact) []interface{} {
	m := make(map[string]interface{})
	m["country"] = in.Country
	m["email"] = in.Email
	m["family_name"] = in.FamilyName
	m["given_name"] = in.GivenName
	m["street_addr"] = in.StreetAddr
	m["type"] = flattenContactType(in.ContactType)

	return []interface{}{m}
}

func flattenContactType(cnt int) (ret string) {
	switch cnt {
	case 0:
		ret = "person"
	case 1:
		ret = "company"
	case 2:
		ret = "association"
	case 3:
		ret = "public body"
	case 4:
		ret = "reseller"
	}
	return
}

func expandContact(in interface{}) *domain.Contact {
	set := in.(*schema.Set)
	contact := set.List()[0].(map[string]interface{})
	cnt := domain.Contact{
		Country:     contact["country"].(string),
		Email:       contact["email"].(string),
		FamilyName:  contact["family_name"].(string),
		GivenName:   contact["given_name"].(string),
		StreetAddr:  contact["street_addr"].(string),
		ContactType: expandContactType(contact["type"].(string)),
	}
	return &cnt
}

func expandContactType(cnt string) (ret int) {
	switch cnt {
	case "person":
		ret = 0
	case "company":
		ret = 1
	case "association":
		ret = 2
	case "public body":
		ret = 3
	case "reseller":
		ret = 4
	}
	return
}

func expandNameServers(ns []interface{}) (ret []string) {
	for _, v := range ns {
		ret = append(ret, v.(string))
	}
	return
}

func validateContactType(val interface{}, key string) (warns []string, errs []error) {
	expected := val.(string)
	found := false
	types := []string{"person", "company", "association", "public body", "reseller"}
	for _, v := range types {
		if expected == v {
			found = true
		}
	}
	if !found {
		errs = append(errs, fmt.Errorf("%q must be one of %s. Got %s", key, strings.Join(types, ", "), expected))
	}
	return
}

func validateCountryCode(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if len(v) != 2 {
		errs = append(errs, fmt.Errorf("%q must be a two letter country code. Got %s", key, v))
	}
	return
}
