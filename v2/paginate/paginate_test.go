// paginate_test.go
package paginate

import (
	"reflect"
	"testing"
	"time"
)

type S struct {
	DataCriacao time.Time `json:"dataCriacao" paginate:"desktop_log.data_criacao"`
	Modulo      string    `json:"modulo" paginate:"desktop_log.modulo"`
	NomeCliente string    `json:"nomeCliente" paginate:"cliente.nome"`
}

func TestPaginQuery(t *testing.T) {
	// Test case 1: Valid parameters with dynamic schema
	schema := "dynamic_schema" // Use dynamic schema
	params, err := PaginQuery(
		WithStruct(S{}),
		WithSchema(schema), // Use dynamic schema
		WithTable("desktop_log"),
		WithColumn("desktop_log.*"),
		WithPage(2),
		WithItemsPerPage(1),
		WithSort([]string{"dataCriacao", "nomeCliente"}, []string{"true", "false"}),
		WithSearch("oficina"),
		WithSearchFields([]string{"nomeCliente"}),
		WithVacuum(true),
		WithMapArgs(map[string]interface{}{
			"dataCriacao": "2023-09-12",
			"id":          23591765,
			"nomeCliente": "EMPRESA. LTDA",
		}),
		WithWhereClause("teste = ?", "argumento"),
		WithNoOffset(true),
	)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if params == nil {
		t.Error("Expected non-nil paginQueryParams")
	}

	// Test case 2: Missing required parameters (expect error)
	_, err = PaginQuery(
		WithSchema(schema), // Use dynamic schema
		WithTable("desktop_log"),
		WithItemsPerPage(10),
		// Missing WithStruct
	)

	if err == nil {
		t.Error("Expected error for missing struct")
	}
}

func TestGenerateSQL(t *testing.T) {
	// Test case 1: Basic query with dynamic schema
	schema := "dynamic_schema"
	params := &paginQueryParams{
		Page:         1,
		ItemsPerPage: 10,
		Table:        "desktop_log",
		Schema:       schema, // Use dynamic schema
		Struct:       S{},
	}
	sql, args := GenerateSQL(params)

	expectedSQL := "SELECT * FROM dynamic_schema.desktop_log LIMIT $1 OFFSET $2"
	expectedArgs := []interface{}{10, 0}

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, Got: %s", expectedSQL, sql)
	}

	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v, Got: %v", expectedArgs, args)
	}

	// Test case 2: Query with search conditions, custom columns, joins, sorting, and dynamic schema
	params = &paginQueryParams{
		Page:           2,
		ItemsPerPage:   1,
		Search:         "example",
		SearchFields:   []string{"nomeCliente"},
		Columns:        []string{"desktop_log.*", "cliente.nome as nome_cliente"},
		Joins:          []string{"INNER JOIN cliente cliente on cliente.id = desktop_log.id_cliente"},
		SortColumns:    []string{"dataCriacao", "nomeCliente"},
		SortDirections: []string{"true", "false"},
		WhereClauses:   []string{"teste = ?"},
		WhereArgs:      []interface{}{"argumento"},
		Table:          "desktop_log",
		Schema:         schema, // Use dynamic schema
		Struct:         S{},
	}
	sql, args = GenerateSQL(params)

	expectedSQL = "SELECT desktop_log.*, cliente.nome as nome_cliente FROM dynamic_schema.desktop_log " +
		"INNER JOIN cliente cliente on cliente.id = dynamic_schema.desktop_log.id_cliente " +
		"WHERE (cliente.nome::TEXT ILIKE $1) AND teste = $2 " +
		"ORDER BY desktop_log.data_criacao DESC, cliente.nome ASC " +
		"LIMIT $3 OFFSET $4"
	expectedArgs = []interface{}{"%example%", "argumento", 1, 1}

	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, Got: %s", expectedSQL, sql)
	}

	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v, Got: %v", expectedArgs, args)
	}
}

func TestGenerateCountQuery(t *testing.T) {
	// Test case 1: Basic count query with dynamic schema
	schema := "dynamic_schema"
	params := &paginQueryParams{
		Page:         1,
		ItemsPerPage: 10,
		Table:        "desktop_log",
		Schema:       schema, // Use dynamic schema
		Struct:       S{},
	}
	countSQL, countArgs := GenerateCountQuery(params)

	expectedCountSQL := "SELECT COUNT(id) FROM dynamic_schema.desktop_log"
	expectedCountArgs := []interface{}{}

	if countSQL != expectedCountSQL {
		t.Errorf("Expected count SQL: %s, Got: %s", expectedCountSQL, countSQL)
	}

	if !reflect.DeepEqual(countArgs, expectedCountArgs) {
		t.Errorf("Expected count args: %v, Got: %v", expectedCountArgs, countArgs)
	}

	// Test case 2: Count query with search conditions, joins, and dynamic schema
	params = &paginQueryParams{
		Page:         2,
		ItemsPerPage: 1,
		Search:       "example",
		SearchFields: []string{"nomeCliente"},
		Columns:      []string{"desktop_log.*", "cliente.nome as nome_cliente"},
		Joins:        []string{"INNER JOIN cliente cliente on cliente.id = desktop_log.id_cliente"},
		Table:        "desktop_log",
		Schema:       schema, // Use dynamic schema
		Struct:       S{},
	}
	countSQL, countArgs = GenerateCountQuery(params)

	expectedCountSQL = "SELECT COUNT(id) FROM dynamic_schema.desktop_log " +
		"INNER JOIN cliente cliente on cliente.id = dynamic_schema.desktop_log.id_cliente " +
		"WHERE (cliente.nome::TEXT ILIKE $1)"
	expectedCountArgs = []interface{}{"%example%"}

	if countSQL != expectedCountSQL {
		t.Errorf("Expected count SQL: %s, Got: %s", expectedCountSQL, countSQL)
	}

	if !reflect.DeepEqual(countArgs, expectedCountArgs) {
		t.Errorf("Expected count args: %v, Got: %v", expectedCountArgs, countArgs)
	}
}

func TestGetFieldName(t *testing.T) {
	// Test case 1: Valid field tag
	fieldName := getFieldName("dataCriacao", "json", "paginate", S{})
	if fieldName != "desktop_log.data_criacao" {
		t.Errorf("Unexpected field name: %s", fieldName)
	}

	// Test case 2: Missing field tag
	fieldName = getFieldName("invalidField", "json", "paginate", S{})
	if fieldName != "" {
		t.Errorf("Expected empty field name, got: %s", fieldName)
	}
}
