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
				Config: testAccConfig(resourceName, domainName, ""),
			},
			{
				Config:            testAccConfig(resourceName, domainName, ""),
				ImportState:       true,
				ResourceName:      fmt.Sprintf("gandi_domain.%s", resourceName),
				ImportStateId:     domainName,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDomain_tags(t *testing.T) {
	id := uuid.New()
	domainName := fmt.Sprintf("terraform-provider-gandi-%s.com", id)
	resourceName := fmt.Sprintf("terraform_provider_gandi_%s_com", id)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		PreCheck:   func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(resourceName, domainName, `tags = ["tag1"]`),
			},
			{
				Config: testAccConfig(resourceName, domainName, `tags = ["tag2"]`),
			},
			{
				Config: testAccConfig(resourceName, domainName, ``),
			},
		},
	})
}

func testAccConfig(resourceName, domainName string, tags string) string {
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
	          family_name      = "lewo"
	          given_name       = "lewo"
	          mail_obfuscated  = false
	          phone            = "+33.606060606"
	          street_addr      = "Paris"
	          type             = "person"
	          zip              = "75000"
	      }
            %s
            }
	`, resourceName, domainName, tags)
}
