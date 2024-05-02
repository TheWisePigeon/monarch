package main

import "testing"

type testCase struct {
	URL      string
	Expected DB
}

var testCases = []testCase{
	{
		URL:      "postgresql://username:password@localhost:5432/database_name",
		Expected: Postgres,
	},
	{
		URL:      "mysql://username:password@tcp(localhost:3306)/database_name",
		Expected: MySQL,
	},
	{
		URL:      "sqlite://database.db",
		Expected: SQLite,
	},
}

func TestUtils(t *testing.T) {
	for _, tc := range testCases {
		got := detectSQLDB(tc.URL)
		if got != tc.Expected {
			t.Fatalf("Wanted %q got %q", tc.Expected, got)
		}
	}
}
