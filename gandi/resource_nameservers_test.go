package gandi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNameservers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		PreCheck:   func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigNameservers(),
			},
		},
	})
}

func testAccConfigNameservers() string {
	return `
	  resource "gandi_nameservers" "terraform_provider_gandi_com" {
	    domain = "terraform-provider-gandi.com"
            nameservers = ["ns1.example.foo", "ns2.example.foo"]
          }
	`
}
