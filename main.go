package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	fmt.Println("Let's do this")
	_, err := sql.Open("sqlite3", "statements.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
