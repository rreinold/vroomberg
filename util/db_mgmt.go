package util

import (
	"database/sql"
	"fmt"
	"os"
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
// Type 1: TSLA NetIncomeLoss
// Type 2: TSLA NetIncomeLoss / TSLA OperatingLeasePayments
// Type 3: TSLA *
func GenerateSQLFromInput(db *sql.DB, input string) (string, error) {
	var output string
	queryTypes := [4]string{
		`([a-zA-Z].*?) ([<>]) ([0-9-]*)`,
		`([a-zA-Z].*?) ([a-zA-Z].*)`,
		`([a-zA-Z].*) ([a-zA-Z].*) [\/] ([a-zA-Z].*) ([a-zA-Z]*)`,
		`([a-zA-Z].*?) \*`}

	queryTemplates := [4]string{
		`SELECT distinct(company) from financials where key='%v' and value %v %v order by end_date desc;`,
		`SELECT value from financials where company = '%v' and key='%v' order by end_date desc limit 1;`,
		`SELECT A.value / B.value AS value FROM   (SELECT value from financials where company = '%v' and key='%v' order by end_date desc limit 1) A,(SELECT value from financials where company = '%v' and key='%v' order by end_date desc limit 1) B;`,
		`SELECT distinct(key), value from financials where company = '%v' order by end_date desc;`}

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
			q = fmt.Sprintf(t, matches[1], matches[2], matches[3], matches[4])
			break
		case 3:
			q = fmt.Sprintf(t, matches[1])
			break
		}

		rows, err = db.Query(q)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			var value string
			var key string
			var err error
			var row string
			switch i {
			case 0:
				var company string
				err = rows.Scan(&company)
				row = company
				break
			case 3:
				err = rows.Scan(&key, &value)
				row = key + ", " + value
				break
			default:
				err = rows.Scan(&value)
				row = value
				break
			}
			if err != nil {
				return "", err
			}

			rowsArray = append(rowsArray, row)
		}
		output = strings.Join(rowsArray, "\n")
	}

	return output, nil
}

func CleanDB(dbFilepath string) {
	os.Remove(dbFilepath)
}
