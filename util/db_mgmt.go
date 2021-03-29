package util

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Create the single table in DB
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

// Insert a single row into financials
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

// Use regex to detect query type, and structure it into a SQL query
// Note: SQL Injection is a concern, but since we run it local, user already has edit access. If this becomes a REST API, add SQL Injection checks
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

	// TODO Update SQL to more accurately query last statement
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
		// We found a matching regex, now let's make the SQL query
		t := queryTemplates[i]
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

		rows, err := db.Query(q)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		defer rows.Close()
		for rows.Next() {
			var err error
			var key, value, row, company string
			switch i {
			case 0:
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

// If we need to initialize DB, clean the local db file prior
func CleanDB(dbFilepath string) {
	os.Remove(dbFilepath)
}
