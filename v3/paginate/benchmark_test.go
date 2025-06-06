package paginate

import (
	"net/url"
	"testing"
)

// BenchmarkUser represents a user model for benchmarking
type BenchmarkUser struct {
	ID        int    `json:"id" paginate:"users.id"`
	Name      string `json:"name" paginate:"users.name"`
	Email     string `json:"email" paginate:"users.email"`
	Age       int    `json:"age" paginate:"users.age"`
	Status    string `json:"status" paginate:"users.status"`
	Salary    int    `json:"salary" paginate:"users.salary"`
	DeptID    int    `json:"dept_id" paginate:"users.dept_id"`
	CreatedAt string `json:"created_at" paginate:"users.created_at"`
}

// BenchmarkFluentAPI benchmarks the new fluent API
func BenchmarkFluentAPI(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, err := NewBuilder().
			Table("users").
			Model(&BenchmarkUser{}).
			Page(2).
			Limit(20).
			Search("john", "name", "email").
			WhereEquals("status", "active").
			WhereGreaterThan("age", 18).
			OrderBy("name").
			OrderByDesc("created_at").
			BuildSQL()

		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkTraditionalAPI benchmarks the traditional API
func BenchmarkTraditionalAPI(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		p, err := NewPaginator(
			WithTable("users"),
			WithStruct(BenchmarkUser{}),
			WithPage(2),
			WithItemsPerPage(20),
			WithSearch("john"),
			WithSearchFields([]string{"name", "email"}),
			WithSort([]string{"name", "created_at"}, []string{"true", "false"}),
			WithWhereClause("status = ?", "active"),
			WithWhereClause("age > ?", 18),
		)

		if err != nil {
			b.Fatal(err)
		}

		_, _ = p.GenerateSQL()
	}
}

// BenchmarkAutomaticBinding benchmarks the automatic binding feature
func BenchmarkAutomaticBinding(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	// Simulate query parameters
	queryParams := url.Values{
		"page":            {"2"},
		"limit":           {"20"},
		"search":          {"john"},
		"search_fields":   {"name,email"},
		"sort":            {"name,-created_at"},
		"likeor[status]":  {"active", "pending"},
		"gte[age]":        {"18"},
		"columns":         {"id,name,email,age"},
	}

	for i := 0; i < b.N; i++ {
		// Bind query params to struct
		params, err := BindQueryParamsToStruct(queryParams)
		if err != nil {
			b.Fatal(err)
		}

		// Create paginator using traditional API with bound params
		options := []Option{
			WithTable("users"),
			WithStruct(BenchmarkUser{}),
			WithPage(params.Page),
			WithItemsPerPage(params.Limit),
			WithSearch(params.Search),
			WithSearchFields(params.SearchFields),
			WithLikeOr(params.LikeOr),
			WithGte(params.Gte),
		}

		// Add columns if specified
		for _, col := range params.Columns {
			options = append(options, WithColumn(col))
		}

		p, err := NewPaginator(options...)
		if err != nil {
			b.Fatal(err)
		}

		_, _ = p.GenerateSQL()
	}
}

// BenchmarkSQLGeneration benchmarks just the SQL generation part
func BenchmarkSQLGeneration(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	// Pre-create a paginator to avoid initialization overhead
	p, err := NewPaginator(
		WithTable("users"),
		WithStruct(BenchmarkUser{}),
		WithPage(2),
		WithItemsPerPage(20),
		WithSearch("john"),
		WithSearchFields([]string{"name", "email"}),
		WithSort([]string{"name", "created_at"}, []string{"true", "false"}),
		WithWhereClause("status = ?", "active"),
		WithWhereClause("age > ?", 18),
	)

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = p.GenerateSQL()
	}
}