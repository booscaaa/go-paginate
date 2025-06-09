package paginate

import (
	"testing"
)

// TestNewOperators tests all the newly implemented operators
func TestNewOperators(t *testing.T) {
	tests := []struct {
		name    string
		options []Option
	}{
		{
			name: "Like operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithLike(map[string][]string{
					"name": {"john"},
				}),
			},
		},
		{
			name: "Eq operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithEq(map[string][]any{
					"status": {"active"},
				}),
			},
		},
		{
			name: "In operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithIn(map[string][]any{
					"age": {25, 30, 35},
				}),
			},
		},
		{
			name: "NotIn operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithNotIn(map[string][]any{
					"status": {"deleted", "banned"},
				}),
			},
		},
		{
			name: "Between operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithBetween(map[string][2]any{
					"age": {18, 65},
				}),
			},
		},
		{
			name: "IsNull operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithIsNull([]string{"role_id"}),
			},
		},
		{
			name: "IsNotNull operator",
			options: []Option{
				WithStruct(TestUser{}),
				WithTable("users"),
				WithIsNotNull([]string{"role_id"}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := NewPaginator(tt.options...)
			if err != nil {
				t.Fatalf("NewPaginator failed: %v", err)
			}

			query, _ := params.GenerateSQL()
			if query == "" {
				t.Fatal("Generated query is empty")
			}

			// Just verify that query generation doesn't fail
			t.Logf("Generated query for %s: %s", tt.name, query)
		})
	}
}

// TestBuilderNewOperators tests the builder pattern with new operators
func TestBuilderNewOperators(t *testing.T) {
	builder := NewBuilder().Model(TestUser{}).Table("users")

	// Test all new builder methods
	builder.
		WhereLike("name", "john").
		WhereIn("age", 25, 30, 35).
		WhereNotIn("status", "deleted", "banned").
		WhereBetween("age", 18, 65).
		WhereIsNull("role_id").
		WhereIsNotNull("email")

	params, err := builder.Build()
	if err != nil {
		t.Fatalf("Builder.Build() failed: %v", err)
	}

	query, _ := params.GenerateSQL()
	if query == "" {
		t.Fatal("Generated query is empty")
	}

	// Just verify that query generation doesn't fail and contains WHERE clause
	if !contains(query, "WHERE") {
		t.Error("Expected query to contain WHERE clause")
	}

	t.Logf("Generated builder query: %s", query)
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
				s[len(s)-len(substr):] == substr || 
				containsInMiddle(s, substr))))
}

func containsInMiddle(s, substr string) bool {
	for i := 1; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}