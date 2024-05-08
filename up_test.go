package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestIsMigrationFile(t *testing.T) {
	testCases := []struct {
		FileName string
		IsValid  bool
		Kind     string
	}{
		{
			FileName: "1_test.up.sql",
			IsValid:  true,
			Kind:     "up",
		},
		{
			FileName: "1_test_migration_file.up.sql",
			IsValid:  true,
			Kind:     "up",
		},
		{
			FileName: "1_test.down.sql",
			IsValid:  true,
			Kind:     "down",
		},
		{
			FileName: "1_test_migration_file.down.sql",
			IsValid:  true,
			Kind:     "down",
		},
		{
			FileName: "1_test_up.sql",
			IsValid:  false,
			Kind:     "up",
		},
		{
			FileName: "a_test_up.sql",
			IsValid:  false,
			Kind:     "up",
		},
	}

	for _, tc := range testCases {
		_, isValid := MigrationFromFile(tc.FileName, "testdata", tc.Kind)
		if isValid != tc.IsValid {
			t.Fatalf("Expected %v but got %v for file name %v", tc.IsValid, isValid, tc.FileName)
		}
	}
}

func TestUpMigration(t *testing.T) {
	source := "migrations"
  db := "sqlite://testdata/test.db"
	upTo := -1
	safe := false
	code := runUp(source, db, upTo, safe)
	if code != 0 {
		t.Fatalf("Expected 0 but got %v", code)
	}
}
