package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func runUp(source, db string, upTo int, safe bool) int {
	if db == "" {
		fmt.Println("Missing `db` flag")
		return 1
	}
	driver := string(getDriver(db))
	dbConn, err := sql.Open(driver, db)
	if err != nil {
		fmt.Println("Failed to connect to db:", err)
		return 1
	}
	err = dbConn.Ping()
	if err != nil {
		fmt.Println("Failed to connect to db:", err)
		return 1
	}
	entries, err := os.ReadDir(source)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println(source, "folder not found")
			return 1
		}
		fmt.Println("Error while reading your migration files:", err)
		return 1
	}
	for _, dirEntry := range entries {
		isFile := !dirEntry.IsDir()
		if isFile{
		}
	}
	return 0
}
