package util

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

func CreateTable(db *sql.DB) error {
	tableSQL := `CREATE TABLE financials (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"company" TEXT,		
		"start_date" TEXT,		
		"end_date" TEXT,
		"key" TEXT,
		"value" DOUBLE PRECISION
	  );`

	statement, err := db.Prepare(tableSQL)
	if err != nil {
		fmt.Printf("Failed to create table: %v", err)
		return err
	}
	statement.Exec()
	return nil
}

func InsertLineItem(db *sql.DB, l LineItem) error {
	insertSQL := `INSERT INTO financials(company, start_date, end_date, key, value) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		fmt.Printf("Failed to prepare insert sql statement: %v", err)
		return err
	}
	_, err = statement.Exec(l.Company, l.Start_date, l.End_date, l.Key, l.Value)
	if err != nil {
		fmt.Printf("Failed to insert row: %v", err)
		return err
	}
	return nil
}

// Type 0: NetIncomeLoss > -400000000
// Type 1: TSLA NetIncomeLoss, TSLA *
// Type 2: TSLA NetIncomeLoss / TSLA OperatingLeasePayments
func GenerateSQLFromInput(db *sql.DB, input string) (string, error) {
	var output string
	queryTypes := [3]string{
		`([a-zA-Z].*?) ([<>]) ([0-9-]*)`,
		`([a-zA-Z].*?) ([a-zA-Z*].*)`,
		`([a-zA-Z].*) ([a-zA-Z].*) ([\/]) ([a-zA-Z].*) ([a-zA-Z]*)`}

	queryTemplates := [3]string{
		`SELECT distinct(company) from financials where key='%v' and value %v %v order by end_date desc;`,
		`SELECT value from financials where company = '%v' and key='%v' order by end_date desc limit 1`,
		`SELECT (value from financials where company = '%v' and key='%v' order by end_date desc limit 1) / (SELECT value from financials where company = '%v' and key='%v' order by end_date desc limit 1)`}

	var rowsArray []string

	for i, queryType := range queryTypes {
		r := regexp.MustCompile(queryType)
		matches := r.FindStringSubmatch(input)
		if matches == nil {
			continue
		}
		t := queryTemplates[i]
		var rows *sql.Rows
		var err error
		var q string
		switch i {
		case 0:
			q = fmt.Sprintf(t, matches[1], matches[2], matches[3])
			break
		case 1:
			q = fmt.Sprintf(t, matches[1], matches[2])
			break
		case 2:
			q = fmt.Sprintf(t, matches[1], matches[2], matches[4], matches[5])
			break
		}

		rows, err = db.Query(q)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			var item string
			if err := rows.Scan(&item); err != nil {
				return "", err
			}
			rowsArray = append(rowsArray, item)
		}
		output = strings.Join(rowsArray, "\n")
	}

	return output, nil
}
