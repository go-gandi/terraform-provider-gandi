package gandi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataDomain_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDataDomain(),
			},
		},
	})
}

func testAccConfigDataDomain() string {
	return `
	  data "gandi_domain" "terraform_provider_gandi_com" {
	    name = "terraform-provider-gandi.com"
          }
	`
}
