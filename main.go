package main

import (
	"database/sql"
	"encoding/json"
	"fin_analysis/util"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DefaultInitFile string = ""

func main() {
	var initFile string
	flag.StringVar(&initFile, "init", DefaultInitFile, "Initialize DB with JSON filepath")
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

}

func initializeDB(db *sql.DB, initFile string) error {
	jsonFile, err := os.Open(initFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer jsonFile.Close()
	var jsonMap map[string]interface{}

	jsonBytes, errFile := ioutil.ReadAll(jsonFile)
	if errFile != nil {
		fmt.Println("Failed to read json init file")
		return errFile
	}
	json.Unmarshal(jsonBytes, &jsonMap)

	structured, errStructure := util.RestructureGAAP(jsonMap)
	if errStructure != nil {
		fmt.Printf("Failed to restructure JSON: %v", err)
		return errStructure
	}

	util.CreateTable(db)

	fmt.Println(structured)
	return nil
}
