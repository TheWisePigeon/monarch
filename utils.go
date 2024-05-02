package main

import "strings"

type DB string

const (
	Postgres DB = "postgres"
	MySQL    DB = "mysql"
	SQLite   DB = "sqlite"
)

func detectSQLDB(dbURL string) DB {
	var db DB
	if strings.HasPrefix(dbURL, "postgres") {
		db = Postgres
	}
	if strings.HasPrefix(dbURL, "mysql") {
		db = MySQL
	}
	if strings.HasPrefix(dbURL, "sqlite") {
		db = SQLite
	}
	return db
}
