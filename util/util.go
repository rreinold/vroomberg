package util

import (
	"errors"
	"fmt"
	"math"
	"regexp"
)

type LineItem struct {
	Company    string
	Start_date string
	End_date   string
	Key        string
	Value      string `json:"value"`
}

func extractCompanyName(raw string) (string, error) {
	regex := regexp.MustCompile(`(.*?).xbrl`)
	matches := regex.FindStringSubmatch(raw)
	if len(matches) != 2 {
		return "", errors.New("Failed to find company name in expected format")
	}
	return matches[1], nil

}

func extractStartAndEndDates(raw string) (string, string, error) {
	regex := regexp.MustCompile(`\('([0-9-]*)', '([0-9-]*)'\)`)
	matches := regex.FindStringSubmatch(raw)
	if len(matches) != 3 {
		return "", "", errors.New("Failed to find start and enddates in expected format")
	}
	return matches[1], matches[2], nil
}

func calculateValue(value float64, decimalExponent int) float64 {
	return value * math.Pow10(decimalExponent)
}

func RestructureGAAP(root map[string]interface{}) ([]LineItem, error) {
	fmt.Println("Restructuring")
	var output []LineItem
	for company, companyChildren := range root {
		fmt.Println(company)
		for dateRange, dateRangeChildren := range companyChildren.(map[string]interface{}) {
			fmt.Println(dateRange)
			for metric, metricChildren := range dateRangeChildren.(map[string]interface{}) {
				companyName, err := extractCompanyName(company)
				if err != nil {
					return []LineItem{}, err
				}
				startDate, endDate, err2 := extractStartAndEndDates(dateRange)
				if err2 != nil {
					return []LineItem{}, err2
				}
				values := metricChildren.(map[string]interface{})
				value := values["value"].(float64)
				if values["decimal"] != nil {
					rawDecimalExponent := values["decimal"].(int)
					value = calculateValue(value, rawDecimalExponent)
				}
				fmt.Println(values)
				lineItem := LineItem{
					Company:    companyName,
					Start_date: startDate,
					End_date:   endDate,
					Key:        metric,
					Value:      fmt.Sprintf("%d", value)}
				fmt.Printf("%v", lineItem)
				output = append(output, lineItem)
			}
		}
	}
	return output, nil
}
