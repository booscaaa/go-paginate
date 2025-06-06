# Query Parameters Bind

This functionality allows binding URL query parameters to pagination structs in a simple and efficient way.

## Features

- ✅ Basic parameter binding (page, limit, search, etc.)
- ✅ Array and slice support
- ✅ Complex parameters with array syntax (`likeor[field]`, `eqor[field]`, etc.)
- ✅ Automatic type conversion (int, bool, string)
- ✅ Custom struct support
- ✅ Type validation
- ✅ Default values

## Basic Usage

### 1. Bind to PaginationParams (default struct)

```go
package main

import (
    "fmt"
    "log"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

func main() {
    // From a query string
    queryString := "page=2&limit=25&search=john&search_fields=name,email"
    params, err := paginate.BindQueryStringToStruct(queryString)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Page: %d\n", params.Page)           // 2
    fmt.Printf("Limit: %d\n", params.Limit)         // 25
    fmt.Printf("Search: %s\n", params.Search)       // "john"
    fmt.Printf("Fields: %v\n", params.SearchFields) // ["name", "email"]
}
```

### 2. Bind using url.Values

```go
import (
    "net/url"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

func handler() {
    queryParams := url.Values{
        "page":   {"3"},
        "limit":  {"50"},
        "search": {"admin"},
        "vacuum": {"true"},
    }

    params, err := paginate.BindQueryParamsToStruct(queryParams)
    if err != nil {
        log.Fatal(err)
    }

    // Use params...
}
```

## Supported Parameters

### Basic Parameters

| Parameter         | Type     | Description                 | Example                        |
| ----------------- | -------- | --------------------------- | ------------------------------ |
| `page`            | int      | Page number                 | `page=2`                       |
| `limit`           | int      | Items per page              | `limit=25`                     |
| `items_per_page`  | int      | Alias for limit             | `items_per_page=25`            |
| `search`          | string   | Search term                 | `search=john`                  |
| `search_fields`   | []string | Fields for search           | `search_fields=name,email`     |
| `sort_columns`    | []string | Columns for sorting         | `sort_columns=name,created_at` |
| `sort_directions` | []string | Sort directions             | `sort_directions=ASC,DESC`     |
| `columns`         | []string | Columns for selection       | `columns=id,name,email`        |
| `vacuum`          | bool     | Use count estimation        | `vacuum=true`                  |
| `no_offset`       | bool     | Disable OFFSET              | `no_offset=false`              |

### Complex Parameters (Array Syntax)

| Parameter           | Type                | Description         | Example                                              |
| ------------------- | ------------------- | ------------------- | ---------------------------------------------------- |
| `likeor[field]`     | map[string][]string | OR search by field  | `likeor[status]=active&likeor[status]=pending`       |
| `likeand[field]`    | map[string][]string | AND search by field | `likeand[name]=john`                                 |
| `eqor[field]`       | map[string][]any    | OR equality         | `eqor[age]=25&eqor[age]=30`                          |
| `eqand[field]`      | map[string][]any    | AND equality        | `eqand[role]=admin`                                  |
| `gte[field]`        | map[string]any      | Greater or equal    | `gte[age]=18`                                        |
| `gt[field]`         | map[string]any      | Greater than        | `gt[score]=80`                                       |
| `lte[field]`        | map[string]any      | Less or equal       | `lte[price]=100.50`                                  |
| `lt[field]`         | map[string]any      | Less than           | `lt[date]=2023-12-31`                                |

## Advanced Examples

### 1. Complex Parameters

```go
queryString := "page=1&likeor[status]=active&likeor[status]=pending&eqor[age]=25&eqor[age]=30&gte[created_at]=2023-01-01"
params, err := paginate.BindQueryStringToStruct(queryString)

// Result:
// params.LikeOr["status"] = ["active", "pending"]
// params.EqOr["age"] = [25, 30]
// params.Gte["created_at"] = "2023-01-01"
```

### 2. Custom Struct

```go
type CustomParams struct {
    Page     int      `query:"page"`
    Limit    int      `query:"limit"`
    Search   string   `query:"search"`
    Filters  []string `query:"filters"`
    Active   bool     `query:"active"`
}

queryParams := url.Values{
    "page":    {"4"},
    "limit":   {"100"},
    "search":  {"custom"},
    "filters": {"filter1,filter2,filter3"},
    "active":  {"true"},
}

customParams := &CustomParams{}
err := paginate.BindQueryParams(queryParams, customParams)
```

### 3. Usage in HTTP Handler

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    // Extract query parameters from request
    queryParams := r.URL.Query()

    // Bind to pagination struct
    paginationParams, err := paginate.BindQueryParamsToStruct(queryParams)
    if err != nil {
        http.Error(w, "Invalid parameters", http.StatusBadRequest)
        return
    }

    // Use the parameters to build the query
    // ...
}
```

## Type Conversion

The library performs automatic type conversion:

- **Strings**: Used directly
- **Integers**: Converted with `strconv.Atoi()`
- **Booleans**: Converted with `strconv.ParseBool()`
- **Floats**: Converted with `strconv.ParseFloat()`
- **Slices**: Multiple values or comma-separated values

## Error Handling

- Invalid values are ignored (keeps default value)
- Incompatible types are ignored
- Query string parsing errors are returned
- Invalid targets (non-pointer or non-struct) return error

## Default Values

The `PaginationParams` struct has default values:

```go
params := &PaginationParams{
    Page:         1,  // default page
    Limit:        10, // default limit
    ItemsPerPage: 10, // default items per page
}
```

## Run Example

To see the functionality in action:

```bash
go run example_bind.go
```

## Run Tests

```bash
go test -v ./paginate -run TestBind
```
