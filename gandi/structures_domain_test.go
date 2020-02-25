package gandi

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tiramiseb/go-gandi/domain"
)

func TestValidateContactType(t *testing.T) {
	v := ""
	_, errors := validateContactType(v, "type")
	if len(errors) == 0 {
		t.Fatalf("The empty string is not a valid contact type")
	}

	valid := []string{"person", "company", "association", "public body", "reseller"}
	for _, v := range valid {
		_, errors := validateContactType(v, "type")
		if len(errors) != 0 {
			t.Fatalf("'%s' should be a valid contact type: %q", v, errors)
		}
	}
	invalid := []string{"p", "misc", "superhero"}
	for _, v := range invalid {
		_, errors := validateContactType(v, "type")
		if len(errors) == 0 {
			t.Fatalf("'%s' should not be a valid contact type", v)
		}
	}
}

func TestValidateCountryCode(t *testing.T) {
	v := ""
	_, errors := validateCountryCode(v, "country")
	if len(errors) == 0 {
		t.Fatalf("The empty string is not a valid country code")
	}

	_, errors = validateCountryCode("GB", "country")
	if len(errors) != 0 {
		t.Fatalf("'GB' should be a valid country code: %q", errors)
	}

	for _, v := range []string{"g", "gbb", "great britain"} {
		_, errors := validateCountryCode(v, "country")
		if len(errors) == 0 {
			t.Fatalf("%q is not a valid country code", v)
		}
	}
}

func TestRoundTripContactType(t *testing.T) {
	types := []string{"person", "company", "association", "public body", "reseller"}
	for _, v := range types {
		i := expandContactType[v]
		ret := flattenContactType[i]
		if ret != v {
			t.Errorf("Contact Type '%s' failed to roundtrip. Finalized as '%s'", v, ret)
		}
	}
	for i := 0; i < 5; i++ {
		val := flattenContactType[i]
		ret := expandContactType[val]
		if ret != i {
			t.Errorf("Contact val %d failed to roundtrip. Finalized as %d", i, ret)
		}
	}
}

func TestFlattenContact(t *testing.T) {
	type args struct {
		in *domain.Contact
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{name: "valid",
			args: args{
				in: &domain.Contact{
					ContactType: 1,
					Country:     "GB",
					Email:       "test@example.com",
					FamilyName:  "User",
					GivenName:   "Test",
					StreetAddr:  "1 Uncanny Valley",
					OrgName:     "Test Org",
					Phone:       "+1.2123333444",
					City:        "Libreville",
					Zip:         "12345",
				}},
			want: map[string]interface{}{
				"country":      "GB",
				"email":        "test@example.com",
				"family_name":  "User",
				"given_name":   "Test",
				"street_addr":  "1 Uncanny Valley",
				"type":         "company",
				"organisation": "Test Org",
				"phone":        "+1.2123333444",
				"city":         "Libreville",
				"zip":          "12345",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := flattenContact(tt.args.in)
			if !reflect.DeepEqual(got, []interface{}{tt.want}) {
				t.Errorf("flattenContact() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestExpandContact(t *testing.T) {
	want := &domain.Contact{
		ContactType: 1,
		Country:     "GB",
		Email:       "test@example.com",
		FamilyName:  "User",
		GivenName:   "Test",
		StreetAddr:  "1 Uncanny Valley",
		OrgName:     "Test Org",
		Phone:       "+1.2123333444",
		City:        "Libreville",
		Zip:         "12345",
	}

	s := contactSchema()
	contact := s.ZeroValue().(*schema.Set)

	contact.Add(map[string]interface{}{
		"country":      "GB",
		"email":        "test@example.com",
		"family_name":  "User",
		"given_name":   "Test",
		"street_addr":  "1 Uncanny Valley",
		"type":         "company",
		"organisation": "Test Org",
		"phone":        "+1.2123333444",
		"city":         "Libreville",
		"zip":          "12345",
	})

	if got := expandContact(contact); !reflect.DeepEqual(got, want) {
		t.Errorf("expandContact() = %#v, want %#v", got, want)
	}
}
