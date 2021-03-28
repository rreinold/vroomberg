package util

import (
	"database/sql"
	"fmt"
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
