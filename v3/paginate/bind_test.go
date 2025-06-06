package paginate

import (
	"net/url"
	"testing"
)

func TestBindQueryParamsToStruct(t *testing.T) {
	// Test basic pagination parameters
	queryParams := url.Values{
		"page":            {"2"},
		"limit":           {"25"},
		"search":          {"john"},
		"search_fields":   {"name,email"},
		"sort_columns":    {"name,created_at"},
		"sort_directions": {"ASC,DESC"},
		"vacuum":          {"true"},
		"no_offset":       {"false"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if parameters were set correctly
	if params.Page != 2 {
		t.Errorf("Expected page = 2, got %d", params.Page)
	}

	if params.Limit != 25 {
		t.Errorf("Expected limit = 25, got %d", params.Limit)
	}

	if params.Search != "john" {
		t.Errorf("Expected search = 'john', got '%s'", params.Search)
	}

	if len(params.SearchFields) != 2 {
		t.Errorf("Expected 2 search fields, got %d", len(params.SearchFields))
	}

	if params.SearchFields[0] != "name" || params.SearchFields[1] != "email" {
		t.Errorf("Expected search fields ['name', 'email'], got %v", params.SearchFields)
	}

	if len(params.SortColumns) != 2 {
		t.Errorf("Expected 2 sort columns, got %d", len(params.SortColumns))
	}

	if params.SortColumns[0] != "name" || params.SortColumns[1] != "created_at" {
		t.Errorf("Expected sort columns ['name', 'created_at'], got %v", params.SortColumns)
	}

	if !params.Vacuum {
		t.Error("Expected vacuum = true")
	}

	if params.NoOffset {
		t.Error("Expected no_offset = false")
	}
}

func TestBindQueryParamsWithNestedParameters(t *testing.T) {
	// Test nested parameters like likeor[field], eqor[field], etc.
	queryParams := url.Values{
		"page":              {"1"},
		"limit":             {"10"},
		"likeor[status]": {"active", "pending"},
		"likeand[name]":  {"john"},
		"eqor[age]":    {"25", "30"},
		"eqand[role]":  {"admin"},
		"gte[created_at]":   {"2023-01-01"},
		"gt[score]":         {"80"},
		"lte[updated_at]":   {"2023-12-31"},
		"lt[price]":         {"100.50"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check likeor parameters
	if len(params.LikeOr["status"]) != 2 {
		t.Errorf("Expected 2 likeor values for status, got %d", len(params.LikeOr["status"]))
	}
	if params.LikeOr["status"][0] != "active" || params.LikeOr["status"][1] != "pending" {
		t.Errorf("Expected likeor status ['active', 'pending'], got %v", params.LikeOr["status"])
	}

	// Check likeand parameters
	if len(params.LikeAnd["name"]) != 1 {
		t.Errorf("Expected 1 likeand value for name, got %d", len(params.LikeAnd["name"]))
	}
	if params.LikeAnd["name"][0] != "john" {
		t.Errorf("Expected likeand name 'john', got '%s'", params.LikeAnd["name"][0])
	}

	// Check eqor parameters
	if len(params.EqOr["age"]) != 2 {
		t.Errorf("Expected 2 eqor values for age, got %d", len(params.EqOr["age"]))
	}

	// Check comparison operators
	if params.Gte["created_at"] != "2023-01-01" {
		t.Errorf("Expected gte created_at '2023-01-01', got %v", params.Gte["created_at"])
	}

	if params.Gt["score"] != 80 {
		t.Errorf("Expected gt score 80, got %v", params.Gt["score"])
	}

	if params.Lte["updated_at"] != "2023-12-31" {
		t.Errorf("Expected lte updated_at '2023-12-31', got %v", params.Lte["updated_at"])
	}

	if params.Lt["price"] != 100.5 {
		t.Errorf("Expected lt price 100.5, got %v", params.Lt["price"])
	}
}

func TestBindQueryStringToStruct(t *testing.T) {
	// Test parsing from raw query string
	queryString := "page=3&limit=50&search=test&search_fields=name,email&likeor[status]=active&likeor[status]=pending&gte[age]=18"

	params, err := BindQueryStringToStruct(queryString)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if parameters were set correctly
	if params.Page != 3 {
		t.Errorf("Expected page = 3, got %d", params.Page)
	}

	if params.Limit != 50 {
		t.Errorf("Expected limit = 50, got %d", params.Limit)
	}

	if params.Search != "test" {
		t.Errorf("Expected search = 'test', got '%s'", params.Search)
	}

	if len(params.LikeOr["status"]) != 2 {
		t.Errorf("Expected 2 likeor values for status, got %d", len(params.LikeOr["status"]))
	}

	if params.Gte["age"] != 18 {
		t.Errorf("Expected gte age 18, got %v", params.Gte["age"])
	}
}

func TestBindQueryParamsWithInvalidValues(t *testing.T) {
	// Test with invalid values
	queryParams := url.Values{
		"page":   {"invalid"},
		"limit":  {"abc"},
		"vacuum": {"not_a_bool"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Invalid values should be ignored, defaults should remain
	if params.Page != 1 {
		t.Errorf("Expected page to remain default (1), got %d", params.Page)
	}

	if params.Limit != 10 {
		t.Errorf("Expected limit to remain default (10), got %d", params.Limit)
	}

	if params.Vacuum {
		t.Error("Expected vacuum to remain default (false)")
	}
}

func TestBindQueryParamsWithMultipleValues(t *testing.T) {
	// Test parameters with multiple values
	queryParams := url.Values{
		"search_fields": {"name", "email", "description"},
		"columns":       {"id", "name", "email"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if multiple values are handled correctly
	if len(params.SearchFields) != 3 {
		t.Errorf("Expected 3 search fields, got %d", len(params.SearchFields))
	}

	expectedSearchFields := []string{"name", "email", "description"}
	for i, field := range expectedSearchFields {
		if i >= len(params.SearchFields) || params.SearchFields[i] != field {
			t.Errorf("Expected search field %d to be '%s', got '%s'", i, field, params.SearchFields[i])
		}
	}

	if len(params.Columns) != 3 {
		t.Errorf("Expected 3 columns, got %d", len(params.Columns))
	}
}

func TestBindQueryParamsWithItemsPerPage(t *testing.T) {
	// Test items_per_page parameter
	queryParams := url.Values{
		"items_per_page": {"20"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if items_per_page was set and copied to limit
	if params.ItemsPerPage != 20 {
		t.Errorf("Expected items_per_page = 20, got %d", params.ItemsPerPage)
	}

	if params.Limit != 20 {
		t.Errorf("Expected limit = 20 (copied from items_per_page), got %d", params.Limit)
	}
}

func TestBindQueryParamsCustomStruct(t *testing.T) {
	// Test binding to a custom struct
	type CustomParams struct {
		Page   int    `query:"page"`
		Limit  int    `query:"limit"`
		Search string `query:"search"`
	}

	queryParams := url.Values{
		"page":   {"5"},
		"limit":  {"100"},
		"search": {"custom"},
	}

	customParams := &CustomParams{}
	err := BindQueryParams(queryParams, customParams)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if customParams.Page != 5 {
		t.Errorf("Expected page = 5, got %d", customParams.Page)
	}

	if customParams.Limit != 100 {
		t.Errorf("Expected limit = 100, got %d", customParams.Limit)
	}

	if customParams.Search != "custom" {
		t.Errorf("Expected search = 'custom', got '%s'", customParams.Search)
	}
}

func TestBindQueryParamsInvalidTarget(t *testing.T) {
	// Test with invalid target (not a pointer to struct)
	queryParams := url.Values{"page": {"1"}}

	// Test with non-pointer
	var notPointer PaginationParams
	err := BindQueryParams(queryParams, notPointer)
	if err == nil {
		t.Error("Expected error when passing non-pointer")
	}

	// Test with pointer to non-struct
	var notStruct int
	err = BindQueryParams(queryParams, &notStruct)
	if err == nil {
		t.Error("Expected error when passing pointer to non-struct")
	}
}
