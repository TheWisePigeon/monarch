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

type DatabaseConn struct {
	db *sql.DB
	tx *sql.Tx
}

func (d *DatabaseConn) Exec(query string) error {
	if d.tx != nil {
		_, err := d.tx.Exec(query)
		return err
	}
	_, err := d.db.Exec(query)
	return err
}

func NewDatabaseConn(db *sql.DB) *DatabaseConn {
	return &DatabaseConn{
		db: db,
		tx: nil,
	}
}

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
	if migrationType != "up" && migrationType != "down" {
		log.Printf("Unrecognized migration type %q", migrationType)
		fmt.Println(helpMessage)
		return
	}
	steps := 0
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

func run(migrationType, source, db string, steps int, safe, verbose bool) int {
	driver := string(getDriver(db))
	if driver == string(SQLite) {
		db = strings.TrimPrefix(db, "sqlite://")
	}
	if verbose {
		log.Println("Connecting to database...")
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
	if verbose {
		log.Println("Connected to database!")
	}
	conn := NewDatabaseConn(dbConn)
	if safe {
		log.Println("'safe' flag set to true: Starting transaction")
		tx, err := conn.db.Begin()
		if err != nil {
			log.Println("Failed to start transaction", err)
			return 1
		}
		conn.tx = tx
		if verbose {
			log.Println("Transaction started")
		}
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
			m, ok := MigrationFromFile(dirEntry.Name(), source, migrationType)
			if !ok {
				continue
			}
			migrations = append(migrations, m)
		}
	}
	for _, m := range migrations {
		switch migrationType {
		case "up":
		case "down":
		}
	}
	return 0
}
