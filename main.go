package main

import (
	"database/sql"
	"fin_analysis/util"
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DefaultInitFile string = ""
var DefaultQuery string = ""

func main() {
	var initFile string
	var query string

	flag.StringVar(&initFile, "init", DefaultInitFile, "Initialize DB with JSON filepath")
	flag.StringVar(&query, "query", DefaultQuery, "Supported query")
	flag.Parse()

	db, err := sql.Open("sqlite3", "statements.db")
	if err != nil {
		fmt.Printf("Failed to open or create db on disk: %v", err)
		os.Exit(1)
	}
	if initFile != DefaultInitFile {
		errInit := initializeDB(db, initFile)
		if errInit != nil {
			fmt.Printf("Failed to init with file: %v", initFile)
			os.Exit(2)
		}
	}
	defer db.Close()

	fmt.Println("DB Prepped. Let's do this")
	fmt.Println(util.GenerateSQLFromInput(db, query))

}

func initializeDB(db *sql.DB, initFile string) error {
	lineItems, err := util.ReadLineItemsFromDisk(initFile)
	if err != nil {
		return err
	}
	util.CreateTable(db)
	for _, lineItem := range lineItems {
		util.InsertLineItem(db, lineItem)
	}

	fmt.Println("initialized db")
	return nil
}
