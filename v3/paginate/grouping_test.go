package paginate

import (
	"net/url"
	"strings"
	"testing"
)

func TestGranularOrGrouping(t *testing.T) {
	type User struct {
		ID    int     `json:"id" paginate:"users.id"`
		Name  string  `json:"name" paginate:"users.name"`
		Email string  `json:"email" paginate:"users.email"`
		Age   int     `json:"age" paginate:"users.age"`
		Price float64 `json:"price" paginate:"users.price"`
	}

	t.Run("Mixed EqOr and LikeOr", func(t *testing.T) {
		builder := NewBuilder().
			Table("users").
			Model(User{}).
			EqOr("name", "John").
			LikeOr("email", "teste@example.com")

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name = $") || !strings.Contains(query, "users.email::TEXT ILIKE $") {
			t.Errorf("Expected query to contain grouped conditions, got: %s", query)
		}
		if !strings.Contains(query, " OR ") {
			t.Errorf("Expected OR between conditions, got: %s", query)
		}
	})

	t.Run("Grouped OR between fields", func(t *testing.T) {
		builder := NewBuilder().
			Table("users").
			Model(User{}).
			EqOr("name", "John").
			EqOr("email", "john@example.com")

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name = $") || !strings.Contains(query, "users.email = $") {
			t.Errorf("Expected query to contain grouped conditions, got: %s", query)
		}
		if !strings.Contains(query, " OR ") {
			t.Errorf("Expected OR between conditions, got: %s", query)
		}
	})

	t.Run("Mixed AND and OR", func(t *testing.T) {
		builder := NewBuilder().
			Table("users").
			Model(User{}).
			EqOr("name", "John").
			EqOr("email", "john@example.com").
			EqAnd("age", 18)

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name = $") || !strings.Contains(query, "users.email = $") || !strings.Contains(query, " OR ") {
			t.Errorf("Expected grouped OR with name and email, got: %s", query)
		}
		if !strings.Contains(query, "users.age = $") || !strings.Contains(query, " AND ") {
			t.Errorf("Expected AND for age, got: %s", query)
		}
	})

	t.Run("Query String Binding", func(t *testing.T) {
		values := url.Values{}
		values.Add("eqor[name]", "teste")
		values.Add("eqor[email]", "teste@example.com")
		values.Add("eq[age]", "20")

		paginationParams, _ := BindQueryParamsToStruct(values)

		builder := NewBuilder().
			Table("users").
			Model(User{}).
			fromMap(map[string]any{
				"eqor": paginationParams.EqOr,
				"eq":   paginationParams.Eq,
			})

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name = $") || !strings.Contains(query, "users.email = $") {
			t.Errorf("Expected grouped OR from query string, got: %s", query)
		}
		if !strings.Contains(query, " OR ") {
			t.Errorf("Expected OR between fields, got: %s", query)
		}
	})

	t.Run("IsNull and IsNotNull", func(t *testing.T) {
		builder := NewBuilder().
			Table("users").
			Model(User{}).
			WhereIsNull("name").
			WhereIsNotNull("email")

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name IS NULL") {
			t.Errorf("Expected users.name IS NULL, got: %s", query)
		}
		if !strings.Contains(query, "users.email IS NOT NULL") {
			t.Errorf("Expected users.email IS NOT NULL, got: %s", query)
		}
	})
	t.Run("IsNullOr and IsNotNullOr grouping", func(t *testing.T) {
		builder := NewBuilder().
			Table("users").
			Model(User{}).
			EqOr("status", "active"). // Using map directly for status since it's common in tests
			WhereIsNullOr("name").
			WhereIsNotNullOr("email")

		// status isn't in struct, so let's use name/email for everything
		builder = NewBuilder().
			Table("users").
			Model(User{}).
			EqOr("name", "John").
			WhereIsNullOr("email").
			WhereIsNotNullOr("age")

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "users.name = $") || !strings.Contains(query, "users.email IS NULL") || !strings.Contains(query, "users.age IS NOT NULL") {
			t.Errorf("Expected grouped OR with null checks, got: %s", query)
		}
		if strings.Count(query, " OR ") != 2 {
			t.Errorf("Expected 2 OR operators in the group, got: %s", query)
		}
	})
}

func TestAllOperatorsOrGrouping(t *testing.T) {
	type Product struct {
		ID        int     `json:"id" paginate:"products.id"`
		Name      string  `json:"name" paginate:"products.name"`
		Email     string  `json:"email" paginate:"products.email"`
		Age       int     `json:"age" paginate:"products.age"`
		Price     float64 `json:"price" paginate:"products.price"`
		Stock     int     `json:"stock" paginate:"products.stock"`
		Category  int     `json:"category" paginate:"products.category"`
		Status    string  `json:"status" paginate:"products.status"`
		DeletedAt string  `json:"deleted_at" paginate:"products.deleted_at"`
		CreatedAt string  `json:"created_at" paginate:"products.created_at"`
	}

	t.Run("FromJSON with all operators", func(t *testing.T) {
		jsonData := `{
			"eqor": {"name": "John"},
			"likeor": {"email": "teste"},
			"gteor": {"age": 18},
			"gtor": {"price": 100},
			"lteor": {"stock": 50},
			"ltor": {"id": 1000},
			"inor": {"category": [1, 2, 3]},
			"notinor": {"status": ["deleted"]},
			"isnullor": ["deleted_at"],
			"isnotnullor": ["created_at"],
			"isnull": ["status"],
			"isnotnull": ["name"]
		}`

		builder := NewBuilder().
			Table("products").
			Model(Product{}).
			FromJSON(jsonData)

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		// Verify some of them exist in the SQL
		expectedParts := []string{
			"products.name = $",
			"products.email::TEXT ILIKE $",
			"products.age >= $",
			"products.price > $",
			"products.stock <= $",
			"products.id < $",
			"products.category IN ($",
			"products.status NOT IN ($",
			"products.deleted_at IS NULL",
			"products.created_at IS NOT NULL",
			"products.status IS NULL",
			"products.name IS NOT NULL",
		}

		for _, part := range expectedParts {
			if !strings.Contains(query, part) {
				t.Errorf("Expected query to contain %s, got: %s", part, query)
			}
		}

		if strings.Count(query, " OR ") != 9 {
			t.Errorf("Expected 9 OR operators (10 total conditions), got: %s", query)
		}
	})

	t.Run("FromStruct with OR operators", func(t *testing.T) {
		type ProductFilter struct {
			EqOr  map[string][]any `json:"eqor"`
			GteOr map[string]any   `json:"gteor"`
		}

		filter := ProductFilter{
			EqOr:  map[string][]any{"name": {"Test"}},
			GteOr: map[string]any{"price": 50.0},
		}

		builder := NewBuilder().
			Table("products").
			Model(Product{}).
			FromStruct(filter)

		params, _ := builder.Build()
		query, _ := params.GenerateSQL()

		if !strings.Contains(query, "products.name = $") || !strings.Contains(query, "products.price >= $") || !strings.Contains(query, " OR ") {
			t.Errorf("Expected grouped OR from struct, got: %s", query)
		}
	})
}
