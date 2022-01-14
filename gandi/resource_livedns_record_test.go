package gandi

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRecord_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		PreCheck:   func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRecord(),
			},
		},
	})
}

func testAccConfigRecord() string {
	return `
	  resource "gandi_livedns_record" "terraform_provider_gandi_com" {
	    zone = "terraform-provider-gandi.com"
            name = "www"
            type = "A"
            ttl = 3600
            values = ["192.168.0.1"]
          }
	`
}
