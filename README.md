<p align="center">
  <img src="https://raw.githubusercontent.com/booscaaa/go-paginate/master/assets/icon.png" alt="Go Paginate Logo" width="200">
</p>

<p align="center">
  <h1 align="center">Go Paginate v3 - The Ultimate Go Pagination Library</h1>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/booscaaa/go-paginate/v3"><img alt="Reference" src="https://img.shields.io/badge/go-reference-purple?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/booscaaa/go-paginate.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge"></a>
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/booscaaa/go-paginate/test.yaml?style=for-the-badge">
    <a href="https://codecov.io/gh/booscaaa/go-paginate"><img alt="Coverage" src="https://img.shields.io/codecov/c/github/booscaaa/go-paginate/master.svg?style=for-the-badge"></a>
    <img alt="Go Version" src="https://img.shields.io/badge/go-1.24.2+-blue?style=for-the-badge">
  </p>
</p>

<br>

## 🌟 Why Go Paginate v3?

Go Paginate v3 is the **most powerful and flexible** Go pagination library available. It provides three distinct APIs to fit any use case, from simple REST APIs to complex enterprise applications.

### ✨ Key Features

- 🚀 **3 Powerful APIs**: Fluent Builder, Automatic Binding, Traditional
- 🔍 **Advanced Filtering**: 15+ filter types including LikeOr, LikeAnd, EqOr, EqAnd, Gte, Gt, Lte, Lt
- 🔗 **Automatic HTTP Binding**: Convert query parameters to structs automatically
- 📊 **Complex Joins**: Full support for INNER, LEFT, RIGHT JOINs
- 🎯 **Type Safety**: Compile-time validation and runtime type checking
- 🔄 **100% Backward Compatible**: Seamless migration from v2
- ⚡ **High Performance**: Optimized SQL generation with minimal allocations
- 🛡️ **SQL Injection Safe**: Parameterized queries by default
- 📱 **Modern Sorting**: Support for `sort=name` and `sort=-created_at` patterns
- 🧪 **Thoroughly Tested**: 80%+ test coverage

---

## 📦 Installation

```bash
go get github.com/booscaaa/go-paginate/v3
```

**Requirements**: Go 1.24.2+

---

## 🚀 Quick Start - Advanced Example

```go
package main

import (
    "fmt"
    "log"
    "net/url"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

// User model with comprehensive paginate tags
type User struct {
    ID          int     `json:"id" paginate:"u.id"`
    Name        string  `json:"name" paginate:"u.name"`
    Email       string  `json:"email" paginate:"u.email"`
    Age         int     `json:"age" paginate:"u.age"`
    Status      string  `json:"status" paginate:"u.status"`
    Salary      float64 `json:"salary" paginate:"u.salary"`
    DeptID      int     `json:"dept_id" paginate:"u.dept_id"`
    DeptName    string  `json:"dept_name" paginate:"d.name"`
    CreatedAt   string  `json:"created_at" paginate:"u.created_at"`
    UpdatedAt   string  `json:"updated_at" paginate:"u.updated_at"`
    IsActive    bool    `json:"is_active" paginate:"u.is_active"`
    LastLogin   string  `json:"last_login" paginate:"u.last_login"`
}

// Custom pagination parameters struct
type AdvancedSearchParams struct {
    Page         int                 `json:"page"`
    Limit        int                 `json:"limit"`
    Search       string              `json:"search"`
    SearchFields []string            `json:"search_fields"`
    Sort         []string            `json:"sort"`
    LikeOr       map[string][]string `json:"likeor"`
    LikeAnd      map[string][]string `json:"likeand"`
    EqOr         map[string][]any    `json:"eqor"`
    EqAnd        map[string][]any    `json:"eqand"`
    Gte          map[string]any      `json:"gte"`
    Gt           map[string]any      `json:"gt"`
    Lte          map[string]any      `json:"lte"`
    Lt           map[string]any      `json:"lt"`
    Vacuum       bool                `json:"vacuum"`
}

func main() {
    fmt.Println("=== 🚀 Go Paginate v3 - Advanced Quick Start ===")
    fmt.Println()

    // 🌟 Example 1: Complex Fluent API with Joins
    fmt.Println("1. 🔥 Complex Fluent API with Multiple Joins:")
    complexFluentExample()
    fmt.Println()

    // 🌟 Example 2: FromJSON - Perfect for REST APIs
    fmt.Println("2. 📄 FromJSON - Dynamic Query from JSON:")
    fromJSONExample()
    fmt.Println()

    // 🌟 Example 3: FromStruct - Type-safe parameter binding
    fmt.Println("3. 🏗️ FromStruct - Type-safe Parameter Binding:")
    fromStructExample()
    fmt.Println()

    // 🌟 Example 4: Query String Binding - HTTP Integration
    fmt.Println("4. 🌐 Query String Binding - HTTP Integration:")
    queryStringBindingExample()
    fmt.Println()

    // 🌟 Example 5: Ultimate Complex Query
    fmt.Println("5. 🎯 Ultimate Complex Query - All Features Combined:")
    ultimateComplexExample()
}

func complexFluentExample() {
    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        LeftJoin("user_profiles p", "u.id = p.user_id").
        Page(2).
        Limit(25).
        Search("john", "name", "email").
        LikeOr("status", "active", "pending", "verified").
        LikeAnd("email", "@company.com").
        WhereEquals("is_active", true).
        WhereIn("dept_id", 1, 2, 3, 5).
        WhereGreaterThan("age", 21).
        WhereLessThanOrEqual("salary", 150000).
        WhereBetween("created_at", "2023-01-01", "2024-12-31").
        OrderBy("dept_name").
        OrderByDesc("salary").
        OrderBy("name").
        Vacuum().
        BuildSQL()

    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   SQL: %s\n", sql)
    fmt.Printf("   Args: %v\n", args)

    // Output:
    // SQL: SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at, u.updated_at, 
    // u.is_active, u.last_login FROM users u INNER JOIN departments d ON u.dept_id = d.id 
    // LEFT JOIN user_profiles p ON u.id = p.user_id WHERE (u.name ILIKE $1 OR u.email ILIKE $2) 
    // AND (u.status ILIKE $3 OR u.status ILIKE $4 OR u.status ILIKE $5) AND (u.email ILIKE $6) 
    // AND u.is_active = $7 AND u.dept_id IN ($8, $9, $10, $11) AND u.age > $12 AND u.salary <= $13 
    // AND u.created_at BETWEEN $14 AND $15 ORDER BY d.name ASC, u.salary DESC, u.name ASC LIMIT 25 OFFSET 25
    
    // Args: [%john% %john% %active% %pending% %verified% %@company.com% true 1 2 3 5 21 150000 2023-01-01 2024-12-31]
}

func fromJSONExample() {
    // Complex JSON query - perfect for REST API endpoints
    jsonQuery := `{
        "page": 1,
        "limit": 50,
        "search": "engineer",
        "search_fields": ["name", "email", "dept_name"],
        "likeor": {
            "status": ["active", "pending", "on_leave"],
            "dept_name": ["Engineering", "DevOps", "QA"]
        },
        "likeand": {
            "email": ["@company.com"]
        },
        "eqor": {
            "age": [25, 30, 35, 40],
            "dept_id": [1, 2, 3]
        },
        "gte": {
            "salary": 50000,
            "age": 22
        },
        "lte": {
            "salary": 200000,
            "last_login": "2024-12-31"
        },
        "gt": {
            "created_at": "2020-01-01"
        },
        "lt": {
            "updated_at": "2024-12-31"
        },
        "sort": ["-salary", "dept_name", "-created_at"],
        "vacuum": true
    }`

    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        FromJSON(jsonQuery).
        BuildSQL()

    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   JSON: %s\n", jsonQuery)
    fmt.Printf("   SQL: %s\n", sql)
    fmt.Printf("   Args: %v\n", args)

    // Output:
     // SQL: SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at,
     //  u.updated_at, u.is_active, u.last_login FROM users u INNER JOIN departments d ON u.dept_id = d.id 
     // WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3) AND (u.status ILIKE $4 OR u.status ILIKE $5 
     // OR u.status ILIKE $6 OR d.name ILIKE $7 OR d.name ILIKE $8 OR d.name ILIKE $9) AND (u.email ILIKE $10) 
     // AND (u.age = $11 OR u.age = $12 OR u.age = $13 OR u.age = $14 OR u.dept_id = $15 OR u.dept_id = $16 OR u.dept_id = $17) 
     // AND u.salary >= $18 AND u.age >= $19 AND u.created_at > $20 AND u.salary <= $21 AND u.last_login <= $22 
     // AND u.updated_at < $23 ORDER BY u.salary DESC, d.name ASC, u.created_at DESC LIMIT 50 OFFSET 0
     
     // Args: [%engineer% %engineer% %engineer% %active% %pending% %on_leave% %Engineering% %DevOps% %QA% %@company.com% 25 30 35 40 1 2 3 50000 22 2020-01-01 200000 2024-12-31 2024-12-31]
}

func fromStructExample() {
    // Create complex search parameters struct
    searchParams := &AdvancedSearchParams{
        Page:         3,
        Limit:        30,
        Search:       "senior",
        SearchFields: []string{"name", "email", "dept_name"},
        Sort:         []string{"-salary", "name", "-age"},
        LikeOr: map[string][]string{
            "status":    {"active", "verified"},
            "dept_name": {"Engineering", "Product", "Design"},
        },
        LikeAnd: map[string][]string{
            "email": {"@company.com"},
        },
        EqOr: map[string][]any{
            "age":     {28, 32, 35, 40},
            "dept_id": {1, 2, 4, 7},
        },
        EqAnd: map[string][]any{
            "is_active": {true},
        },
        Gte: map[string]any{
            "salary":     75000,
            "age":        25,
            "created_at": "2021-01-01",
        },
        Gt: map[string]any{
            "last_login": "2024-01-01",
        },
        Lte: map[string]any{
            "salary":     250000,
            "updated_at": "2024-12-31",
        },
        Lt: map[string]any{
            "age": 60,
        },
        Vacuum: true,
    }

    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        LeftJoin("user_roles ur", "u.id = ur.user_id").
        FromStruct(searchParams).
        BuildSQL()

    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   Struct: %+v\n", searchParams)
    fmt.Printf("   SQL: %s\n", sql)
    fmt.Printf("   Args: %v\n", args)

    // Output:
     // SQL: SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at, 
     // u.updated_at, u.is_active, u.last_login FROM users u 
     // INNER JOIN departments d ON u.dept_id = d.id LEFT JOIN user_roles ur ON u.id = ur.user_id 
     // WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3) AND (u.status ILIKE $4 OR u.status ILIKE $5 
     // OR d.name ILIKE $6 OR d.name ILIKE $7 OR d.name ILIKE $8) AND (u.email ILIKE $9) AND (u.age = $10 OR u.age = $11 
     // OR u.age = $12 OR u.age = $13 OR u.dept_id = $14 OR u.dept_id = $15 OR u.dept_id = $16 OR u.dept_id = $17) 
     // AND u.is_active = $18 AND u.salary >= $19 AND u.age >= $20 AND u.created_at >= $21 AND u.last_login > $22 
     // AND u.salary <= $23 AND u.updated_at <= $24 AND u.age < $25 ORDER BY u.salary DESC, u.name ASC, u.age DESC LIMIT 30 OFFSET 60
     
     // Args: [%senior% %senior% %senior% %active% %verified% %Engineering% %Product% %Design% %@company.com% 28 32 35 40 1 2 4 7 true 75000 25 2021-01-01 2024-01-01 250000 2024-12-31 60]
}

func queryStringBindingExample() {
    // Simulate complex HTTP query string (like from a web form or API call)
    queryString := "page=2&limit=40&search=developer&search_fields=name,email,dept_name" +
        "&likeor[status]=active&likeor[status]=pending&likeor[dept_name]=Engineering" +
        "&likeand[email]=@company.com&eqor[age]=25&eqor[age]=30&eqor[age]=35" +
        "&eqand[is_active]=true&gte[salary]=60000&gte[age]=23&lte[salary]=180000" +
        "&gt[created_at]=2022-01-01&lt[updated_at]=2024-12-31" +
        "&sort=-salary&sort=dept_name&sort=-created_at&vacuum=true"

    // Method 1: Bind to default PaginationParams struct
    params, err := paginate.BindQueryStringToStruct(queryString)
    if err != nil {
        log.Printf("❌ Error binding query string: %v", err)
        return
    }

    fmt.Printf("   Query String: %s\n", queryString)
    fmt.Printf("   Bound Params: %+v\n", params)

    // Method 2: Use bound parameters with builder
    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        FromStruct(params).
        BuildSQL()

    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   SQL: %s\n", sql)
    fmt.Printf("   Args: %v\n", args)

    // Method 3: Direct url.Values binding
    queryParams, _ := url.ParseQuery(queryString)
    directParams, err := paginate.BindQueryParamsToStruct(queryParams)
    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   Direct Bound: %+v\n", directParams)

    // Output:
    // SQL: SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at, 
    // u.updated_at, u.is_active, u.last_login FROM users u INNER JOIN departments d ON u.dept_id = d.id 
    // WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3) AND (u.status ILIKE $4 OR u.status ILIKE $5 OR d.name ILIKE $6) AND (u.email ILIKE $7) AND (u.age = $8 
    // OR u.age = $9 OR u.age = $10) AND u.is_active = $11 AND u.salary >= $12 AND u.age >= $13 
    // AND u.created_at > $14 AND u.salary <= $15 AND u.updated_at < $16 ORDER BY u.salary DESC, d.name ASC, u.created_at DESC LIMIT 40 OFFSET 40
    
    // Args: [%developer% %developer% %developer% %active% %pending% %Engineering% %@company.com% 25 30 35 true 60000 23 2022-01-01 180000 2024-12-31]
}

func ultimateComplexExample() {
    // The most complex query possible - showcasing ALL features
    ultimateJSON := `{
        "page": 1,
        "limit": 100,
        "search": "tech lead",
        "search_fields": ["name", "email", "dept_name", "role_name"],
        "likeor": {
            "status": ["active", "verified", "premium"],
            "dept_name": ["Engineering", "DevOps", "Architecture", "Platform"],
            "role_name": ["Senior", "Lead", "Principal", "Staff"]
        },
        "likeand": {
            "email": ["@company.com"],
            "skills": ["golang", "kubernetes"]
        },
        "eqor": {
            "experience_level": ["senior", "lead", "principal"],
            "team_size": [5, 8, 12, 15],
            "location_id": [1, 2, 3, 5, 8]
        },
        "eqand": {
            "is_active": [true],
            "is_verified": [true],
            "has_security_clearance": [true]
        },
        "gte": {
            "salary": 120000,
            "age": 28,
            "years_experience": 5,
            "team_size": 3,
            "created_at": "2020-01-01",
            "performance_score": 8.5
        },
        "gt": {
            "last_promotion_date": "2022-01-01",
            "last_review_score": 4.0
        },
        "lte": {
            "salary": 350000,
            "age": 55,
            "updated_at": "2024-12-31"
        },
        "lt": {
            "days_since_last_login": 30,
            "open_tickets": 5
        },
        "sort": ["-salary", "-performance_score", "dept_name", "name", "-created_at"],
        "vacuum": true
    }`

    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        InnerJoin("user_roles ur", "u.id = ur.user_id").
        InnerJoin("roles r", "ur.role_id = r.id").
        LeftJoin("user_skills us", "u.id = us.user_id").
        LeftJoin("skills s", "us.skill_id = s.id").
        LeftJoin("locations l", "u.location_id = l.id").
        FromJSON(ultimateJSON).
        BuildSQL()

    if err != nil {
        log.Printf("❌ Error: %v", err)
        return
    }

    fmt.Printf("   🎯 Ultimate Complex Query Generated!\n")
    fmt.Printf("   📊 Features Used: Multi-table joins, complex filtering, advanced sorting\n")
    fmt.Printf("   🔍 Search Types: OR/AND search, equals OR/AND, range filters (gte,gt,lte,lt)\n")
    fmt.Printf("   📄 JSON Length: %d characters\n", len(ultimateJSON))
    fmt.Printf("   🗃️ SQL Length: %d characters\n", len(sql))
    fmt.Printf("   📝 Parameters: %d args\n", len(args))
    fmt.Printf("   \n   SQL: %s\n", sql)
    fmt.Printf("   Args: %v\n", args)

    // Also generate count query
    countSQL, countArgs, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        InnerJoin("departments d", "u.dept_id = d.id").
        InnerJoin("user_roles ur", "u.id = ur.user_id").
        InnerJoin("roles r", "ur.role_id = r.id").
        LeftJoin("user_skills us", "u.id = us.user_id").
        LeftJoin("skills s", "us.skill_id = s.id").
        LeftJoin("locations l", "u.location_id = l.id").
        FromJSON(ultimateJSON).
        BuildCountSQL()

    if err != nil {
        log.Printf("❌ Count Error: %v", err)
        return
    }

    fmt.Printf("   \n   Count SQL: %s\n", countSQL)
     fmt.Printf("   Count Args: %v\n", countArgs)

     // Output:
     // SQL: SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at, u.updated_at, 
     // u.is_active, u.last_login FROM users u 
     // INNER JOIN departments d ON u.dept_id = d.id INNER JOIN user_roles ur ON u.id = ur.user_id 
     // INNER JOIN roles r ON ur.role_id = r.id LEFT JOIN user_skills us ON u.id = us.user_id 
     // LEFT JOIN skills s ON us.skill_id = s.id LEFT JOIN locations l ON u.location_id = l.id 
     // WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3 OR r.name ILIKE $4) AND (u.status ILIKE $5 OR u.status ILIKE $6 
     // OR u.status ILIKE $7 OR d.name ILIKE $8 OR d.name ILIKE $9 OR d.name ILIKE $10 OR d.name ILIKE $11 OR r.name ILIKE $12 
     // OR r.name ILIKE $13 OR r.name ILIKE $14 OR r.name ILIKE $15) AND (u.email ILIKE $16 AND s.name ILIKE $17 AND s.name ILIKE $18) 
     // AND (u.experience_level = $19 OR u.experience_level = $20 OR u.experience_level = $21 OR u.team_size = $22 OR u.team_size = $23 
     // OR u.team_size = $24 OR u.team_size = $25 OR u.location_id = $26 OR u.location_id = $27 OR u.location_id = $28 OR u.location_id = $29 
     // OR u.location_id = $30) AND u.is_active = $31 AND u.is_verified = $32 AND u.has_security_clearance = $33 AND u.salary >= $34 AND u.age >= $35 
     // AND u.years_experience >= $36 AND u.team_size >= $37 AND u.created_at >= $38 AND u.performance_score >= $39 AND u.last_promotion_date > $40 
     // AND u.last_review_score > $41 AND u.salary <= $42 AND u.age <= $43 AND u.updated_at <= $44 AND u.days_since_last_login < $45 
     // AND u.open_tickets < $46 ORDER BY u.salary DESC, u.performance_score DESC, d.name ASC, u.name ASC, u.created_at DESC LIMIT 100 OFFSET 0
     
     // Args: [%tech lead% %tech lead% %tech lead% %tech lead% %active% %verified% %premium% %Engineering% %DevOps% %Architecture% %Platform% %Senior% %Lead% %Principal% %Staff% %@company.com% %golang% %kubernetes% senior lead principal 5 8 12 15 1 2 3 5 8 true true true 120000 28 5 3 2020-01-01 8.5 2022-01-01 4.0 350000 55 2024-12-31 30 5]

     // Count SQL: SELECT COUNT(*) FROM users u INNER JOIN departments d ON u.dept_id = d.id INNER JOIN user_roles ur ON u.id = ur.user_id 
     // INNER JOIN roles r ON ur.role_id = r.id LEFT JOIN user_skills us ON u.id = us.user_id LEFT JOIN skills s ON us.skill_id = s.id 
     // LEFT JOIN locations l ON u.location_id = l.id WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3 OR r.name ILIKE $4) 
     // AND (u.status ILIKE $5 OR u.status ILIKE $6 OR u.status ILIKE $7 OR d.name ILIKE $8 OR d.name ILIKE $9 OR d.name ILIKE $10 OR d.name 
     // ILIKE $11 OR r.name ILIKE $12 OR r.name ILIKE $13 OR r.name ILIKE $14 OR r.name ILIKE $15) AND (u.email ILIKE $16 AND s.name ILIKE $17 
     // AND s.name ILIKE $18) AND (u.experience_level = $19 OR u.experience_level = $20 OR u.experience_level = $21 OR u.team_size = $22 
     // OR u.team_size = $23 OR u.team_size = $24 OR u.team_size = $25 OR u.location_id = $26 OR u.location_id = $27 OR u.location_id = $28 
     // OR u.location_id = $29 OR u.location_id = $30) AND u.is_active = $31 AND u.is_verified = $32 AND u.has_security_clearance = $33 
     // AND u.salary >= $34 AND u.age >= $35 AND u.years_experience >= $36 AND u.team_size >= $37 AND u.created_at >= $38 
     // AND u.performance_score >= $39 AND u.last_promotion_date > $40 AND u.last_review_score > $41 AND u.salary <= $42 
     // AND u.age <= $43 AND u.updated_at <= $44 AND u.days_since_last_login < $45 AND u.open_tickets < $46

     // Count Args: [%tech lead% %tech lead% %tech lead% %tech lead% %active% %verified% %premium% %Engineering% %DevOps% %Architecture% %Platform% %Senior% %Lead% %Principal% %Staff% %@company.com% %golang% %kubernetes% senior lead principal 5 8 12 15 1 2 3 5 8 true true true 120000 28 5 3 2020-01-01 8.5 2022-01-01 4.0 350000 55 2024-12-31 30 5]
 }
```

### 🎯 Key Features Demonstrated:

- **🔥 Complex Joins**: INNER, LEFT joins across multiple tables
- **📄 FromJSON**: Dynamic queries from JSON (perfect for REST APIs)
- **🏗️ FromStruct**: Type-safe parameter binding from structs
- **🌐 Query String Binding**: Automatic HTTP query parameter parsing
- **🔍 Advanced Search**: LikeOr, LikeAnd with multiple fields
- **⚖️ Flexible Filtering**: EqOr, EqAnd, Gte, Gt, Lte, Lt
- **📊 Smart Sorting**: Modern `-field` syntax for DESC ordering
- **🎯 Range Queries**: Between, In, complex comparisons
- **⚡ Performance**: Vacuum mode for optimized queries
- **🛡️ SQL Safety**: Parameterized queries prevent injection

### 📈 Output Example:
```sql
SELECT u.id, u.name, u.email, u.age, u.status, u.salary, u.dept_id, d.name, u.created_at, u.updated_at, u.is_active, u.last_login 
FROM users u 
INNER JOIN departments d ON u.dept_id = d.id 
INNER JOIN user_roles ur ON u.id = ur.user_id 
INNER JOIN roles r ON ur.role_id = r.id 
LEFT JOIN user_skills us ON u.id = us.user_id 
LEFT JOIN skills s ON us.skill_id = s.id 
LEFT JOIN locations l ON u.location_id = l.id 
WHERE (u.name ILIKE $1 OR u.email ILIKE $2 OR d.name ILIKE $3 OR r.name ILIKE $4) 
AND (u.status ILIKE $5 OR d.name ILIKE $6 OR r.name ILIKE $7) 
AND (u.email ILIKE $8 AND s.name ILIKE $9) 
AND (u.experience_level = $10 OR u.team_size = $11 OR u.location_id = $12) 
AND u.is_active = $13 AND u.is_verified = $14 AND u.has_security_clearance = $15 
AND u.salary >= $16 AND u.age >= $17 AND u.years_experience >= $18 
AND u.last_promotion_date > $19 AND u.last_review_score > $20 
AND u.salary <= $21 AND u.age <= $22 AND u.updated_at <= $23 
AND u.days_since_last_login < $24 AND u.open_tickets < $25 
ORDER BY u.salary DESC, u.performance_score DESC, d.name ASC, u.name ASC, u.created_at DESC 
LIMIT 100 OFFSET 0
```

---

## 🎯 Three Powerful APIs

### 1. 🌟 Fluent API (Recommended)

**Perfect for**: New projects, complex queries, readable code

```go
// Basic usage
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Page(1).
    Limit(25).
    Search("john", "name", "email").
    OrderBy("name").
    BuildSQL()

// Advanced filtering
sql, args, err := paginate.NewBuilder().
    Table("users u").
    Model(&User{}).
    Select("u.*", "d.name as dept_name").
    LeftJoin("departments d", "u.dept_id = d.id").
    WhereEquals("u.status", "active").
    WhereIn("u.role", "admin", "manager").
    WhereGreaterThanOrEqual("u.age", 21).
    WhereLessThanOrEqual("u.salary", 150000).
    WhereBetween("u.created_at", "2023-01-01", "2023-12-31").
    LikeOr("u.name", "John", "Jane").
    LikeAnd("u.email", "@company.com").
    OrderBy("d.name", "u.name").
    OrderByDesc("u.created_at").
    BuildSQL()
```

### 2. 🔗 Automatic HTTP Binding

**Perfect for**: REST APIs, web frameworks, HTTP handlers

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    // Automatically convert query parameters to struct
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    // Use with Fluent API
    sql, args, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        FromStruct(params).  // 🎯 Magic happens here!
        BuildSQL()

    // Execute query...
}

// Example URLs that work automatically:
// /users?page=2&limit=25&search=john&search_fields=name,email
// /users?likeor[status]=active&likeor[status]=pending
// /users?eqor[age]=25&eqor[age]=30&gte[salary]=50000
// /users?sort=name&sort=-created_at  // New sorting pattern!
```

### 3. 🔧 Traditional API (Compatibility)

**Perfect for**: Migrating from v2, existing codebases

```go
// Traditional API - still fully supported
p, err := paginate.NewPaginator(
    paginate.WithTable("users"),
    paginate.WithStruct(&User{}),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(20),
    paginate.WithSearch("john"),
    paginate.WithSearchFields([]string{"name", "email"}),
    paginate.WithEqOr(map[string][]any{
        "status": {"active", "pending"},
    }),
    paginate.WithGte(map[string]any{"age": 18}),
)

if err != nil {
    log.Fatal(err)
}

sql, args := p.GenerateSQL()
countSQL, countArgs := p.GenerateCountQuery()
```

---

## 🔍 Complete Filter Reference

### Basic Filters

| Method                      | Description                   | Example                                  |
| --------------------------- | ----------------------------- | ---------------------------------------- |
| `Search(term, fields...)`   | Search across multiple fields | `Search("john", "name", "email")`        |
| `WhereEquals(field, value)` | Exact match                   | `WhereEquals("status", "active")`        |
| `Where(clause, args...)`    | Custom WHERE clause           | `Where("age BETWEEN $1 AND $2", 18, 65)` |

### Advanced Equality Filters

| Method                            | Description | Example                                        |
| --------------------------------- | ----------- | ---------------------------------------------- |
| `WhereIn(field, values...)`       | IN clause   | `WhereIn("role", "admin", "manager")`          |
| `WhereEqOr(field, values...)` | OR equality | `WhereEqOr("status", "active", "pending")` |

### Comparison Filters

| Method                                  | Description                | Example                                                  |
| --------------------------------------- | -------------------------- | -------------------------------------------------------- |
| `WhereGreaterThan(field, value)`        | Greater than (>)           | `WhereGreaterThan("age", 18)`                            |
| `WhereGreaterThanOrEqual(field, value)` | Greater than or equal (>=) | `WhereGreaterThanOrEqual("salary", 50000)`               |
| `WhereLessThan(field, value)`           | Less than (<)              | `WhereLessThan("age", 65)`                               |
| `WhereLessThanOrEqual(field, value)`    | Less than or equal (<=)    | `WhereLessThanOrEqual("price", 100)`                     |
| `WhereBetween(field, min, max)`         | Between values             | `WhereBetween("created_at", "2023-01-01", "2023-12-31")` |

### Search Filters

| Method                        | Description           | Example                              |
| ----------------------------- | --------------------- | ------------------------------------ |
| `LikeOr(field, values...)`  | Search with OR logic  | `LikeOr("name", "John", "Jane")`   |
| `LikeAnd(field, values...)` | Search with AND logic | `LikeAnd("email", "@company.com")` |

### Join Operations

| Method                        | Description | Example                                              |
| ----------------------------- | ----------- | ---------------------------------------------------- |
| `LeftJoin(table, condition)`  | LEFT JOIN   | `LeftJoin("departments d", "u.dept_id = d.id")`      |
| `InnerJoin(table, condition)` | INNER JOIN  | `InnerJoin("roles r", "u.role_id = r.id")`           |
| `RightJoin(table, condition)` | RIGHT JOIN  | `RightJoin("profiles p", "u.id = p.user_id")`        |
| `Join(clause)`                | Custom JOIN | `Join("FULL OUTER JOIN logs l ON u.id = l.user_id")` |

### Sorting & Ordering

| Method                          | Description              | Example                                       |
| ------------------------------- | ------------------------ | --------------------------------------------- |
| `OrderBy(column, direction...)` | Sort ascending (default) | `OrderBy("name")` or `OrderBy("name", "ASC")` |
| `OrderByDesc(column)`           | Sort descending          | `OrderByDesc("created_at")`                   |

### Utility Methods

| Method               | Description          | Example                         |
| -------------------- | -------------------- | ------------------------------- |
| `Select(columns...)` | Custom SELECT        | `Select("id", "name", "email")` |
| `WithoutOffset()`    | Disable OFFSET       | `WithoutOffset()`               |
| `WithVacuum()`       | Use count estimation | `WithVacuum()`                  |

---

## 📋 HTTP Query Parameter Reference

### Basic Parameters

| Parameter       | Type     | Description           | Example                     |
| --------------- | -------- | --------------------- | --------------------------- |
| `page`          | int      | Page number (1-based) | `?page=2`                   |
| `limit`         | int      | Items per page        | `?limit=25`                 |
| `search`        | string   | Search term           | `?search=john`              |
| `search_fields` | []string | Fields to search      | `?search_fields=name,email` |
| `vacuum`        | bool     | Use count estimation  | `?vacuum=true`              |
| `no_offset`     | bool     | Disable OFFSET        | `?no_offset=true`           |

### Modern Sorting (v3 New!)

| Parameter | Type     | Description                       | Example                       |
| --------- | -------- | --------------------------------- | ----------------------------- |
| `sort`    | []string | Sort fields (prefix `-` for DESC) | `?sort=name&sort=-created_at` |

### Legacy Sorting (Backward Compatible)

| Parameter         | Type     | Description     | Example                         |
| ----------------- | -------- | --------------- | ------------------------------- |
| `sort_columns`    | []string | Columns to sort | `?sort_columns=name,created_at` |
| `sort_directions` | []string | Sort directions | `?sort_directions=ASC,DESC`     |

### Advanced Filters

| Parameter           | Type                | Description           | Example                                               |
| ------------------- | ------------------- | --------------------- | ----------------------------------------------------- |
| `likeor[field]`     | map[string][]string | OR search             | `?likeor[status]=active&likeor[status]=pending`       |
| `likeand[field]`    | map[string][]string | AND search            | `?likeand[name]=admin`                                |
| `eqor[field]`       | map[string][]any    | OR equality           | `?eqor[age]=25&eqor[age]=30`                          |
| `eqand[field]`      | map[string][]any    | AND equality          | `?eqand[role]=admin`                                  |
| `gte[field]`        | map[string]any      | Greater than or equal | `?gte[age]=18`                                        |
| `gt[field]`         | map[string]any      | Greater than          | `?gt[score]=80`                                       |
| `lte[field]`        | map[string]any      | Less than or equal    | `?lte[price]=100.50`                                  |
| `lt[field]`         | map[string]any      | Less than             | `?lt[date]=2023-12-31`                                |

---

## 🎨 Real-World Examples

### Example 1: E-commerce Product Search

```go
type Product struct {
    ID          int     `json:"id" paginate:"products.id"`
    Name        string  `json:"name" paginate:"products.name"`
    Description string  `json:"description" paginate:"products.description"`
    Price       float64 `json:"price" paginate:"products.price"`
    CategoryID  int     `json:"category_id" paginate:"products.category_id"`
    Category    string  `json:"category" paginate:"categories.name"`
    InStock     bool    `json:"in_stock" paginate:"products.in_stock"`
    CreatedAt   string  `json:"created_at" paginate:"products.created_at"`
}

func searchProducts(w http.ResponseWriter, r *http.Request) {
    // URL: /products?search=laptop&gte[price]=500&lte[price]=2000&eqor[category_id]=1&eqor[category_id]=2&sort=price&sort=-created_at

    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    sql, args, err := paginate.NewBuilder().
        Table("products p").
        Model(&Product{}).
        Select("p.*", "c.name as category").
        LeftJoin("categories c", "p.category_id = c.id").
        WhereEquals("p.in_stock", true).
        FromStruct(params).
        BuildSQL()

    if err != nil {
        http.Error(w, "Query build error", http.StatusInternalServerError)
        return
    }

    // Execute query and return results...
    rows, err := db.Query(sql, args...)
    // ... handle results
}
```

### Example 2: User Management with Complex Filters

```go
type UserWithDepartment struct {
    ID           int    `json:"id" paginate:"u.id"`
    Name         string `json:"name" paginate:"u.name"`
    Email        string `json:"email" paginate:"u.email"`
    Age          int    `json:"age" paginate:"u.age"`
    Salary       int    `json:"salary" paginate:"u.salary"`
    Status       string `json:"status" paginate:"u.status"`
    DepartmentID int    `json:"department_id" paginate:"u.department_id"`
    Department   string `json:"department" paginate:"d.name"`
    Role         string `json:"role" paginate:"r.name"`
    CreatedAt    string `json:"created_at" paginate:"u.created_at"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    // Complex query with multiple joins and filters
    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&UserWithDepartment{}).
        Select("u.*", "d.name as department", "r.name as role").
        InnerJoin("departments d", "u.department_id = d.id").
        LeftJoin("roles r", "u.role_id = r.id").
        WhereEquals("u.status", "active").
        WhereIn("d.type", "engineering", "sales", "marketing").
        WhereGreaterThanOrEqual("u.age", 21).
        WhereLessThanOrEqual("u.salary", 200000).
        WhereBetween("u.created_at", "2023-01-01", "2023-12-31").
        LikeOr("u.name", "John", "Jane", "Admin").
        LikeAnd("u.email", "@company.com").
        OrderBy("d.name").
        OrderBy("u.name").
        OrderByDesc("u.created_at").
        Page(1).
        Limit(50).
        BuildSQL()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Generated SQL will be something like:
    // SELECT u.id, u.name, u.email, u.age, u.salary, u.status, u.department_id, d.name as department, r.name as role, u.created_at
    // FROM users u
    // INNER JOIN departments d ON u.department_id = d.id
    // LEFT JOIN roles r ON u.role_id = r.id
    // WHERE u.status = $1 AND d.type IN ($2, $3, $4) AND u.age >= $5 AND u.salary <= $6
    // AND u.created_at BETWEEN $7 AND $8
    // AND (u.name ILIKE $9 OR u.name ILIKE $10 OR u.name ILIKE $11)
    // AND u.email ILIKE $12
    // ORDER BY d.name ASC, u.name ASC, u.created_at DESC
    // LIMIT 50 OFFSET 0

    rows, err := db.Query(sql, args...)
    // ... process results
}
```

### Example 3: JSON API with FromJSON

```go
func searchFromJSON(w http.ResponseWriter, r *http.Request) {
    // Accept JSON payload for complex searches
    jsonQuery := `{
        "page": 1,
        "limit": 20,
        "search": "john",
        "search_fields": ["name", "email"],
        "eqor": {
            "status": ["active", "pending"],
            "role": ["admin", "manager"]
        },
        "likeor": {
            "name": ["John", "Jane"],
            "email": ["@company.com", "@gmail.com"]
        },
        "gte": {
            "age": 18,
            "salary": 50000
        },
        "lte": {
            "salary": 150000
        },
        "sort": ["-created_at", "name"]
    }`

    sql, args, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        FromJSON(jsonQuery).
        BuildSQL()

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute query...
}
```

### Example 4: Custom Struct Binding

```go
type CustomSearchParams struct {
    Page        int      `query:"page"`
    Limit       int      `query:"limit"`
    Search      string   `query:"search"`
    Category    string   `query:"category"`
    Tags        []string `query:"tags"`
    MinPrice    float64  `query:"min_price"`
    MaxPrice    float64  `query:"max_price"`
    InStock     bool     `query:"in_stock"`
    Featured    bool     `query:"featured"`
    SortBy      string   `query:"sort_by"`
    SortOrder   string   `query:"sort_order"`
}

func customSearch(w http.ResponseWriter, r *http.Request) {
    customParams := &CustomSearchParams{
        Page:     1,
        Limit:    20,
        InStock:  true,  // default values
        SortBy:   "name",
        SortOrder: "asc",
    }

    err := paginate.BindQueryParams(r.URL.Query(), customParams)
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

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

    builder = builder.WhereEquals("in_stock", customParams.InStock)

    if customParams.Featured {
        builder = builder.WhereEquals("featured", true)
    }

    // Apply sorting
    if customParams.SortOrder == "desc" {
        builder = builder.OrderByDesc(customParams.SortBy)
    } else {
        builder = builder.OrderBy(customParams.SortBy)
    }

    sql, args, err := builder.BuildSQL()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute query...
}
```

---

## 🔄 Migration Guide

### From v2 to v3

#### v2 Code (Old)

```go
params, err := paginate.PaginQuery(
    paginate.WithStruct(User{}),
    paginate.WithTable("users"),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(10),
    paginate.WithSearch("john"),
    paginate.WithSearchFields([]string{"name", "email"}),
)
sql, args := paginate.GenerateSQL(params)
```

#### v3 Code (New - Fluent API)

```go
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Page(2).
    Limit(10).
    Search("john", "name", "email").
    BuildSQL()
```

#### v3 Code (Compatibility Mode)

```go
// The traditional API still works!
p, err := paginate.NewPaginator(
    paginate.WithStruct(User{}),
    paginate.WithTable("users"),
    paginate.WithPage(2),
    paginate.WithItemsPerPage(10),
    paginate.WithSearch("john"),
    paginate.WithSearchFields([]string{"name", "email"}),
)
sql, args := p.GenerateSQL()
```

### New Sorting Pattern

#### Old Way

```
?sort_columns=name,created_at&sort_directions=ASC,DESC
```

#### New Way (Recommended)

```
?sort=name&sort=-created_at
```

Both patterns work! The new pattern takes priority when both are present.

---

## 🎯 API Comparison

| Feature              | Fluent API | Traditional API | Automatic Binding |
| -------------------- | ---------- | --------------- | ----------------- |
| **Ease of use**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐          | ⭐⭐⭐⭐⭐        |
| **Type Safety**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐        | ⭐⭐⭐            |
| **Flexibility**      | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐      | ⭐⭐⭐            |
| **Performance**      | ⭐⭐⭐⭐   | ⭐⭐⭐⭐        | ⭐⭐⭐⭐          |
| **HTTP Integration** | ⭐⭐⭐     | ⭐⭐            | ⭐⭐⭐⭐⭐        |
| **Learning Curve**   | ⭐⭐⭐⭐   | ⭐⭐            | ⭐⭐⭐⭐⭐        |
| **Code Readability** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐          | ⭐⭐⭐⭐          |

### When to use each API:

- **🌟 Fluent API**: New projects, complex queries, better readability, type safety
- **🔧 Traditional API**: Migrating existing code, maximum flexibility, compatibility
- **🔗 Automatic Binding**: REST APIs, web frameworks, rapid development

---

## 🚀 Running Examples

Clone the repository and run the examples:

```bash
git clone https://github.com/booscaaa/go-paginate.git
cd go-paginate

# Run examples from the examples folder
cd examples/builder && go run main.go     # Fluent API examples
cd ../bind && go run main.go              # HTTP binding examples  
cd ../debug && go run main.go             # Debug mode examples
cd ../v2 && go run main.go                # Traditional API examples

# Or run from project root
go run examples/builder/main.go           # Fluent API examples
go run examples/bind/main.go              # HTTP binding examples
go run examples/debug/main.go             # Debug mode examples
go run examples/v2/main.go                # Traditional API examples
```

---

## 🧪 Testing

```bash
# Run all tests
go test -v ./paginate

# Run specific test suites
go test -v ./paginate -run TestBuilder
go test -v ./paginate -run TestBind
go test -v ./paginate -run TestPaginate

# Run with coverage
go test -v -cover ./paginate

# Generate coverage report
go test -coverprofile=coverage.out ./paginate
go tool cover -html=coverage.out
```

---

## 🔧 Configuration

### Environment Variables

```bash
# Enable debug mode (prints generated SQL)
export GO_PAGINATE_DEBUG=true

# Set default page size
export GO_PAGINATE_DEFAULT_LIMIT=25

# Set maximum page size
export GO_PAGINATE_MAX_LIMIT=1000
```

### Global Configuration

```go
package main

import "github.com/booscaaa/go-paginate/v3/paginate"

func init() {
    // Set global defaults
    paginate.SetDefaultLimit(25)
    paginate.SetMaxLimit(1000)
    paginate.SetDebugMode(true)
}
```

---

## 📊 Performance Benchmarks

```
BenchmarkFluentAPI-20          204253     5559 ns/op    4384 B/op   105 allocs/op
BenchmarkTraditionalAPI-20     269164     4057 ns/op    3130 B/op    85 allocs/op
BenchmarkAutomaticBinding-20   148874     7935 ns/op    5540 B/op   124 allocs/op
BenchmarkSQLGeneration-20      307086     3750 ns/op    2490 B/op    76 allocs/op
```

_Benchmarks run on 12th Gen Intel(R) Core(TM) i7-12700, Go 1.24.2_

---

## 🛡️ Security

### SQL Injection Protection

Go Paginate v3 uses **parameterized queries** by default, making it safe from SQL injection attacks:

```go
// ✅ Safe - uses parameterized queries
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Search(userInput, "name", "email").  // userInput is safely parameterized
    WhereEquals("status", userStatus).   // userStatus is safely parameterized
    BuildSQL()

// Generated SQL: SELECT ... WHERE (name ILIKE $1 OR email ILIKE $2) AND status = $3
    // Args: ["%userInput%", "%userInput%", "userStatus"]
```

### Input Validation

```go
// Built-in validation
builder := paginate.NewBuilder().
    Page(-1).  // ❌ Error: page must be greater than 0
    Limit(0)   // ❌ Error: limit must be greater than 0

sql, args, err := builder.BuildSQL()
if err != nil {
    // Handle validation errors
}
```

---

## 🤝 Contributing

Contributions are very welcome! Here's how you can help:

### Ways to Contribute

- 🐛 **Report bugs** - Found an issue? Let us know!
- 💡 **Suggest features** - Have an idea? We'd love to hear it!
- 📝 **Improve documentation** - Help make our docs better
- 🔧 **Submit pull requests** - Code contributions are appreciated
- ⭐ **Star the project** - Show your support!
- 📢 **Spread the word** - Tell others about Go Paginate

### Development Setup

```bash
# Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/go-paginate.git
cd go-paginate/v3

# Install dependencies
go mod tidy

# Run tests
go test -v ./paginate

# Run examples
go run example_builder.go

# Create a feature branch
git checkout -b feature/amazing-feature

# Make your changes and commit
git commit -m 'Add some amazing feature'

# Push to your fork
git push origin feature/amazing-feature

# Open a Pull Request
```

### Code Style

- Follow standard Go conventions
- Add tests for new features
- Update documentation
- Run `go fmt` and `go vet`
- Ensure all tests pass

---

## 📚 Additional Resources

- 📖 **[Complete Bind Documentation](v3/BIND_README.md)** - Detailed guide on query parameter binding
- 🔗 **[Go Reference](https://pkg.go.dev/github.com/booscaaa/go-paginate/v3)** - Complete API documentation
- 📝 **[Examples Repository](https://github.com/booscaaa/go-paginate/examples)** - More example code

---

### Support the Project

If Go Paginate v3 has been helpful to you, consider:

- ⭐ **Starring the repository**
- 🐛 **Reporting issues**
- 💡 **Suggesting improvements**
- 📢 **Sharing with others**

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](https://github.com/booscaaa/go-paginate/blob/master/LICENSE) file for details.

```
MIT License

Copyright (c) 2025 Vinícius Boscardin

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNES FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```


<p align="center">
  <a href="#-why-go-paginate-v3">⬆️ Back to top</a>
</p>
