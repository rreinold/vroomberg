package util

import (
	"regexp"
)

type LineItem struct {
        Company    string
        Start_date string
        End_date   string
        Key        string
        Value      string `json:"value"`
}

func ExtractCompanyName(raw string) (string,error) {
	regex := regexp.MustCompile(`(.*?).xbrl`)
	matches := regex.FindStringSubmatch(raw)
	return matches[1],nil

}

func RestructureGAAP(map[string]interface{})([]LineItem,error){

	return []LineItem{},nil
}
