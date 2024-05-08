package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

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
	if driver == string(SQLite) {
		db = strings.TrimPrefix(db, "sqlite://")
	}
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
	migrations := []*Migration{}
	for _, dirEntry := range entries {
		isFile := !dirEntry.IsDir()
		if isFile {
			m, ok := MigrationFromFile(dirEntry.Name(), source, "up")
			if !ok {
				continue
			}
			migrations = append(migrations, m)
		}
	}
	var query string
	for _, m := range migrations {
		query += m.Query()
	}
	fmt.Println(query)
	return 0
}
