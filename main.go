package main

import (
	"database/sql"
	"encoding/json"
	"fin_analysis/util"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	fmt.Println("Let's do this")
	_, err := sql.Open("sqlite3", "statements.db")
	if err != nil {
		fmt.Printf("Failed to open or create db on disk: %v", err)
		os.Exit(1)
	}
	errInit := initializeDB()
	if errInit != nil {
		fmt.Println(errInit)
		os.Exit(2)
	}

}

func initializeDB() error {
	jsonFile, err := os.Open("partial.json")
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
	fmt.Println(structured)
	return nil
}
