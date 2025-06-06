package paginate

import (
	"strings"
	"testing"
)

// TestUser represents a test user model
type TestUser struct {
	ID     int    `json:"id" paginate:"id"`
	Name   string `json:"name" paginate:"name"`
	Email  string `json:"email" paginate:"email"`
	Age    int    `json:"age" paginate:"age"`
	Status string `json:"status" paginate:"status"`
	Salary int    `json:"salary" paginate:"salary"`
	DeptID int    `json:"dept_id" paginate:"dept_id"`
}

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() returned nil")
	}
	if builder.params == nil {
		t.Fatal("Builder params is nil")
	}
	if builder.params.Page != 1 {
		t.Errorf("Expected default page to be 1, got %d", builder.params.Page)
	}
	if builder.params.ItemsPerPage != 10 {
		t.Errorf("Expected default items per page to be 10, got %d", builder.params.ItemsPerPage)
	}
}

func TestBuilderBasicUsage(t *testing.T) {
	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Page(2).
		Limit(20).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT * FROM users LIMIT $1 OFFSET $2"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedArgs := []any{20, 20}
	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(args))
	}
}

// TestPaginationParams represents a test struct for pagination parameters
type TestPaginationParams struct {
	Page         int                 `json:"page"`
	Limit        int                 `json:"limit"`
	Search       string              `json:"search"`
	SearchFields []string            `json:"search_fields"`
	LikeOr       map[string][]string `json:"likeor"`
	LikeAnd      map[string][]string `json:"likeand"`
	EqOr         map[string][]any    `json:"eqor"`
	EqAnd        map[string][]any    `json:"eqand"`
	Gte          map[string]any      `json:"gte"`
	Gt           map[string]any      `json:"gt"`
	Lte          map[string]any      `json:"lte"`
	Lt           map[string]any      `json:"lt"`
	Sort         []string            `json:"sort"`
}

func TestFromStruct(t *testing.T) {
	// Test basic struct conversion
	params := TestPaginationParams{
		Page:         2,
		Limit:        25,
		Search:       "john",
		SearchFields: []string{"name", "email"},
		LikeOr: map[string][]string{
			"status": {"active", "pending"},
		},
		Gte: map[string]any{
			"age": 18,
		},
		Sort: []string{"name", "-created_at"},
	}

	builder := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromStruct(params)

	if builder.err != nil {
		t.Fatalf("Unexpected error: %v", builder.err)
	}

	// Verify the parameters were set correctly
	if builder.params.Page != 2 {
		t.Errorf("Expected page to be 2, got %d", builder.params.Page)
	}

	if builder.params.ItemsPerPage != 25 {
		t.Errorf("Expected limit to be 25, got %d", builder.params.ItemsPerPage)
	}

	if builder.params.Search != "john" {
		t.Errorf("Expected search to be 'john', got '%s'", builder.params.Search)
	}

	if len(builder.params.SearchFields) != 2 {
		t.Errorf("Expected 2 search fields, got %d", len(builder.params.SearchFields))
	}

	if len(builder.params.LikeOr["status"]) != 2 {
		t.Errorf("Expected 2 likeor values for status, got %d", len(builder.params.LikeOr["status"]))
	}

	if builder.params.Gte["age"] != 18 {
		t.Errorf("Expected gte age to be 18, got %v", builder.params.Gte["age"])
	}

	// Test with pointer to struct
	builder2 := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromStruct(&params)

	if builder2.err != nil {
		t.Fatalf("Unexpected error with pointer: %v", builder2.err)
	}

	if builder2.params.Page != 2 {
		t.Errorf("Expected page to be 2 with pointer, got %d", builder2.params.Page)
	}
}

func TestFromStructWithNilAndInvalidTypes(t *testing.T) {
	// Test with nil
	builder := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromStruct(nil)

	if builder.err != nil {
		t.Fatalf("Unexpected error with nil: %v", builder.err)
	}

	// Test with invalid type
	builder2 := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromStruct("not a struct")

	if builder2.err == nil {
		t.Fatal("Expected error with invalid type, got nil")
	}

	if !strings.Contains(builder2.err.Error(), "expected struct") {
		t.Errorf("Expected error message to contain 'expected struct', got: %v", builder2.err)
	}
}

func TestFromStructWithJSONTags(t *testing.T) {
	// Test struct with json tags
	type CustomParams struct {
		PageNumber   int    `json:"page"`
		PageSize     int    `json:"limit"`
		SearchTerm   string `json:"search"`
		IgnoredField string `json:"-"`
		EmptyField   string `json:"empty_field"`
	}

	params := CustomParams{
		PageNumber:   3,
		PageSize:     50,
		SearchTerm:   "test",
		IgnoredField: "should be ignored",
		// EmptyField is left empty (zero value)
	}

	builder := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromStruct(params)

	if builder.err != nil {
		t.Fatalf("Unexpected error: %v", builder.err)
	}

	if builder.params.Page != 3 {
		t.Errorf("Expected page to be 3, got %d", builder.params.Page)
	}

	if builder.params.ItemsPerPage != 50 {
		t.Errorf("Expected limit to be 50, got %d", builder.params.ItemsPerPage)
	}

	if builder.params.Search != "test" {
		t.Errorf("Expected search to be 'test', got '%s'", builder.params.Search)
	}
}

func TestBuilderSearch(t *testing.T) {
	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Search("john", "name", "email").
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if !strings.Contains(sql, "name") || !strings.Contains(sql, "email") {
		t.Error("Expected SQL to contain search fields")
	}
	if len(args) == 0 {
		t.Error("Expected search args")
	}
}

func TestBuilderWhereConditions(t *testing.T) {
	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		WhereEquals("status", "active").
		WhereGreaterThan("age", 25).
		WhereLessThanOrEqual("salary", 100000).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if len(args) < 3 {
		t.Errorf("Expected at least 3 args, got %d", len(args))
	}
}

func TestBuilderWhereIn(t *testing.T) {
	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		WhereIn("dept_id", 1, 2, 3).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if len(args) < 3 {
		t.Errorf("Expected at least 3 args, got %d", len(args))
	}
}

func TestBuilderWhereBetween(t *testing.T) {
	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		WhereBetween("age", 18, 65).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if len(args) < 2 {
		t.Errorf("Expected at least 2 args, got %d", len(args))
	}
}

func TestBuilderOrderBy(t *testing.T) {
	sql, _, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		OrderBy("name").
		OrderByDesc("age").
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "ORDER BY") {
		t.Error("Expected SQL to contain ORDER BY clause")
	}
	if !strings.Contains(sql, "name") {
		t.Error("Expected SQL to contain name in ORDER BY")
	}
	if !strings.Contains(sql, "DESC") {
		t.Error("Expected SQL to contain DESC")
	}
}

func TestBuilderJoins(t *testing.T) {
	sql, _, err := NewBuilder().
		Table("users u").
		Model(&TestUser{}).
		LeftJoin("departments d", "u.dept_id = d.id").
		InnerJoin("roles r", "u.role_id = r.id").
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "LEFT JOIN") {
		t.Error("Expected SQL to contain LEFT JOIN")
	}
	if !strings.Contains(sql, "INNER JOIN") {
		t.Error("Expected SQL to contain INNER JOIN")
	}
	if !strings.Contains(sql, "departments d") {
		t.Error("Expected SQL to contain departments table")
	}
	if !strings.Contains(sql, "roles r") {
		t.Error("Expected SQL to contain roles table")
	}
}

func TestBuilderSelect(t *testing.T) {
	sql, _, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Select("id", "name", "email").
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if strings.Contains(sql, "SELECT *") {
		t.Error("Expected SQL not to contain SELECT *")
	}
	if !strings.Contains(sql, "SELECT id, name, email") {
		t.Error("Expected SQL to contain specific columns")
	}
}

func TestBuilderFromJSON(t *testing.T) {
	jsonQuery := `{
		"page": 2,
		"limit": 15,
		"search": "john",
		"search_fields": ["name", "email"],
		"eqor": {
			"status": ["active", "pending"]
		},
		"gte": {
			"age": 18
		},
		"sort": ["name", "-created_at"]
	}`

	sql, args, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		FromJSON(jsonQuery).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if !strings.Contains(sql, "ORDER BY") {
		t.Error("Expected SQL to contain ORDER BY clause")
	}
	if !strings.Contains(sql, "LIMIT") {
		t.Error("Expected SQL to contain LIMIT clause")
	}
	if len(args) == 0 {
		t.Error("Expected args from JSON query")
	}
}

func TestBuilderFromMap(t *testing.T) {
	queryMap := map[string]any{
		"page":          3,
		"limit":         25,
		"search":        "test",
		"search_fields": []string{"name", "description"},
		"eqor": map[string]any{
			"category": []string{"tech", "business"},
		},
		"gt": map[string]any{
			"price": 100,
		},
	}

	sql, args, err := NewBuilder().
		Table("products").
		Model(&TestUser{}).
		FromMap(queryMap).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if !strings.Contains(sql, "LIMIT") {
		t.Error("Expected SQL to contain LIMIT clause")
	}
	if len(args) == 0 {
		t.Error("Expected args from map query")
	}
}

func TestBuilderValidation(t *testing.T) {
	// Test missing table
	_, _, err := NewBuilder().
		Model(&TestUser{}).
		BuildSQL()

	if err == nil {
		t.Error("Expected error for missing table")
	}

	// Test missing model
	_, _, err = NewBuilder().
		Table("users").
		BuildSQL()

	if err == nil {
		t.Error("Expected error for missing model")
	}

	// Test invalid page
	builder := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Page(0)

	if builder.err == nil {
		t.Error("Expected error for invalid page")
	}

	// Test invalid limit
	builder = NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Limit(0)

	if builder.err == nil {
		t.Error("Expected error for invalid limit")
	}
}

func TestBuilderCountSQL(t *testing.T) {
	countSQL, countArgs, err := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		WhereEquals("status", "active").
		WhereGreaterThan("age", 18).
		BuildCountSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(countSQL, "SELECT COUNT") {
		t.Error("Expected count SQL to contain SELECT COUNT")
	}
	if !strings.Contains(countSQL, "WHERE") {
		t.Error("Expected count SQL to contain WHERE clause")
	}
	if strings.Contains(countSQL, "LIMIT") {
		t.Error("Expected count SQL not to contain LIMIT")
	}
	if strings.Contains(countSQL, "OFFSET") {
		t.Error("Expected count SQL not to contain OFFSET")
	}
	if len(countArgs) == 0 {
		t.Error("Expected count args")
	}
}

func TestBuilderComplexQuery(t *testing.T) {
	// Test a complex query with multiple conditions
	sql, args, err := NewBuilder().
		Table("users u").
		Model(&TestUser{}).
		Select("u.id", "u.name", "d.name as dept_name").
		LeftJoin("departments d", "u.dept_id = d.id").
		Search("john", "name", "email").
		WhereEquals("status", "active").
		WhereIn("dept_id", 1, 2, 3).
		WhereGreaterThan("age", 25).
		WhereLessThanOrEqual("salary", 100000).
		Where("u.created_at >= ?", "2023-01-01").
		OrderBy("name").
		OrderByDesc("age").
		Page(2).
		Limit(20).
		BuildSQL()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify all parts are present
	if !strings.Contains(sql, "SELECT u.id, u.name, d.name as dept_name") {
		t.Error("Expected SQL to contain custom SELECT")
	}
	if !strings.Contains(sql, "LEFT JOIN departments d") {
		t.Error("Expected SQL to contain LEFT JOIN")
	}
	if !strings.Contains(sql, "WHERE") {
		t.Error("Expected SQL to contain WHERE clause")
	}
	if !strings.Contains(sql, "ORDER BY") {
		t.Error("Expected SQL to contain ORDER BY clause")
	}
	if !strings.Contains(sql, "LIMIT") {
		t.Error("Expected SQL to contain LIMIT clause")
	}
	if !strings.Contains(sql, "OFFSET") {
		t.Error("Expected SQL to contain OFFSET clause")
	}

	// Should have multiple args
	if len(args) < 3 {
		t.Errorf("Expected at least 3 args, got %d", len(args))
	}
}

func TestBuilderErrorHandling(t *testing.T) {
	// Test that errors are propagated correctly
	builder := NewBuilder().
		Table("users").
		Model(&TestUser{}).
		Page(-1) // This should set an error

	if builder.err == nil {
		t.Error("Expected error to be set")
	}

	// Further operations should not execute
	builder = builder.Limit(10).OrderBy("name")

	_, _, err := builder.BuildSQL()
	if err == nil {
		t.Error("Expected error to be returned from BuildSQL")
	}
}

func TestBuilderChaining(t *testing.T) {
	// Test that all methods return the builder for chaining
	builder := NewBuilder()

	result := builder.
		Table("users").
		Model(&TestUser{}).
		Page(1).
		Limit(10).
		Search("test", "name").
		WhereEquals("status", "active").
		OrderBy("name")

	if result != builder {
		t.Error("Expected method chaining to return the same builder instance")
	}
}
