# Go Paginate v3 - New Operators & Filters Guide

<p align="center">
  <img src="https://raw.githubusercontent.com/booscaaa/go-paginate/master/assets/icon.png" alt="Go Paginate Logo" width="150">
</p>

<p align="center">
  <h1 align="center">🚀 Go Paginate v3 - New Operators Guide</h1>
  <p align="center">
    <strong>Complete guide to the 7 new powerful operators added in v3</strong>
  </p>
</p>

---

## 🎯 What's New in v3

Go Paginate v3 introduces **7 powerful new operators** that significantly expand your querying capabilities:

- ✨ **`Like`** - Simple LIKE pattern matching
- ✨ **`Eq`** - Simple equality operator
- ✨ **`In`** - IN clause for multiple values
- ✨ **`NotIn`** - NOT IN clause for exclusions
- ✨ **`Between`** - BETWEEN clause for ranges
- ✨ **`IsNull`** - IS NULL checks
- ✨ **`IsNotNull`** - IS NOT NULL checks

These operators complement the existing advanced operators (`LikeOr`, `LikeAnd`, `EqOr`, `EqAnd`, `Gte`, `Gt`, `Lte`, `Lt`) to give you **22+ total filtering options**.

---

## 📚 Complete Operator Reference

### 🆕 New Simple Operators

#### `Like` - Pattern Matching
```go
// Builder API
builder.WhereLike("name", "john%")

// Functional API
paginate.WithLike(map[string][]string{
    "name": {"john%", "jane%"},
    "email": {"%@company.com"},
})

// HTTP Query
// ?like[name]=john%&like[email]=%@company.com
```

#### `Eq` - Simple Equality
```go
// Builder API
builder.WhereEquals("status", "active")

// Functional API
paginate.WithEq(map[string][]any{
    "status": {"active", "pending"},
    "role_id": {1, 2, 3},
})

// HTTP Query
// ?eq[status]=active&eq[status]=pending&eq[role_id]=1
```

#### `In` - Multiple Value Matching
```go
// Builder API
builder.WhereIn("age", 25, 30, 35, 40)

// Functional API
paginate.WithIn(map[string][]any{
    "age": {25, 30, 35, 40},
    "department_id": {1, 2, 3},
})

// HTTP Query
// ?in[age]=25&in[age]=30&in[age]=35&in[department_id]=1
```

#### `NotIn` - Exclusion Matching
```go
// Builder API
builder.WhereNotIn("status", "deleted", "banned", "suspended")

// Functional API
paginate.WithNotIn(map[string][]any{
    "status": {"deleted", "banned"},
    "role_id": {99, 100}, // Exclude admin roles
})

// HTTP Query
// ?notin[status]=deleted&notin[status]=banned&notin[role_id]=99
```

#### `Between` - Range Queries
```go
// Builder API
builder.WhereBetween("age", 18, 65)
builder.WhereBetween("salary", 50000, 150000)
builder.WhereBetween("created_at", "2023-01-01", "2023-12-31")

// Functional API
paginate.WithBetween(map[string][2]any{
    "age": {18, 65},
    "salary": {50000, 150000},
    "created_at": {"2023-01-01", "2023-12-31"},
})

// HTTP Query
// ?between[age][0]=18&between[age][1]=65&between[salary][0]=50000&between[salary][1]=150000
```

#### `IsNull` - Null Value Checks
```go
// Builder API
builder.WhereIsNull("deleted_at")
builder.WhereIsNull("archived_at")

// Functional API
paginate.WithIsNull([]string{"deleted_at", "archived_at"})

// HTTP Query
// ?isnull=deleted_at&isnull=archived_at
```

#### `IsNotNull` - Non-Null Value Checks
```go
// Builder API
builder.WhereIsNotNull("email")
builder.WhereIsNotNull("phone")

// Functional API
paginate.WithIsNotNull([]string{"email", "phone", "verified_at"})

// HTTP Query
// ?isnotnull=email&isnotnull=phone&isnotnull=verified_at
```

---

## 🔥 Real-World Examples

### Example 1: User Management with New Operators

```go
type User struct {
    ID        int     `json:"id" paginate:"users.id"`
    Name      string  `json:"name" paginate:"users.name"`
    Email     string  `json:"email" paginate:"users.email"`
    Age       int     `json:"age" paginate:"users.age"`
    Status    string  `json:"status" paginate:"users.status"`
    RoleID    *int    `json:"role_id" paginate:"users.role_id"`
    Salary    int     `json:"salary" paginate:"users.salary"`
    DeletedAt *string `json:"deleted_at" paginate:"users.deleted_at"`
}

func advancedUserSearch() {
    // Using all new operators together
    sql, args, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        // Simple pattern matching
        WhereLike("name", "john%").
        WhereLike("email", "%@company.com").
        // Multiple value matching
        WhereIn("age", 25, 30, 35, 40).
        // Exclusions
        WhereNotIn("status", "deleted", "banned", "suspended").
        // Range queries
        WhereBetween("salary", 50000, 150000).
        // Null checks
        WhereIsNull("deleted_at").  // Only active users
        WhereIsNotNull("email").   // Must have email
        WhereIsNotNull("role_id"). // Must have role assigned
        OrderBy("name").
        Page(1).
        Limit(25).
        BuildSQL()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("SQL: %s\n", sql)
    fmt.Printf("Args: %v\n", args)
    
    // Output SQL:
    // SELECT * FROM users WHERE 
    //   (name::TEXT ILIKE $1) AND 
    //   (email::TEXT ILIKE $2) AND 
    //   age IN ($3, $4, $5, $6) AND 
    //   status NOT IN ($7, $8, $9) AND 
    //   salary BETWEEN $10 AND $11 AND 
    //   deleted_at IS NULL AND 
    //   email IS NOT NULL AND 
    //   role_id IS NOT NULL 
    // ORDER BY name ASC LIMIT 25 OFFSET 0
}
```

### Example 2: E-commerce Product Filtering

```go
type Product struct {
    ID          int     `json:"id" paginate:"products.id"`
    Name        string  `json:"name" paginate:"products.name"`
    Description string  `json:"description" paginate:"products.description"`
    Price       float64 `json:"price" paginate:"products.price"`
    CategoryID  int     `json:"category_id" paginate:"products.category_id"`
    Brand       string  `json:"brand" paginate:"products.brand"`
    InStock     bool    `json:"in_stock" paginate:"products.in_stock"`
    Rating      float64 `json:"rating" paginate:"products.rating"`
    Tags        string  `json:"tags" paginate:"products.tags"`
    DiscountID  *int    `json:"discount_id" paginate:"products.discount_id"`
}

func productSearch() {
    // Complex product filtering
    params, err := paginate.NewPaginator(
        paginate.WithStruct(Product{}),
        paginate.WithTable("products"),
        
        // Search for electronics or gadgets
        paginate.WithLike(map[string][]string{
            "name": {"%electronics%", "%gadget%"},
            "tags": {"%smartphone%", "%laptop%"},
        }),
        
        // Specific categories
        paginate.WithIn(map[string][]any{
            "category_id": {1, 2, 3, 5}, // Electronics categories
        }),
        
        // Exclude certain brands
        paginate.WithNotIn(map[string][]any{
            "brand": {"BrandX", "BrandY", "Discontinued"},
        }),
        
        // Price range
        paginate.WithBetween(map[string][2]any{
            "price": {100.0, 2000.0},
            "rating": {3.5, 5.0},
        }),
        
        // Must be in stock and have description
        paginate.WithEq(map[string][]any{
            "in_stock": {true},
        }),
        paginate.WithIsNotNull([]string{"description"}),
        
        // Optional: products with discounts
        // paginate.WithIsNotNull([]string{"discount_id"}),
        
        paginate.WithPage(1),
        paginate.WithLimit(20),
        paginate.WithSort([]string{"-rating", "price", "name"}),
    )

    if err != nil {
        log.Fatal(err)
    }

    sql, args := params.GenerateSQL()
    fmt.Printf("Product Search SQL: %s\n", sql)
    fmt.Printf("Args: %v\n", args)
}
```

### Example 3: HTTP API Integration

```go
func handleUserSearch(w http.ResponseWriter, r *http.Request) {
    // Example URL:
    // /users?like[name]=john%&in[age]=25&in[age]=30&in[age]=35&notin[status]=deleted&notin[status]=banned&between[salary][0]=50000&between[salary][1]=150000&isnull=deleted_at&isnotnull=email&page=1&limit=25
    
    // Automatically bind query parameters
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, "Invalid query parameters", http.StatusBadRequest)
        return
    }

    // Build query with bound parameters
    sql, args, err := paginate.NewBuilder().
        Table("users u").
        Model(&User{}).
        LeftJoin("departments d", "u.dept_id = d.id").
        LeftJoin("roles r", "u.role_id = r.id").
        FromStruct(params).
        BuildSQL()

    if err != nil {
        http.Error(w, "Query build error", http.StatusInternalServerError)
        return
    }

    // Execute query
    rows, err := db.Query(sql, args...)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Process results...
    var users []User
    for rows.Next() {
        var user User
        // Scan into user struct...
        users = append(users, user)
    }

    // Return JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "users": users,
        "total": len(users),
        "page":  params.Page,
        "limit": params.Limit,
    })
}
```

---

## 🔄 Migration from v2

All new operators are **100% backward compatible**. Your existing v2 code will continue to work without any changes.

### Adding New Operators to Existing Code

```go
// v2 code (still works)
builder := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    WhereEquals("status", "active").
    WhereGreaterThan("age", 18)

// v3 enhancements (add these)
builder = builder.
    WhereLike("name", "john%").           // New!
    WhereIn("department_id", 1, 2, 3).    // New!
    WhereNotIn("role", "admin", "super"). // New!
    WhereBetween("salary", 50000, 100000). // New!
    WhereIsNull("deleted_at").            // New!
    WhereIsNotNull("email")               // New!
```

---

## 📊 Performance Notes

### Optimized SQL Generation
All new operators generate optimized, parameterized SQL:

```sql
-- Like operator
WHERE name::TEXT ILIKE $1

-- In operator  
WHERE age IN ($1, $2, $3, $4)

-- NotIn operator
WHERE status NOT IN ($1, $2, $3)

-- Between operator
WHERE salary BETWEEN $1 AND $2

-- Null checks
WHERE deleted_at IS NULL
WHERE email IS NOT NULL
```

### Best Practices

1. **Use `In` instead of multiple `Eq` conditions**
   ```go
   // ❌ Less efficient
   builder.WhereEquals("status", "active").
           WhereEquals("status", "pending")
   
   // ✅ More efficient
   builder.WhereIn("status", "active", "pending")
   ```

2. **Use `Between` for ranges**
   ```go
   // ❌ Less efficient
   builder.WhereGreaterThanOrEqual("age", 18).
           WhereLessThanOrEqual("age", 65)
   
   // ✅ More efficient
   builder.WhereBetween("age", 18, 65)
   ```

3. **Use `IsNull`/`IsNotNull` for existence checks**
   ```go
   // ✅ Efficient null checks
   builder.WhereIsNull("deleted_at").      // Active records
           WhereIsNotNull("email").        // Must have email
           WhereIsNotNull("verified_at")   // Verified users
   ```

---

## 🧪 Testing Your Queries

Use the built-in debug mode to see generated SQL:

```go
// Enable debug mode
paginate.SetDebugMode(true)

// Your query will now log the generated SQL
sql, args, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    WhereLike("name", "john%").
    WhereIn("age", 25, 30, 35).
    WhereIsNull("deleted_at").
    BuildSQL()

// Output will include:
// [DEBUG] Generated SQL: SELECT * FROM users WHERE (name::TEXT ILIKE $1) AND age IN ($2, $3, $4) AND deleted_at IS NULL
// [DEBUG] Parameters: [john% 25 30 35]
```

---

## 📝 Complete HTTP Query Examples

### Simple Filtering
```
/users?like[name]=john%&eq[status]=active&isnotnull=email
```

### Complex Filtering
```
/products?like[name]=%phone%&in[category_id]=1&in[category_id]=2&in[category_id]=3&notin[brand]=BrandX&notin[brand]=BrandY&between[price][0]=100&between[price][1]=1000&eq[in_stock]=true&isnotnull=description&isnull=discontinued_at&sort=-rating&sort=price&page=1&limit=20
```

### User Management
```
/users?like[email]=%@company.com&in[department_id]=1&in[department_id]=2&notin[status]=deleted&notin[status]=suspended&between[age][0]=25&between[age][1]=55&between[salary][0]=50000&between[salary][1]=150000&isnull=deleted_at&isnotnull=role_id&sort=name&sort=-created_at
```

---

## 🎯 Summary

Go Paginate v3 adds **7 powerful new operators** that make your database queries more expressive and efficient:

- ✅ **22+ total operators** for comprehensive filtering
- ✅ **100% backward compatible** with v2
- ✅ **Optimized SQL generation** with parameterized queries
- ✅ **Full HTTP binding support** for web APIs
- ✅ **Type-safe** Go interfaces
- ✅ **Thoroughly tested** with comprehensive test coverage

Upgrade to v3 today and supercharge your pagination capabilities!

```bash
go get github.com/booscaaa/go-paginate/v3
```

---

<p align="center">
  <strong>Happy Paginating! 🚀</strong>
</p>