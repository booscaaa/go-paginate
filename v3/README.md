<p align="center">
  <h1 align="center">Go Paginate v3 - Go package to generate query pagination</h1>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/booscaaa/go-paginate/v3"><img alt="Reference" src="https://img.shields.io/badge/go-reference-purple?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/booscaaa/go-paginate.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge"></a>
    <img alt="GitHub Workflow Status (with event)" src="https://img.shields.io/github/actions/workflow/status/booscaaa/go-paginate/test.yaml?style=for-the-badge">
    <a href="https://codecov.io/gh/booscaaa/go-paginate"><img alt="Coverage" src="https://img.shields.io/codecov/c/github/booscaaa/go-paginate/master.svg?style=for-the-badge"></a>
  </p>
</p>

<br>

## Why?

This project is part of my personal portfolio, so, I'll be happy if you could provide me any feedback about the project, code, structure or anything that you can report that could make me a better developer!

Email-me: boscardinvinicius@gmail.com

Connect with me at [LinkedIn](https://www.linkedin.com/in/booscaaa/).

<br>

# Go Paginate v3

A powerful and flexible Go package for generating paginated SQL queries with advanced filters. v3 introduces a modern **fluent API**, **automatic query parameter binding**, and maintains compatibility with the traditional API.

## 🚀 What's New in v3

- ✨ **Fluent API (Builder Pattern)** - Modern and intuitive interface
- 🔗 **Automatic Binding** - Converts HTTP query parameters automatically
- 🔍 **Advanced Filters** - SearchOr, SearchAnd, EqualsOr, EqualsAnd, Gte, Gt, Lte, Lt
- 🔄 **Compatibility** - Maintains support for traditional API
- 📊 **Complex Joins** - Full support for JOINs
- 🎯 **Type Safety** - Compile-time type validation

## 📦 Installation

```bash
go get github.com/booscaaa/go-paginate/v3
```

## 🎯 Usage Methods

### 1. 🌟 Fluent API (Recommended)

The new fluent API offers a modern and intuitive interface:

```go
package main

import (
    "fmt"
    "log"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

type User struct {
    ID        int    `json:"id" paginate:"users.id"`
    Name      string `json:"name" paginate:"users.name"`
    Email     string `json:"email" paginate:"users.email"`
    Age       int    `json:"age" paginate:"users.age"`
    Status    string `json:"status" paginate:"users.status"`
    CreatedAt string `json:"created_at" paginate:"users.created_at"`
}

func main() {
    // Fluent API - Basic usage
    sql, args, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        Page(2).
        Limit(20).
        Search("john", "name", "email").
        OrderBy("name").
        OrderByDesc("created_at").
        BuildSQL()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("SQL: %s\n", sql)
    fmt.Printf("Args: %v\n", args)
}
```

#### Advanced Filters with Fluent API

```go
// Complex filters
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    WhereEquals("status", "active").
    WhereIn("dept_id", 1, 2, 3).
    WhereGreaterThan("age", 25).
    WhereLessThanOrEqual("salary", 100000).
    WhereBetween("created_at", "2023-01-01", "2023-12-31").
    SearchOr("name", "John", "Jane").
    SearchAnd("email", "@company.com").
    BuildSQL()
```

#### Joins with Fluent API

```go
// Complex joins
sql, args, err := paginate.NewBuilder().
    Table("users u").
    Model(&User{}).
    InnerJoin("departments d", "u.dept_id = d.id").
    LeftJoin("profiles p", "u.id = p.user_id").
    Select("u.*", "d.name as dept_name", "p.avatar").
    WhereEquals("u.status", "active").
    OrderBy("u.name").
    BuildSQL()
```

### 2. 🔗 Automatic Query Parameter Binding

Automatically convert HTTP query parameters to pagination structs:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
    // Extract and convert query parameters automatically
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    // Use the converted parameters
    fmt.Printf("Page: %d\n", params.Page)
    fmt.Printf("Limit: %d\n", params.Limit)
    fmt.Printf("Search: %s\n", params.Search)
    fmt.Printf("SearchOr: %v\n", params.SearchOr)
    fmt.Printf("EqualsOr: %v\n", params.EqualsOr)
}

// Example URL:
// /users?page=2&limit=25&search=john&search_or[status]=active&search_or[status]=pending&equals_or[age]=25&equals_or[age]=30
```

#### Binding from Query String

```go
// From a query string
queryString := "page=2&limit=25&search=john&search_fields=name,email&vacuum=true"
params, err := paginate.BindQueryStringToStruct(queryString)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Page: %d\n", params.Page)           // 2
fmt.Printf("Limit: %d\n", params.Limit)         // 25
fmt.Printf("Search: %s\n", params.Search)       // "john"
fmt.Printf("Fields: %v\n", params.SearchFields) // ["name", "email"]
```

#### Supported Complex Parameters

| Parameter           | Type                | Example                                              |
| ------------------- | ------------------- | ---------------------------------------------------- |
| `search_or[field]`  | map[string][]string | `search_or[status]=active&search_or[status]=pending` |
| `search_and[field]` | map[string][]string | `search_and[name]=john`                              |
| `equals_or[field]`  | map[string][]any    | `equals_or[age]=25&equals_or[age]=30`                |
| `equals_and[field]` | map[string][]any    | `equals_and[role]=admin`                             |
| `gte[field]`        | map[string]any      | `gte[age]=18`                                        |
| `gt[field]`         | map[string]any      | `gt[score]=80`                                       |
| `lte[field]`        | map[string]any      | `lte[price]=100.50`                                  |
| `lt[field]`         | map[string]any      | `lt[date]=2023-12-31`                                |

### 3. 🔧 Traditional API (Compatibility)

Maintains full compatibility with previous versions:

```go
// API tradicional com opções
p, err := paginate.NewPaginator(
    paginate.WithTable("users"),
    paginate.WithStruct(User{}),
    paginate.WithSearchOr(map[string][]string{
        "name": {"vini", "joao"},
    }),
    paginate.WithEqualsOr(map[string][]any{
        "age": {25, 30, 35},
    }),
    paginate.WithGte(map[string]any{"id": 1}),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(20),
)
if err != nil {
    log.Fatal(err)
}

query, args := p.GenerateSQL()
countQuery, countArgs := p.GenerateCountQuery()
```

## 🔍 Available Advanced Filters

### Search Filters

- **SearchOr**: Search for multiple values using OR
- **SearchAnd**: Search for multiple values using AND

### Equality Filters

- **EqualsOr**: Equality with multiple values using OR
- **EqualsAnd**: Equality with multiple values using AND

### Comparison Filters

- **Gte**: Greater than or equal (>=)
- **Gt**: Greater than (>)
- **Lte**: Less than or equal (<=)
- **Lt**: Less than (<)

## 📋 Complete Examples

### Example 1: Combined Complex Filters

```go
// Combining multiple filters
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Page(1).
    Limit(50).
    // Search OR: name contains "john" OR "jane"
    SearchOr("name", "john", "jane").
    // Search AND: email contains "@company.com"
    SearchAnd("email", "@company.com").
    // Equality OR: status is "active" OR "pending"
    WhereIn("status", "active", "pending").
    // Comparisons: age between 25 and 65
    WhereGreaterThanOrEqual("age", 25).
    WhereLessThanOrEqual("age", 65).
    // Ordering
    OrderBy("name").
    OrderByDesc("created_at").
    BuildSQL()
```

### Example 2: Joins with Filters

```go
type UserWithDepartment struct {
    ID       int    `json:"id" paginate:"u.id"`
    Name     string `json:"name" paginate:"u.name"`
    Email    string `json:"email" paginate:"u.email"`
    DeptName string `json:"dept_name" paginate:"d.name"`
    Salary   int    `json:"salary" paginate:"u.salary"`
}

sql, args, err := paginate.NewBuilder().
    Table("users u").
    Model(&UserWithDepartment{}).
    InnerJoin("departments d", "u.dept_id = d.id").
    LeftJoin("salaries s", "u.id = s.user_id").
    Select("u.*", "d.name as dept_name", "s.amount as salary").
    WhereEquals("u.status", "active").
    WhereIn("d.type", "engineering", "sales", "marketing").
    WhereGreaterThan("s.amount", 50000).
    Search("john", "u.name", "u.email").
    OrderBy("d.name", "u.name").
    BuildSQL()
```

### Example 3: Using with HTTP Handler

```go
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Automatic binding of query parameters
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    // 2. Build query using Builder
    builder := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        Page(params.Page).
        Limit(params.Limit)

    // 3. Apply filters if provided
    if params.Search != "" {
        builder = builder.Search(params.Search, "name", "email")
    }

    // 4. Apply advanced filters
    for field, values := range params.SearchOr {
        builder = builder.SearchOr(field, values...)
    }

    for field, values := range params.EqualsOr {
        builder = builder.WhereIn(field, values...)
    }

    // 5. Generate SQL
    sql, args, err := builder.BuildSQL()
    if err != nil {
        http.Error(w, "Query build error", http.StatusInternalServerError)
        return
    }

    // 6. Execute query on database
    rows, err := db.Query(sql, args...)
    // ... process results
}
```

### Example 4: Custom Struct for Binding

```go
type CustomSearchParams struct {
    Page      int      `query:"page"`
    Limit     int      `query:"limit"`
    Search    string   `query:"search"`
    Category  string   `query:"category"`
    Tags      []string `query:"tags"`
    MinPrice  float64  `query:"min_price"`
    MaxPrice  float64  `query:"max_price"`
    Active    bool     `query:"active"`
}

func productSearchHandler(w http.ResponseWriter, r *http.Request) {
    customParams := &CustomSearchParams{
        Page:  1,
        Limit: 20,
        Active: true, // default value
    }

    err := paginate.BindQueryParams(r.URL.Query(), customParams)
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    // Use custom parameters
    builder := paginate.NewBuilder().
        Table("products").
        Model(&Product{}).
        Page(customParams.Page).
        Limit(customParams.Limit)

    if customParams.Search != "" {
        builder = builder.Search(customParams.Search, "name", "description")
    }

    if customParams.Category != "" {
        builder = builder.WhereEquals("category", customParams.Category)
    }

    if len(customParams.Tags) > 0 {
        builder = builder.WhereIn("tag", customParams.Tags...)
    }

    if customParams.MinPrice > 0 {
        builder = builder.WhereGreaterThanOrEqual("price", customParams.MinPrice)
    }

    if customParams.MaxPrice > 0 {
        builder = builder.WhereLessThanOrEqual("price", customParams.MaxPrice)
    }

    builder = builder.WhereEquals("active", customParams.Active)

    sql, args, err := builder.BuildSQL()
    // ... execute query
}
```

## 🎯 API Comparison

| Feature              | Fluent API | Traditional API | Automatic Binding |
| -------------------- | ---------- | --------------- | ----------------- |
| **Ease of use**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐          | ⭐⭐⭐⭐⭐        |
| **Type Safety**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐        | ⭐⭐⭐            |
| **Flexibility**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐      | ⭐⭐⭐            |
| **Performance**      | ⭐⭐⭐⭐   | ⭐⭐⭐⭐        | ⭐⭐⭐⭐          |
| **HTTP Integration** | ⭐⭐⭐     | ⭐⭐            | ⭐⭐⭐⭐⭐        |

### When to use each API:

- **🌟 Fluent API**: For new projects, complex queries, better readability
- **🔧 Traditional API**: For migrating existing code, compatibility
- **🔗 Automatic Binding**: For REST APIs, integration with web frameworks

## 🚀 Running Examples

To see the features in action, run the provided examples:

```bash
# Fluent API (Builder)
go run example_builder.go

# Query Parameter Binding
go run example_bind.go

# Traditional API
go run example_usage.go
```

## 🧪 Running Tests

```bash
# All tests
go test -v ./paginate

# Specific tests
go test -v ./paginate -run TestBuilder
go test -v ./paginate -run TestBind
go test -v ./paginate -run TestPaginate

# With coverage
go test -v -cover ./paginate
```

## 📚 Additional Documentation

- [📖 Complete Bind Documentation](BIND_README.md) - Detailed guide on query parameter binding
- [🔗 Go Reference](https://pkg.go.dev/github.com/booscaaa/go-paginate/v3) - API documentation
- [📝 Examples](https://github.com/booscaaa/go-paginate/tree/master/v3) - Complete example code

## 🔄 Migration from v2 to v3

### v2 Code (Old)

```go
params, err := paginate.PaginQuery(
    paginate.WithStruct(User{}),
    paginate.WithTable("users"),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(10),
    paginate.WithSearch("john"),
)
sql, args := paginate.GenerateSQL(params)
```

### v3 Code (New - Fluent API)

```go
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Page(2).
    Limit(10).
    Search("john", "name", "email").
    BuildSQL()
```

### v3 Code (Compatibility)

```go
// The traditional API still works!
p, err := paginate.NewPaginator(
    paginate.WithStruct(User{}),
    paginate.WithTable("users"),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(10),
    paginate.WithSearch("john"),
)
sql, args := p.GenerateSQL()
```

## 🤝 Contributing

Contributions are very welcome! Feel free to:

- 🐛 Report bugs
- 💡 Suggest new features
- 📝 Improve documentation
- 🔧 Submit pull requests
- ⭐ Star the project

### How to Contribute

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Development

```bash
# Clone the repository
git clone https://github.com/booscaaa/go-paginate.git
cd go-paginate/v3

# Install dependencies
go mod tidy

# Run tests
go test -v ./paginate

# Run examples
go run example_builder.go
```

## 📞 Contact

This project is part of my personal portfolio. I'll be happy to receive feedback about the project, code, structure, or anything that could make me a better developer!

- 📧 Email: boscardinvinicius@gmail.com
- 💼 LinkedIn: [booscaaa](https://www.linkedin.com/in/booscaaa/)
- 🐙 GitHub: [booscaaa](https://github.com/booscaaa)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/booscaaa/go-paginate/blob/master/LICENSE) file for details.

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/booscaaa">Vinícius Boscardin</a>
</p>
