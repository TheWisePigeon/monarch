package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var helpMessage = `monarch: Simple and easy to use migration tool
Usage:
  monarch up: Runs all the up migrations
  monarch down: Runs all the down migrations
  monarch -source <migrations_folder> up/down: Run all the up/down migrations in the specified folder
  monarch up n: Runs all the first n up migrations
  monarch down n: Runs all the down migrations starting from the last down to the nth
`

func main() {
	source := flag.String("source", "migrations/", "The migrations folder")
	db := flag.String("db", "", "The connection URL to the database")
	safe := flag.Bool("safe", true, "The connection URL to the database")
	verbose := flag.Bool("v", false, "While true will cause CLI to print more information during migrations")

	flag.Parse()

	if *db == "" {
		log.Println("Missing required flag: 'db'")
		return
	}

	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	migrationType := flag.Args()[0]
	var steps int
	if len(flag.Args()) >= 2 {
		if migrationSteps, err := strconv.Atoi(flag.Args()[1]); err != nil || migrationSteps < 0 {
			log.Printf("The number of steps if specified must be a positive number! You passed %q", flag.Args()[1])
			return
		} else {
			steps = migrationSteps
		}
	}
	code := run(
		migrationType,
		*source,
		*db,
		steps,
		*safe,
		*verbose,
	)
	os.Exit(code)
}

func run(mType, source, db string, steps int, safe, verbose bool) int {
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
	defer dbConn.Close()
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
	if safe {
		log.Println("up started with 'safe' argument. Starting transaction...")
		tx, err := dbConn.Begin()
		if err != nil {
			log.Println("Error while starting transaction", err)
			return 1
		}
		log.Println("Transaction started!")
		_, err = tx.Exec(query)
		if err != nil {
			log.Println("Error while running migration", err)
			log.Println("Rolling back changes")
			err = tx.Rollback()
			if err != nil {
				log.Println("Error while rolling back changes", err)
			}
			return 1
		}
		log.Println("Migration completed! Commiting changes")
		err = tx.Commit()
		if err != nil {
			log.Println("Error while commiting changes", err)
			return 1
		}
		return 0
	}
	_, err = dbConn.Exec(query)
	if err != nil {
		log.Println("Error while running migration", err)
		return 1
	}
	log.Println("Migration completed!")

	return 0
}
