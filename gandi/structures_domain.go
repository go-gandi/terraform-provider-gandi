package gandi

import (
	"fmt"
	"strings"

	"github.com/go-gandi/go-gandi/domain"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func flattenContact(in *domain.Contact) []interface{} {
	m := make(map[string]interface{})
	m["country"] = in.Country
	m["state"] = in.State
	m["mail_obfuscated"] = *in.MailObfuscated
	m["data_obfuscated"] = *in.DataObfuscated
	m["extra_parameters"] = in.ExtraParameters
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
	list := in.(*schema.Set).List()
	// We are sure the TypeSet contains a single element thanks to
	// the MaxItems=1 constraint
	contact := list[0].(map[string]interface{})

	dataObfuscated := contact["data_obfuscated"].(bool)
	mailObfuscated := contact["mail_obfuscated"].(bool)

	cnt := domain.Contact{
		Country:         contact["country"].(string),
		State:           contact["state"].(string),
		DataObfuscated:  &dataObfuscated,
		MailObfuscated:  &mailObfuscated,
		Email:           contact["email"].(string),
		FamilyName:      contact["family_name"].(string),
		GivenName:       contact["given_name"].(string),
		StreetAddr:      contact["street_addr"].(string),
		Phone:           contact["phone"].(string),
		City:            contact["city"].(string),
		OrgName:         contact["organisation"].(string),
		Zip:             contact["zip"].(string),
		ContactType:     expandContactType[contact["type"].(string)],
		ExtraParameters: contact["extra_parameters"].(map[string]interface{}),
	}
	return &cnt
}

func expandNameServers(ns []interface{}) (ret []string) {
	// We need to allocate at least 0 element. Otherwise, the
	// empty list is json encoded to null instead of [].
	// See https://apoorvam.github.io/blog/2017/golang-json-marshal-slice-as-empty-array-not-null/
	ret = make([]string, 0)
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
