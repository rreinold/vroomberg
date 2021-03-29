package util

import (
	"encoding/json"
	"testing"
)

func TestExtractCompanyName(t *testing.T) {
	inputs := [...]string{"EIX.xbrl"}
	expectations := [...]string{"EIX"}
	for i, input := range inputs {
		expected := expectations[i]
		output, err := extractCompanyName(input)
		if err != nil {
			t.Errorf("Encountered error extracting company name: %v, %v", err, input)
		}
		if output != expected {
			t.Errorf("Did not extract company correctly: %v, %v", input, output)
		}
	}
}

func TestExtractStartAndEndDates(t *testing.T) {
	inputs := [...]string{"('2020-01-01', '2020-12-31')"}
	expectationsStart := [...]string{"2020-01-01"}
	expectationsEnd := [...]string{"2020-12-31"}
	for i, input := range inputs {
		expectedStart := expectationsStart[i]
		expectedEnd := expectationsEnd[i]
		start, end, err := extractStartAndEndDates(input)
		if err != nil {
			t.Errorf("Encountered error extracting stand and end dates: %v, %v", err, input)
		}
		if start != expectedStart || end != expectedEnd {
			t.Errorf("Did not extract company correctly: %v", input)
		}
	}
}

func TestSmokeRestructure(t *testing.T) {
	inputs := [...]string{`{
			"EIX.xbrl": {
				"('2020-01-01', '2020-12-31')": {
				"Revenues": {
					"value": 13578000000,
					"context_ref": "i171001617d444fec804dc14a1178d10d_D20200101-20201231",
					"is_instant_date": false,
					"decimals": "-6"
				},
				"UtilitiesOperatingExpenseMaintenanceAndOperations": {
					"value": 3609000000,
					"context_ref": "i171001617d444fec804dc14a1178d10d_D20200101-20201231",
					"is_instant_date": false,
					"decimals": "-6"
				},
				"LossFromCatastrophes": {
					"value": 1328000000,
					"context_ref": "i171001617d444fec804dc14a1178d10d_D20200101-20201231",
					"is_instant_date": false,
					"decimals": "-6"
				}
			}
		}
	}
`}
	for _, input := range inputs {
		var jsonMap map[string]interface{}
		json.Unmarshal([]byte(input), &jsonMap)
		if _, err := restructureGAAP(jsonMap); err != nil {
			t.Errorf("Encountered error extracting stand and end dates: %v, %v", err, input)
		}
	}
}
