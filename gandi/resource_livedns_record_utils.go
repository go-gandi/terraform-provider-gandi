package gandi

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func isRecordWrappedWithQuotes(record string) bool {
	return strings.HasPrefix(record, "\"") && strings.HasSuffix(record, "\"")
}

func removeRecordFromValuesList(records []string, index int) []string {
	return append(records[:index], records[index+1:]...)
}

func keepUniqueRecords(recordsList []string) []string {
	keys := make(map[string]bool)
	uniqueRecords := []string{}
	for _, entry := range recordsList {
		if _, exists := keys[entry]; !exists {
			keys[entry] = true
			uniqueRecords = append(uniqueRecords, entry)
		}
	}
	return uniqueRecords
}

func keepRecordsInApiAndTF(tfValues []string, apiValues []string) []string {
	var apiRecordsWithoutQuotes []string
	for _, v := range apiValues {
		if isRecordWrappedWithQuotes(v) {
			apiRecordsWithoutQuotes = append(apiRecordsWithoutQuotes, strings.Trim(v, "\""))
		} else {
			apiRecordsWithoutQuotes = append(apiRecordsWithoutQuotes, v)
		}
	}

	var values []string
	for _, tfv := range tfValues {
		for _, apiv := range apiRecordsWithoutQuotes {
			if tfv == apiv {
				values = append(values, apiv)
			}
		}
	}
	return values
}

func containsRecord(recordsList []string, recordToFind string) (int, bool) {
	for i, rec := range recordsList {
		if rec == recordToFind {
			return i, true
		}
	}
	return 0, false
}

func wrapRecordsWithQuotes(records []string) []string {
	var recordsWithQuotes []string
	for i := range records {
		record := fmt.Sprintf("%v", records[i])
		if isRecordWrappedWithQuotes(record) {
			recordsWithQuotes = append(recordsWithQuotes, record)
		} else {
			recordsWithQuotes = append(recordsWithQuotes, "\""+record+"\"")
		}
	}
	return recordsWithQuotes
}

func areStringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	a_copy := make([]string, len(a))
	b_copy := make([]string, len(b))
	copy(a_copy, a)
	copy(b_copy, b)
	sort.Strings(a_copy)
	sort.Strings(b_copy)
	return reflect.DeepEqual(a_copy, b_copy)
}

func getUpdatedTXTRecordsList(stateRecords, apiRecords, newRecords []string) []string {
	currentRecordsWithQuotes := wrapRecordsWithQuotes(stateRecords)
	apiRecordsWithQuotes := wrapRecordsWithQuotes(apiRecords)
	for _, v := range currentRecordsWithQuotes {
		index, exists := containsRecord(apiRecordsWithQuotes, v)
		if exists {
			apiRecordsWithQuotes = removeRecordFromValuesList(apiRecordsWithQuotes, index)
		}
	}

	records := append(wrapRecordsWithQuotes(newRecords), apiRecordsWithQuotes...)
	return keepUniqueRecords(records)
}
