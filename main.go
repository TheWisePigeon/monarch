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
	upCMD := flag.NewFlagSet("up", flag.ExitOnError)
	downCMD := flag.NewFlagSet("down", flag.ExitOnError)
	source := flag.String("source", "migrations/", "The migrations folder")
	db := flag.String("db", "", "The connection URL to the database")
	safe := flag.Bool("safe", true, "The connection URL to the database")

	upTo := upCMD.Int("to", -1, "Specifies the maximum migration version to run up to")
	downTo := downCMD.Int("to", -1, "Specifies the maximum migration version to run down to")

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	cmd := flag.Args()[0]
	switch cmd {
	case "up":
		code := runUp(*source, *db, *upTo, *safe)
		os.Exit(code)
	case "down":
		code := runUp(*source, *db, *downTo, *safe)
		os.Exit(code)
	case "visualize":
		return
	}
}
