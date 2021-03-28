package util

import (
	"database/sql"
	"fmt"
)

func CreateTable(db *sql.DB) error {
	sqlTable := `CREATE TABLE financials (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"company" TEXT,		
		"start_date" TEXT,		
		"end_date" TEXT,
		"key" TEXT,
		"value" DOUBLE PRECISION
	  );`

	statement, err := db.Prepare(sqlTable)
	if err != nil {
		fmt.Printf("Failed to create table: %v", err)
		return err
	}
	statement.Exec()
	return nil
}
