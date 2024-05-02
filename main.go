package main

import (
	"flag"
	"fmt"
	"os"
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

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		runUp(*source, *db, false)
	}
}
