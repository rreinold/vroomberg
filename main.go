package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type LineItem struct {
	Company    string
	Start_date string
	End_date   string
	Key        string 
	Value      string `json:"value"`
}

func main() {

	fmt.Println("Let's do this")
	_, err := sql.Open("sqlite3", "statements.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	initializeDB()

}

func initializeDB() {
	jsonFile, err := os.Open("raw.json")
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal(byteValue, &result)
	fmt.Println(byteValue)
}
