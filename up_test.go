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
	}{
		{
			FileName: "1_test.up.sql",
			IsValid:  true,
		},
		{
			FileName: "1_test_migration_file.up.sql",
			IsValid:  true,
		},
		{
			FileName: "1_test.down.sql",
			IsValid:  true,
		},
		{
			FileName: "1_test_migration_file.down.sql",
			IsValid:  true,
		},
		{
			FileName: "1_test_up.sql",
			IsValid:  false,
		},
		{
			FileName: "a_test_up.sql",
			IsValid:  false,
		},
	}

	for _, tc := range testCases {
		_, isValid := MigrationFromFile(tc.FileName)
		if isValid != tc.IsValid {
			t.Fatalf("Expected %v but got %v for file name %v", tc.IsValid, isValid, tc.FileName)
		}
	}
}
