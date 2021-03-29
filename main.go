package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"vroomberg/util"

	_ "github.com/mattn/go-sqlite3"
)

var DefaultInitFile string = ""
var DefaultQuery string = ""
var DBFilepath string = "./statements.db"

func main() {
	var initFile string
	var query string

	flag.StringVar(&initFile, "init", DefaultInitFile, "Initialize DB (./statements.db) with JSON filepath")
	flag.StringVar(&query, "query", DefaultQuery, "Supported query")
	flag.Parse()

	db, err := sql.Open("sqlite3", DBFilepath)
	if err != nil {
		fmt.Printf("Failed to open or create db on disk: %v", err)
		os.Exit(1)
	}
	if initFile != DefaultInitFile {
		util.CleanDB(DBFilepath)
		errInit := initializeDB(db, initFile)
		if errInit != nil {
			fmt.Printf("Failed to init with file: %v", initFile)
			os.Exit(2)
		}
	}
	defer db.Close()

	// If no query provided, exit
	if query == DefaultQuery {
		os.Exit(0)
	}

	output, err := util.GenerateSQLFromInput(db, query)
	if err != nil {
		fmt.Printf("Failed to query DB: %v", err)
		os.Exit(3)
	}
	fmt.Println(output)

}

// Initialize DB by reading JSON, running INSERTS
func initializeDB(db *sql.DB, initFile string) error {
	fmt.Println("Structuring data...")
	lineItems, err := util.ReadLineItemsFromDisk(initFile)
	if err != nil {
		return err
	}
	// TODO Roadmap Item #2: Bulk INSERT for initializing DB
	fmt.Printf("Initializing db...")
	util.CreateTable(db)
	for i, lineItem := range lineItems {
		util.InsertLineItem(db, lineItem)
		if i%10000 == 0 {
			fmt.Printf(".")
		}
	}

	fmt.Println("initialized db")
	return nil
}
