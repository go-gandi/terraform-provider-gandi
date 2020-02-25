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
	m["phone"] = in.Phone
	m["city"] = in.City
	m["organisation"] = in.OrgName
	m["zip"] = in.Zip
	m["type"] = flattenContactType[in.ContactType]

	return []interface{}{m}
}

var expandContactType = map[string]int{
	"person":      0,
	"company":     1,
	"association": 2,
	"public body": 3,
	"reseller":    4,
}

var flattenContactType = []string{
	"person",
	"company",
	"association",
	"public body",
	"reseller",
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
		Phone:       contact["phone"].(string),
		City:        contact["city"].(string),
		OrgName:     contact["organisation"].(string),
		Zip:         contact["zip"].(string),
		ContactType: expandContactType[contact["type"].(string)],
	}
	return &cnt
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
