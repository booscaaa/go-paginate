// paginate_test.go
package paginate

import (
	"reflect"
	"strings"
	"testing"
)

// User struct used for testing.
type User struct {
	ID    int    `json:"id" paginate:"users.id"`
	Name  string `json:"name" paginate:"users.name"`
	Email string `json:"email" paginate:"users.email"`
	Age   int    `json:"age" paginate:"users.age"`
}

// TestNewPaginator tests the NewPaginator function.
func TestNewPaginator(t *testing.T) {
	// Test case: Missing table should return an error.
	_, err := NewPaginator(
		WithStruct(User{}),
	)
	if err == nil || !strings.Contains(err.Error(), "principal table is required") {
		t.Errorf("Expected error about missing table, got: %v", err)
	}

	// Test case: Missing struct should return an error.
	_, err = NewPaginator(
		WithTable("users"),
	)
	if err == nil || !strings.Contains(err.Error(), "struct is required") {
		t.Errorf("Expected error about missing struct, got: %v", err)
	}

	// Test case: Valid paginator initialization.
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if p.Page != 1 || p.ItemsPerPage != 10 {
		t.Errorf("Unexpected default values: Page=%d, ItemsPerPage=%d", p.Page, p.ItemsPerPage)
	}
}

// TestGenerateSQL tests the GenerateSQL method.
func TestGenerateSQL(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithPage(2),
		WithItemsPerPage(5),
		WithSearch("john"),
		WithSearchFields([]string{"name", "email"}),
		WithSort([]string{"name"}, []string{"false"}),
		WithWhereClause("age > ?", 30),
		WithJoin("INNER JOIN orders ON users.id = orders.user_id"),
		WithColumn("users.id"),
		WithColumn("users.name"),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT users.id, users.name FROM users INNER JOIN orders ON users.id = orders.user_id WHERE (users.name::TEXT ILIKE $1 OR users.email::TEXT ILIKE $2) AND age > $3 ORDER BY users.name ASC LIMIT $4 OFFSET $5"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{"%john%", "%john%", 30, 5, 5}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestGenerateCountQuery tests the GenerateCountQuery method.
func TestGenerateCountQuery(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSearch("doe"),
		WithSearchFields([]string{"name", "email"}),
		WithWhereClause("age > ?", 25),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateCountQuery()
	expectedQuery := "SELECT COUNT(users.id) FROM users WHERE (users.name::TEXT ILIKE $1 OR users.email::TEXT ILIKE $2) AND age > $3"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{"%doe%", "%doe%", 25}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestNoOffset tests the NoOffset option.
func TestNoOffset(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithNoOffset(true),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	if strings.Contains(query, "OFFSET") {
		t.Errorf("Expected query without OFFSET, got: %s", query)
	}
	if len(args) != 1 || args[0] != 10 {
		t.Errorf("Expected args: [10], got: %v", args)
	}
}

// TestVacuumCountQuery tests the Vacuum option in GenerateCountQuery.
func TestVacuumCountQuery(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithVacuum(true),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, _ := p.GenerateCountQuery()
	if !strings.Contains(query, "count_estimate") {
		t.Errorf("Expected count_estimate in query, got: %s", query)
	}
}

// TestWithMapArgs tests the WithMapArgs option.
func TestWithMapArgs(t *testing.T) {
	mapArgs := map[string]any{
		"key1": "value1",
		"key2": 42,
	}
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithMapArgs(mapArgs),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(p.MapArgs, mapArgs) {
		t.Errorf("Expected MapArgs: %v\nGot: %v", mapArgs, p.MapArgs)
	}
}

// TestWithWhereCombining tests the WhereCombining option.
func TestWithWhereCombining(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithWhereCombining("OR"),
		WithWhereClause("age > ?", 20),
		WithWhereClause("age < ?", 30),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	if !strings.Contains(query, "age > $1 OR age < $2") {
		t.Errorf("Expected WHERE clause with OR, got: %s", query)
	}
	expectedArgs := []any{20, 30, 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestGetFieldNameInvalidType tests getFieldName with an invalid type.
func TestGetFieldNameInvalidType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid type, but did not panic")
		}
	}()

	getFieldName("id", "json", "paginate", "not a struct")
}

// TestReplacePlaceholders tests the replacePlaceholders function.
func TestReplacePlaceholders(t *testing.T) {
	query := "SELECT * FROM users WHERE name = ? AND age > ?"
	args := []any{"John", 30}
	expectedQuery := "SELECT * FROM users WHERE name = $1 AND age > $2"

	resultQuery, resultArgs := replacePlaceholders(query, args)
	if resultQuery != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, resultQuery)
	}
	if !reflect.DeepEqual(resultArgs, args) {
		t.Errorf("Expected args: %v\nGot: %v", args, resultArgs)
	}
}

// TestReplacePlaceholdersNoPlaceholders tests replacePlaceholders with no placeholders.
func TestReplacePlaceholdersNoPlaceholders(t *testing.T) {
	query := "SELECT * FROM users"
	args := []any{}
	resultQuery, resultArgs := replacePlaceholders(query, args)
	if resultQuery != query {
		t.Errorf("Expected query unchanged, got: %s", resultQuery)
	}
	if len(resultArgs) != 0 {
		t.Errorf("Expected no args, got: %v", resultArgs)
	}
}

// TestGetFieldName tests the getFieldName function.
func TestGetFieldName(t *testing.T) {
	s := User{}
	fieldName := getFieldName("name", "json", "paginate", s)
	expected := "users.name"
	if fieldName != expected {
		t.Errorf("Expected field name: %s\nGot: %s", expected, fieldName)
	}

	// Test non-existent field.
	fieldName = getFieldName("nonexistent", "json", "paginate", s)
	if fieldName != "" {
		t.Errorf("Expected empty field name, got: %s", fieldName)
	}
}

// TestWithSortInvalidDirections tests WithSort with mismatched directions.
func TestWithSortInvalidDirections(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSort([]string{"name", "age"}, []string{"false"}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, _ := p.GenerateSQL()
	if strings.Contains(query, "ORDER BY") {
		t.Errorf("Expected no ORDER BY clause due to mismatched sort directions, got: %s", query)
	}
}

// TestWithInvalidSearchFields tests invalid search fields.
func TestWithInvalidSearchFields(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSearch("john"),
		WithSearchFields([]string{"nonexistent"}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	if strings.Contains(query, "ILIKE") {
		t.Errorf("Expected no ILIKE clause due to invalid search fields, got: %s", query)
	}
	if len(args) != 2 {
		t.Errorf("Expected no args, got: %v", args)
	}
}

// TestWithEmptySortColumns tests empty sort columns.
func TestWithEmptySortColumns(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSort([]string{}, []string{}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, _ := p.GenerateSQL()
	if strings.Contains(query, "ORDER BY") {
		t.Errorf("Expected no ORDER BY clause, got: %s", query)
	}
}

// TestSchemaUsage tests usage of the Schema option.
func TestSchemaUsage(t *testing.T) {
	p, err := NewPaginator(
		WithSchema("public"),
		WithTable("users"),
		WithStruct(User{}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, _ := p.GenerateSQL()
	if !strings.Contains(query, "FROM public.users") {
		t.Errorf("Expected schema in FROM clause, got: %s", query)
	}
}

// TestComplexWhereClause tests a complex WHERE clause.
func TestComplexWhereClause(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithWhereClause("(age > ? AND age < ?) OR email LIKE ?", 20, 30, "%@example.com"),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQueryPart := "(age > $1 AND age < $2) OR email LIKE $3"
	if !strings.Contains(query, expectedQueryPart) {
		t.Errorf("Expected complex WHERE clause, got: %s", query)
	}
	expectedArgs := []any{20, 30, "%@example.com", 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestMultipleJoinsAndPostInstanceWhereClause tests multiple joins and adding a where clause after paginator instance creation.
func TestMultipleJoinsAndPostInstanceWhereClause(t *testing.T) {
	// Define structs representing the database tables.
	type Order struct {
		ID     int     `json:"id" paginate:"orders.id"`
		UserID int     `json:"user_id" paginate:"orders.user_id"`
		Total  float64 `json:"total" paginate:"orders.total"`
	}

	type Product struct {
		ID      int     `json:"id" paginate:"products.id"`
		OrderID int     `json:"order_id" paginate:"products.order_id"`
		Name    string  `json:"name" paginate:"products.name"`
		Price   float64 `json:"price" paginate:"products.price"`
	}

	// Create the paginator instance with initial options.
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithColumn("users.id"),
		WithColumn("users.name"),
		WithJoin("INNER JOIN orders ON users.id = orders.user_id"),
		WithJoin("INNER JOIN products ON orders.id = products.order_id"),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Add a where clause after the paginator has been created.
	WithWhereClause("products.price > ?", 100.0)(p)
	WithWhereClause("orders.total < ?", 1000.0)(p)

	// Generate the SQL query and arguments.
	query, args := p.GenerateSQL()

	// Expected SQL query.
	expectedQuery := "SELECT users.id, users.name FROM users INNER JOIN orders ON users.id = orders.user_id INNER JOIN products ON orders.id = products.order_id WHERE products.price > $1 AND orders.total < $2 LIMIT $3 OFFSET $4"

	// Check if the generated query matches the expected query.
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	// Expected arguments.
	expectedArgs := []any{100.0, 1000.0, 10, 0}

	// Check if the generated arguments match the expected arguments.
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestNewFilters tests the new filter functionalities.
func TestNewFilters(t *testing.T) {
	// Test SearchOr
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSearchOr(map[string][]string{
			"name": {"vini", "joao"},
		}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE (users.name::TEXT ILIKE $1 OR users.name::TEXT ILIKE $2) LIMIT $3 OFFSET $4"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{"%vini%", "%joao%", 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestSearchAnd tests the SearchAnd filter.
func TestSearchAnd(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSearchAnd(map[string][]string{
			"name": {"john", "doe"},
		}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE (users.name::TEXT ILIKE $1 AND users.name::TEXT ILIKE $2) LIMIT $3 OFFSET $4"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{"%john%", "%doe%", 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestEqualsOr tests the EqualsOr filter.
func TestEqualsOr(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithEqualsOr(map[string][]any{
			"age": {25, 30, 35},
		}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE (users.age = $1 OR users.age = $2 OR users.age = $3) LIMIT $4 OFFSET $5"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{25, 30, 35, 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestEqualsAnd tests the EqualsAnd filter.
func TestEqualsAnd(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithEqualsAnd(map[string][]any{
			"id": {1, 2},
		}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE (users.id = $1 AND users.id = $2) LIMIT $3 OFFSET $4"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{1, 2, 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestComparisonFilters tests the Gte, Gt, Lte, Lt filters.
func TestComparisonFilters(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithGte(map[string]any{"age": 18}),
		WithGt(map[string]any{"id": 0}),
		WithLte(map[string]any{"age": 65}),
		WithLt(map[string]any{"id": 1000}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE users.age >= $1 AND users.id > $2 AND users.age <= $3 AND users.id < $4 LIMIT $5 OFFSET $6"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{18, 0, 65, 1000, 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestCombinedNewFilters tests multiple new filters combined.
func TestCombinedNewFilters(t *testing.T) {
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(User{}),
		WithSearchOr(map[string][]string{
			"name": {"vini", "joao"},
		}),
		WithEqualsOr(map[string][]any{
			"age": {25, 30},
		}),
		WithGte(map[string]any{"id": 1}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	query, args := p.GenerateSQL()
	expectedQuery := "SELECT * FROM users WHERE (users.name::TEXT ILIKE $1 OR users.name::TEXT ILIKE $2) AND (users.age = $3 OR users.age = $4) AND users.id >= $5 LIMIT $6 OFFSET $7"
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	expectedArgs := []any{"%vini%", "%joao%", 25, 30, 1, 10, 0}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}
}

// TestFullComplexPaginator tests a paginator using all properties in a complex scenario.
func TestFullComplexPaginator(t *testing.T) {
	// Define structs representing the database tables.
	type User struct {
		ID       int    `json:"id" paginate:"u.id"`
		Name     string `json:"name" paginate:"u.name"`
		Email    string `json:"email" paginate:"u.email"`
		Age      int    `json:"age" paginate:"u.age"`
		RoleID   int    `json:"role_id" paginate:"u.role_id"`
		IsActive bool   `json:"is_active" paginate:"u.is_active"`
	}

	type Role struct {
		ID   int    `json:"id" paginate:"r.id"`
		Name string `json:"name" paginate:"r.name"`
	}

	type Order struct {
		ID     int     `json:"id" paginate:"o.id"`
		UserID int     `json:"user_id" paginate:"o.user_id"`
		Total  float64 `json:"total" paginate:"o.total"`
		Status string  `json:"status" paginate:"o.status"`
	}

	// Initialize MapArgs with custom parameters.
	mapArgs := map[string]any{
		"min_age":     25,
		"max_age":     35,
		"active_only": true,
		"statuses":    []string{"completed", "shipped"},
	}

	// Create the paginator instance with all properties.
	p, err := NewPaginator(
		WithSchema("public"),
		WithTable("users"),
		WithStruct(User{}),
		WithPage(3),
		WithItemsPerPage(15),
		WithSearch("john doe"),
		WithSearchFields([]string{"name", "email"}),
		WithVacuum(false),
		WithNoOffset(false),
		WithMapArgs(mapArgs),
		WithColumn("u.id"),
		WithColumn("u.name"),
		WithColumn("u.email"),
		WithColumn("r.name AS role_name"),
		WithColumn("SUM(o.total) AS total_spent"),
		WithJoin("INNER JOIN roles r ON u.role_id = r.id"),
		WithJoin("LEFT JOIN orders o ON u.id = o.user_id"),
		WithSort([]string{"name"}, []string{"false"}),
	)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Add complex where clauses after paginator creation.
	WithWhereClause("u.age BETWEEN ? AND ?", mapArgs["min_age"], mapArgs["max_age"])(p)
	WithWhereClause("u.is_active = ?", mapArgs["active_only"])(p)
	WithWhereClause("o.status IN (?)", mapArgs["statuses"])(p)
	WithWhereCombining("AND")(p)

	// Generate the SQL query and arguments.
	query, args := p.GenerateSQL()

	// Expected SQL query.
	expectedQuery := "SELECT u.id, u.name, u.email, r.name AS role_name, SUM(o.total) AS total_spent FROM public.users INNER JOIN roles r ON u.role_id = r.id LEFT JOIN orders o ON u.id = o.user_id WHERE (u.name::TEXT ILIKE $1 OR u.email::TEXT ILIKE $2) AND u.age BETWEEN $3 AND $4 AND u.is_active = $5 AND o.status IN ($6) ORDER BY u.name ASC LIMIT $7 OFFSET $8"

	// Check if the generated query matches the expected query.
	if query != expectedQuery {
		t.Errorf("Expected query:\n%s\nGot:\n%s", expectedQuery, query)
	}

	// Expected arguments.
	expectedArgs := []any{
		"%john doe%",       // Search for name
		"%john doe%",       // Search for email
		mapArgs["min_age"], // u.age BETWEEN ?
		mapArgs["max_age"],
		mapArgs["active_only"], // u.is_active = ?
		mapArgs["statuses"],    // o.status IN (?)
		15,                     // LIMIT
		30,                     // OFFSET (page 3 with 15 items per page)
	}

	// Check if the generated arguments match the expected arguments.
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("Expected args: %v\nGot: %v", expectedArgs, args)
	}

	// Generate the count query.
	countQuery, countArgs := p.GenerateCountQuery()

	// Expected count query.
	expectedCountQuery := "SELECT COUNT(u.id) FROM public.users INNER JOIN roles r ON u.role_id = r.id LEFT JOIN orders o ON u.id = o.user_id WHERE (u.name::TEXT ILIKE $1 OR u.email::TEXT ILIKE $2) AND u.age BETWEEN $3 AND $4 AND u.is_active = $5 AND o.status IN ($6)"

	// Check if the generated count query matches the expected count query.
	if countQuery != expectedCountQuery {
		t.Errorf("Expected count query:\n%s\nGot:\n%s", expectedCountQuery, countQuery)
	}

	// Check if the generated count arguments match the expected arguments.
	if !reflect.DeepEqual(countArgs, expectedArgs[:6]) {
		t.Errorf("Expected count args: %v\nGot: %v", expectedArgs[:6], countArgs)
	}
}
