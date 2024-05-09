package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Driver string

const (
	Postgres Driver = "postgres"
	MySQL    Driver = "mysql"
	SQLite   Driver = "sqlite3"
)

func getDriver(dbURL string) Driver {
	var db Driver
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

type Migration struct {
	//version
	version string
	//title of the migration for printing
	title string
	//up or down migration
	kind string
	//content of the migration file
	sql []byte
}

func (m *Migration) Query() string {
	return string(m.sql)
}

func MigrationFromFile(fileName, source, kind string) (*Migration, bool) {
	m := Migration{}
	parts := strings.SplitN(fileName, "_", 2)
	if len(parts) != 2 {
		return nil, false
	}
	if _, err := strconv.Atoi(parts[0]); err != nil {
		return nil, false
	}
	m.version = parts[0]
	parts = strings.SplitN(parts[1], ".", 2)
	if len(parts) != 2 {
		return nil, false
	}
	m.title = parts[0]
	parts = strings.Split(parts[1], ".")
	if len(parts) != 2 {
		return nil, false
	}
	if parts[1] != "sql" {
		return nil, false
	}
	if parts[0] != kind {
		return nil, false
	}
	m.kind = parts[0]
	file, err := os.Open(fmt.Sprintf("%s/%s", source, fileName))
	if err != nil {
		return nil, false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return nil, false
	}
	m.sql = make([]byte, stat.Size())
	_, err = file.Read(m.sql)
	if err != nil {
		return nil, false
	}
	return &m, true
}

func GetDBVersion() {

}

func SetDBVersion() {

}
