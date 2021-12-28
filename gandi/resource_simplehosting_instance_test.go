package gandi

import (
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSimpleHostingInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigInstance(),
				SkipFunc: func() (bool, error) {
					url := os.Getenv("GANDI_URL")
					// The sandbox currently doesn't support the Simple Hosting API.
					// So, this test only has to be executed on the Gandi internal stage platform.
					return !strings.Contains(url, "stage"), nil
				},
			},
		},
	})
}

func testAccConfigInstance() string {
	return `
	  resource "gandi_simplehosting_instance" "create" {
	    name = "create"
	    size = "s+"
	    database_name = "mysql"
	    language_name = "php"
	    location = "FR"
          }
	`
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GANDI_KEY"); v == "" {
		t.Fatal("GANDI_KEY must be set for acceptance tests")
	}
}
