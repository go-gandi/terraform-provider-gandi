package gandi

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"

	"github.com/go-gandi/go-gandi"
	"github.com/go-gandi/go-gandi/config"
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

func deleteRecord() {
	config := config.Config{
		APIURL: os.Getenv("GANDI_URL"),
		APIKey: os.Getenv("GANDI_KEY"),
		Debug:  logging.IsDebugOrHigher(),
	}

	liveDNS := gandi.NewLiveDNSClient(config)
	err := liveDNS.DeleteDomainRecord(
		"terraform-provider-gandi.com",
		"www",
		"A")
	// To make golangci-lint happy :/
	if err != nil {
		return
	}
}

// TestAccRecord_manually_removed is a non regression test for
// https://github.com/go-gandi/terraform-provider-gandi/issues/100
// When a resource is manually ressource, Terraform has to recreate it.
func TestAccRecord_manually_removed(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		PreCheck:   func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigRecord(),
			},
			{
				// The record is removed. Terraform has to recreate it.
				PreConfig: deleteRecord,
				Config:    testAccConfigRecord(),
			},
		},
	})
}
