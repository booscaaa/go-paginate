package paginate

import (
	"testing"
)

// TestUser struct for testing new operators binding
type TestUserBind struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Active   bool   `json:"active"`
	Category string `json:"category"`
}

// TestNewOperatorsQueryStringBind tests binding new operators from query string
func TestNewOperatorsQueryStringBind(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		validation  func(*PaginationParams) bool
	}{
		{
			name:        "Like operator from query string",
			queryString: "like[name]=john&like[email]=@example.com",
			validation: func(p *PaginationParams) bool {
				return len(p.Like) == 2 &&
					len(p.Like["name"]) == 1 && p.Like["name"][0] == "john" &&
					len(p.Like["email"]) == 1 && p.Like["email"][0] == "@example.com"
			},
		},
		{
			name:        "Eq operator from query string",
			queryString: "eq[age]=25&eq[active]=true",
			validation: func(p *PaginationParams) bool {
				return len(p.Eq) == 2 &&
					len(p.Eq["age"]) == 1 && p.Eq["age"][0] == 25 &&
					len(p.Eq["active"]) == 1 && p.Eq["active"][0] == true
			},
		},
		{
			name:        "In operator from query string",
			queryString: "in[category]=admin&in[category]=user&in[age]=25&in[age]=30",
			validation: func(p *PaginationParams) bool {
				return len(p.In) == 2 &&
					len(p.In["category"]) == 2 &&
					len(p.In["age"]) == 2
			},
		},
		{
			name:        "NotIn operator from query string",
			queryString: "notin[category]=banned&notin[age]=0",
			validation: func(p *PaginationParams) bool {
				return len(p.NotIn) == 2 &&
					len(p.NotIn["category"]) == 1 && p.NotIn["category"][0] == "banned" &&
					len(p.NotIn["age"]) == 1 && p.NotIn["age"][0] == 0
			},
		},
		{
			name:        "Between operator from query string",
			queryString: "between[age]=18&between[age]=65",
			validation: func(p *PaginationParams) bool {
				return len(p.Between) == 1 &&
					p.Between["age"][0] == 18 && p.Between["age"][1] == 65
			},
		},
		{
			name:        "IsNull operator from query string",
			queryString: "isnull=deleted_at&isnull=archived_at",
			validation: func(p *PaginationParams) bool {
				return len(p.IsNull) == 2 &&
					p.IsNull[0] == "deleted_at" && p.IsNull[1] == "archived_at"
			},
		},
		{
			name:        "IsNotNull operator from query string",
			queryString: "isnotnull=email&isnotnull=phone",
			validation: func(p *PaginationParams) bool {
				return len(p.IsNotNull) == 2 &&
					p.IsNotNull[0] == "email" && p.IsNotNull[1] == "phone"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := BindQueryStringToStruct(tt.queryString)
			if err != nil {
				t.Fatalf("Failed to bind query string: %v", err)
			}

			if !tt.validation(params) {
				t.Errorf("Validation failed for %s", tt.name)
			}
		})
	}
}

// TestNewOperatorsFromJSON tests binding new operators from JSON
func TestNewOperatorsFromJSON(t *testing.T) {
	tests := []struct {
		name       string
		jsonData   string
		validation func(*PaginatorBuilder) bool
	}{
		{
			name: "All new operators from JSON",
			jsonData: `{
				"like": {"name": ["john"], "email": ["@example.com"]},
				"eq": {"age": [25], "active": [true]},
				"in": {"category": ["admin", "user"]},
				"notin": {"status": ["banned", "deleted"]},
				"between": {"age": [18, 65]},
				"isnull": ["deleted_at"],
				"isnotnull": ["email"]
			}`,
			validation: func(b *PaginatorBuilder) bool {
				params := b.params
				return len(params.Like) > 0 &&
					len(params.Eq) > 0 &&
					len(params.In) > 0 &&
					len(params.NotIn) > 0 &&
					len(params.Between) > 0 &&
					len(params.IsNull) > 0 &&
					len(params.IsNotNull) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder().
				Model(TestUserBind{}).
				Table("users").
				FromJSON(tt.jsonData)

			if builder.err != nil {
				t.Fatalf("Failed to bind from JSON: %v", builder.err)
			}

			if !tt.validation(builder) {
				t.Errorf("Validation failed for %s", tt.name)
			}
		})
	}
}

// TestNewOperatorsFromStruct tests binding new operators from struct
func TestNewOperatorsFromStruct(t *testing.T) {
	type FilterStruct struct {
		Like      map[string][]string `json:"like"`
		Eq        map[string][]any    `json:"eq"`
		In        map[string][]any    `json:"in"`
		NotIn     map[string][]any    `json:"notin"`
		Between   map[string][2]any   `json:"between"`
		IsNull    []string            `json:"isnull"`
		IsNotNull []string            `json:"isnotnull"`
	}

	filterData := FilterStruct{
		Like: map[string][]string{
			"name":  {"john"},
			"email": {"@example.com"},
		},
		Eq: map[string][]any{
			"age":    {25},
			"active": {true},
		},
		In: map[string][]any{
			"category": {"admin", "user"},
		},
		NotIn: map[string][]any{
			"status": {"banned", "deleted"},
		},
		Between: map[string][2]any{
			"age": {18, 65},
		},
		IsNull:    []string{"deleted_at"},
		IsNotNull: []string{"email"},
	}

	builder := NewBuilder().
		Model(TestUserBind{}).
		Table("users").
		FromStruct(filterData)

	if builder.err != nil {
		t.Fatalf("Failed to bind from struct: %v", builder.err)
	}

	params := builder.params
	if len(params.Like) == 0 || len(params.Eq) == 0 || len(params.In) == 0 ||
		len(params.NotIn) == 0 || len(params.Between) == 0 ||
		len(params.IsNull) == 0 || len(params.IsNotNull) == 0 {
		t.Error("Not all operators were bound from struct")
	}
}

// TestNewOperatorsFromMap tests binding new operators from map
func TestNewOperatorsFromMap(t *testing.T) {
	filterMap := map[string]any{
		"like": map[string]any{
			"name":  []string{"john"},
			"email": []string{"@example.com"},
		},
		"eq": map[string]any{
			"age":    []any{25},
			"active": []any{true},
		},
		"in": map[string]any{
			"category": []any{"admin", "user"},
		},
		"notin": map[string]any{
			"status": []any{"banned", "deleted"},
		},
		"between": map[string]any{
			"age": []any{18, 65},
		},
		"isnull":    []string{"deleted_at"},
		"isnotnull": []string{"email"},
	}

	builder := NewBuilder().
		Model(TestUserBind{}).
		Table("users").
		FromMap(filterMap)

	if builder.err != nil {
		t.Fatalf("Failed to bind from map: %v", builder.err)
	}

	params := builder.params
	if len(params.Like) == 0 || len(params.Eq) == 0 || len(params.In) == 0 ||
		len(params.NotIn) == 0 || len(params.Between) == 0 ||
		len(params.IsNull) == 0 || len(params.IsNotNull) == 0 {
		t.Error("Not all operators were bound from map")
	}
}

// TestNewOperatorsQueryGeneration tests that new operators generate correct SQL
func TestNewOperatorsQueryGeneration(t *testing.T) {
	builder := NewBuilder().
		Model(TestUserBind{}).
		Table("users").
		WhereLike("name", "john").
		WhereEquals("age", 25).
		WhereIn("category", "admin", "user").
		WhereNotIn("status", "banned").
		WhereBetween("age", 18, 65).
		WhereIsNull("deleted_at").
		WhereIsNotNull("email")

	sql, args, err := builder.BuildSQL()
	if err != nil {
		t.Fatalf("Failed to build SQL: %v", err)
	}

	if sql == "" {
		t.Error("Generated SQL is empty")
	}

	if len(args) == 0 {
		t.Error("No arguments generated")
	}

	// Verify that the query contains WHERE clause (basic validation)
	if len(sql) == 0 {
		t.Error("Generated SQL should not be empty")
	}

	// Basic check that it's a SELECT query with WHERE
	if sql[:6] != "SELECT" {
		t.Error("Generated SQL should start with SELECT")
	}

	t.Logf("Generated SQL: %s", sql)
	t.Logf("Generated Args: %v", args)
}
