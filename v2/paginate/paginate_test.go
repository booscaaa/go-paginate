package paginate

import (
	"reflect"
	"strings"
	"testing"
)

func TestPaginQuery(t *testing.T) {
	t.Run("ValidOptions", func(t *testing.T) {
		// Test PaginQuery with valid options
		params, err := PaginQuery(
			WithStruct(struct{}{}),
			WithTable("myTable"),
		)

		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify that the returned params are correctly set
		expectedParams := &paginQueryParams{
			Page:           1,
			ItemsPerPage:   10,
			WhereCombining: "AND",
			Table:          "myTable",
			Struct:         struct{}{},
		}

		if !reflect.DeepEqual(params, expectedParams) {
			t.Errorf("Expected params %+v, got %+v", expectedParams, params)
		}
	})

	t.Run("MissingTable", func(t *testing.T) {
		// Test PaginQuery with missing table
		_, err := PaginQuery(
			WithStruct(struct{}{}),
		)

		if err == nil || err.Error() != "principal table is required" {
			t.Errorf("Expected 'principal table is required' error, got: %v", err)
		}
	})

	t.Run("MissingStruct", func(t *testing.T) {
		// Test PaginQuery with missing struct
		_, err := PaginQuery(
			WithTable("myTable"),
		)

		if err == nil || err.Error() != "struct is required" {
			t.Errorf("Expected 'struct is required' error, got: %v", err)
		}
	})
}

func TestGenerateSQL(t *testing.T) {
	// Create a sample paginQueryParams for testing
	params := &paginQueryParams{
		Page:           1,
		ItemsPerPage:   10,
		WhereCombining: "AND",
		Table:          "myTable",
		Struct:         struct{}{},
	}

	t.Run("BasicTest", func(t *testing.T) {
		// Test GenerateSQL with basic parameters
		sql, _ := GenerateSQL(params)

		// You can add assertions here to check the generated SQL
		// For example, check if it contains "SELECT" and "FROM myTable"
		if !strings.Contains(sql, "SELECT") || !strings.Contains(sql, "FROM myTable") {
			t.Errorf("Generated SQL doesn't match the expected format")
		}
	})

	// Add more test cases to cover other scenarios
}

func TestGenerateCountQuery(t *testing.T) {
	// Create a sample paginQueryParams for testing
	params := &paginQueryParams{
		Page:           1,
		ItemsPerPage:   10,
		WhereCombining: "AND",
		Table:          "myTable",
		Struct:         struct{}{},
	}

	t.Run("BasicTest", func(t *testing.T) {
		// Test GenerateCountQuery with basic parameters
		sql, _ := GenerateCountQuery(params)

		// You can add assertions here to check the generated SQL
		// For example, check if it contains "SELECT COUNT(*)" and "FROM myTable"
		if !strings.Contains(sql, "SELECT COUNT(*)") || !strings.Contains(sql, "FROM myTable") {
			t.Errorf("Generated SQL doesn't match the expected format")
		}
	})

	// Add more test cases to cover other scenarios
}

func TestGetFieldName(t *testing.T) {
	t.Run("ValidTag", func(t *testing.T) {
		// Test GetFieldName with a valid tag
		fieldname := getFieldName("myTag", "json", "paginate", struct {
			MyTag string `json:"myTag" paginate:"myField"`
		}{})

		if fieldname != "myField" {
			t.Errorf("Expected 'myField', got: %s", fieldname)
		}
	})

	t.Run("InvalidTag", func(t *testing.T) {
		// Test GetFieldName with an invalid tag
		fieldname := getFieldName("invalidTag", "json", "paginate", struct {
			MyTag string `json:"myTag" paginate:"myField"`
		}{})

		if fieldname != "" {
			t.Errorf("Expected empty string, got: %s", fieldname)
		}
	})

	// Add more test cases to cover other scenarios
}
