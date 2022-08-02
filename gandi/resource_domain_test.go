package gandi

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDomain_basic(t *testing.T) {
	id := uuid.New()
	domainName := fmt.Sprintf("terraform-provider-gandi-%s.com", id)
	resourceName := fmt.Sprintf("terraform_provider_gandi_%s_com", id)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		PreCheck:   func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigDomain(resourceName, domainName),
			},
			{
				Config:            testAccConfigDomain(resourceName, domainName),
				ImportState:       true,
				ResourceName:      fmt.Sprintf("gandi_domain.%s", resourceName),
				ImportStateId:     domainName,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigDomain(resourceName, domainName string) string {
	return fmt.Sprintf(`
	  resource "gandi_domain" "%s" {
	    name = "%s"
	      owner {
	          city             = "Paris"
	          country          = "FR"
	          data_obfuscated  = false
	          email            = "admin@example.com"
	          extra_parameters = {
	              "birth_city"               = ""
	              "birth_country"            = ""
	              "birth_date"               = ""
	              "birth_department"         = ""
	          }
	          family_name      = "Bar"
	          given_name       = "Foo"
	          mail_obfuscated  = false
	          phone            = "+33.606060606"
	          street_addr      = "Paris"
	          type             = "person"
	          zip              = "75000"
	      }
            }
	`, resourceName, domainName)
}
