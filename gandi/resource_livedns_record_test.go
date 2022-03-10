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

func TestKeepUniqueRecords(t *testing.T) {
	recordsListFromApi := []string{"record_one", "record_two", "record_three"}
	recordsListFromTerraform := []string{"record_one", "tf_record_two"}

	t.Run("Remove duplicated record from records list", func(t *testing.T) {
		recordsList := append(recordsListFromApi, recordsListFromTerraform...)
		shortenedList := keepUniqueRecords(recordsList)

		if !(len(shortenedList) == 4) {
			t.Errorf("Amount of records should have been decreased by one.")
		}
	})
}

func TestIfRecordISWrappedWithQuotes(t *testing.T) {
	t.Run("wrapped with quotes", func(t *testing.T) {
		wrappedRecord := "\"192.168.0.1\""
		if !isRecordWrappedWithQuotes(wrappedRecord) {
			t.Errorf("%s record is wrapped with quotes.", wrappedRecord)
		}
	})

	t.Run("suffix quote", func(t *testing.T) {
		wrappedRecord := "192.168.0.1\""
		if isRecordWrappedWithQuotes(wrappedRecord) {
			t.Errorf("%s record is not wrapped with quotes.", wrappedRecord)
		}
	})

	t.Run("prefix quote", func(t *testing.T) {
		wrappedRecord := "\"192.168.0.1"
		if isRecordWrappedWithQuotes(wrappedRecord) {
			t.Errorf("%s record is not wrapped with quotes.", wrappedRecord)
		}
	})

	t.Run("no quotes", func(t *testing.T) {
		wrappedRecord := "192.168.0.1"
		if isRecordWrappedWithQuotes(wrappedRecord) {
			t.Errorf("%s record is not wrapped with quotes.", wrappedRecord)
		}
	})
}

func TestWrappingRecordsWithQuotes(t *testing.T) {
	t.Run("wrapped with quotes", func(t *testing.T) {
		records := []string{"\"192.168.0.1\"", "192.168.0.2", "192.168.0.3", "\"192.168.0.1\""}
		wrappedRecords := wrapRecordsWithQuotes(records)
		awaitedResult := []string{"\"192.168.0.1\"", "\"192.168.0.2\"", "\"192.168.0.3\"", "\"192.168.0.1\""}
		if !areStringSlicesEqual(wrappedRecords, awaitedResult) {
			t.Errorf("%s records are not wrapped with quotes.", wrappedRecords)
		}
	})
}

func TestIfRecordsListContainsRecord(t *testing.T) {
	recordToCheck := "10.10.0.0"
	t.Run("contains record at first index", func(t *testing.T) {
		recordsList := []string{"192.168.1.1", "10.10.0.0", "0.0.0.0"}
		index, exists := containsRecord(recordsList, recordToCheck)
		if !exists && index != 1 {
			t.Errorf("%s should be existing and being at index 1 of the records list", recordToCheck)
		}
	})

	t.Run("records list does not contain the desired record", func(t *testing.T) {
		recordsList := []string{"192.168.1.1", "10.10.10.10", "0.0.0.0"}
		_, exists := containsRecord(recordsList, recordToCheck)
		if exists {
			t.Errorf("records list should not contain the desired record")
		}
	})
}

func TestRemoveRecordFromValuesList(t *testing.T) {
	t.Run("remove record at first index", func(t *testing.T) {
		recordsList := []string{"192.168.1.1", "10.10.10.10", "0.0.0.0"}
		shortenedList := removeRecordFromValuesList(recordsList, 0)
		awaitedList := []string{"10.10.10.10", "0.0.0.0"}
		if !areStringSlicesEqual(shortenedList, awaitedList) {
			t.Errorf("shortenedList should only contains 2 elements")
		}
	})

	t.Run("remove record at third index", func(t *testing.T) {
		recordsList := []string{"192.168.1.1", "10.10.10.10", "0.0.0.0"}
		shortenedList := removeRecordFromValuesList(recordsList, 2)
		awaitedList := []string{"192.168.1.1", "10.10.10.10"}
		if !areStringSlicesEqual(shortenedList, awaitedList) {
			t.Errorf("shortenedList should only contains 2 elements")
		}
	})
}

func TestGetUpdatedTXTRecordsList(t *testing.T) {
	currentStateRecords := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	managedByHandRecords := []string{"10.10.10.10", "0.0.0.0"}
	apiRecords := append(managedByHandRecords, currentStateRecords...)

	t.Run("remove records in terraform", func(t *testing.T) {
		nextStateRecords := []string{"192.168.1.1"}
		updatedRecordsList := getUpdatedTXTRecordsList(currentStateRecords, apiRecords, nextStateRecords)
		awaitedRecordsList := wrapRecordsWithQuotes(append(managedByHandRecords, nextStateRecords...))
		if !areStringSlicesEqual(updatedRecordsList, awaitedRecordsList) {
			t.Errorf("records list should only contains records managed by hand and new terraform records")
		}
	})

	t.Run("add records in terraform", func(t *testing.T) {
		nextStateRecords := append(currentStateRecords, "192.168.1.4")
		updatedRecordsList := getUpdatedTXTRecordsList(currentStateRecords, apiRecords, nextStateRecords)
		awaitedRecordsList := wrapRecordsWithQuotes(append(managedByHandRecords, nextStateRecords...))
		if !areStringSlicesEqual(updatedRecordsList, awaitedRecordsList) {
			t.Errorf("records list should only contains records managed by hand and new terraform records")
		}
	})
}

func TestKeepRecordsInApiAndTF(t *testing.T) {
	terraformRecords := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	managedByHandRecords := []string{"10.10.10.10", "0.0.0.0"}
	apiRecords := append(terraformRecords, managedByHandRecords...)
	apiRecordsWithQuotes := wrapRecordsWithQuotes(apiRecords)

	t.Run("remove terraform record by hand", func(t *testing.T) {
		awaitedRecords := terraformRecords[1:]
		// api returns records wrapped with quotes
		recordsInBoth := keepRecordsInApiAndTF(terraformRecords, apiRecordsWithQuotes[1:])
		if !areStringSlicesEqual(recordsInBoth, awaitedRecords) {
			t.Errorf("should only contains values that are both in api and terraform")
		}
	})
}
