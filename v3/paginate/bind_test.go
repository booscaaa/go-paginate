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
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verificar se os parâmetros foram definidos corretamente
	if params.Page != 2 {
		t.Errorf("Esperado page = 2, obtido %d", params.Page)
	}

	if params.Limit != 25 {
		t.Errorf("Esperado limit = 25, obtido %d", params.Limit)
	}

	if params.Search != "john" {
		t.Errorf("Esperado search = 'john', obtido '%s'", params.Search)
	}

	if len(params.SearchFields) != 2 {
		t.Errorf("Esperado 2 search fields, obtido %d", len(params.SearchFields))
	}

	if params.SearchFields[0] != "name" || params.SearchFields[1] != "email" {
		t.Errorf("Esperado search fields ['name', 'email'], obtido %v", params.SearchFields)
	}

	if len(params.SortColumns) != 2 {
		t.Errorf("Esperado 2 sort columns, obtido %d", len(params.SortColumns))
	}

	if params.SortColumns[0] != "name" || params.SortColumns[1] != "created_at" {
		t.Errorf("Esperado sort columns ['name', 'created_at'], obtido %v", params.SortColumns)
	}

	if !params.Vacuum {
		t.Error("Esperado vacuum = true")
	}

	if params.NoOffset {
		t.Error("Esperado no_offset = false")
	}
}

func TestBindQueryParamsWithNestedParameters(t *testing.T) {
	// Test nested parameters like search_or[field], equals_or[field], etc.
	queryParams := url.Values{
		"page":              {"1"},
		"limit":             {"10"},
		"search_or[status]": {"active", "pending"},
		"search_and[name]":  {"john"},
		"equals_or[age]":    {"25", "30"},
		"equals_and[role]":  {"admin"},
		"gte[created_at]":   {"2023-01-01"},
		"gt[score]":         {"80"},
		"lte[updated_at]":   {"2023-12-31"},
		"lt[price]":         {"100.50"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verificar parâmetros search_or
	if len(params.SearchOr["status"]) != 2 {
		t.Errorf("Esperado 2 valores search_or para status, obtido %d", len(params.SearchOr["status"]))
	}

	if params.SearchOr["status"][0] != "active" || params.SearchOr["status"][1] != "pending" {
		t.Errorf("Esperado search_or status ['active', 'pending'], obtido %v", params.SearchOr["status"])
	}

	// Verificar parâmetros search_and
	if len(params.SearchAnd["name"]) != 1 {
		t.Errorf("Esperado 1 valor search_and para name, obtido %d", len(params.SearchAnd["name"]))
	}

	if params.SearchAnd["name"][0] != "john" {
		t.Errorf("Esperado search_and name 'john', obtido '%s'", params.SearchAnd["name"][0])
	}

	// Verificar parâmetros equals_or
	if len(params.EqualsOr["age"]) != 2 {
		t.Errorf("Esperado 2 valores equals_or para age, obtido %d", len(params.EqualsOr["age"]))
	}

	// Verificar operadores de comparação
	if params.Gte["created_at"] != "2023-01-01" {
		t.Errorf("Esperado gte created_at '2023-01-01', obtido %v", params.Gte["created_at"])
	}

	if params.Gt["score"] != 80 {
		t.Errorf("Esperado gt score 80, obtido %v", params.Gt["score"])
	}

	if params.Lte["updated_at"] != "2023-12-31" {
		t.Errorf("Esperado lte updated_at '2023-12-31', obtido %v", params.Lte["updated_at"])
	}

	if params.Lt["price"] != 100.5 {
		t.Errorf("Esperado lt price 100.5, obtido %v", params.Lt["price"])
	}
}

func TestBindQueryStringToStruct(t *testing.T) {
	// Test parsing from raw query string
	queryString := "page=3&limit=50&search=test&search_fields=name,email&search_or[status]=active&search_or[status]=pending&gte[age]=18"

	params, err := BindQueryStringToStruct(queryString)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verificar se os parâmetros foram definidos corretamente
	if params.Page != 3 {
		t.Errorf("Esperado page = 3, obtido %d", params.Page)
	}

	if params.Limit != 50 {
		t.Errorf("Esperado limit = 50, obtido %d", params.Limit)
	}

	if params.Search != "test" {
		t.Errorf("Esperado search = 'test', obtido '%s'", params.Search)
	}

	if len(params.SearchOr["status"]) != 2 {
		t.Errorf("Esperado 2 valores search_or para status, obtido %d", len(params.SearchOr["status"]))
	}

	if params.Gte["age"] != 18 {
		t.Errorf("Esperado gte age 18, obtido %v", params.Gte["age"])
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
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Valores inválidos devem ser ignorados, padrões devem permanecer
	if params.Page != 1 {
		t.Errorf("Esperado page permanecer padrão (1), obtido %d", params.Page)
	}

	if params.Limit != 10 {
		t.Errorf("Esperado limit permanecer padrão (10), obtido %d", params.Limit)
	}

	if params.Vacuum {
		t.Error("Esperado vacuum permanecer padrão (false)")
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
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verificar se múltiplos valores são tratados corretamente
	if len(params.SearchFields) != 3 {
		t.Errorf("Esperado 3 search fields, obtido %d", len(params.SearchFields))
	}

	expectedSearchFields := []string{"name", "email", "description"}
	for i, field := range expectedSearchFields {
		if i >= len(params.SearchFields) || params.SearchFields[i] != field {
			t.Errorf("Esperado search field %d ser '%s', obtido '%s'", i, field, params.SearchFields[i])
		}
	}

	if len(params.Columns) != 3 {
		t.Errorf("Esperado 3 columns, obtido %d", len(params.Columns))
	}
}

func TestBindQueryParamsWithItemsPerPage(t *testing.T) {
	// Test items_per_page parameter
	queryParams := url.Values{
		"items_per_page": {"20"},
	}

	params, err := BindQueryParamsToStruct(queryParams)
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	// Verificar se items_per_page foi definido e copiado para limit
	if params.ItemsPerPage != 20 {
		t.Errorf("Esperado items_per_page = 20, obtido %d", params.ItemsPerPage)
	}

	if params.Limit != 20 {
		t.Errorf("Esperado limit = 20 (copiado de items_per_page), obtido %d", params.Limit)
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
		t.Fatalf("Erro inesperado: %v", err)
	}

	if customParams.Page != 5 {
		t.Errorf("Esperado page = 5, obtido %d", customParams.Page)
	}

	if customParams.Limit != 100 {
		t.Errorf("Esperado limit = 100, obtido %d", customParams.Limit)
	}

	if customParams.Search != "custom" {
		t.Errorf("Esperado search = 'custom', obtido '%s'", customParams.Search)
	}
}

func TestBindQueryParamsInvalidTarget(t *testing.T) {
	// Test with invalid target (not a pointer to struct)
	queryParams := url.Values{"page": {"1"}}

	// Test with non-pointer
	var notPointer PaginationParams
	err := BindQueryParams(queryParams, notPointer)
	if err == nil {
		t.Error("Esperado erro ao passar não-ponteiro")
	}

	// Test with pointer to non-struct
	var notStruct int
	err = BindQueryParams(queryParams, &notStruct)
	if err == nil {
		t.Error("Esperado erro ao passar ponteiro para não-struct")
	}
}
