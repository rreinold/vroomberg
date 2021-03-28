package util

import (
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
