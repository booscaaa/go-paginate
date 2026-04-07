<p align="center">
  <img src="https://raw.githubusercontent.com/booscaaa/go-paginate/master/assets/icon.png" alt="Go Paginate Logo" width="200">
</p>

<p align="center">
  <h1 align="center">Go Paginate v4 — The Ultimate Go Pagination Library</h1>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/booscaaa/go-paginate/v4"><img alt="Reference" src="https://img.shields.io/badge/go-reference-purple?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/booscaaa/go-paginate.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge"></a>
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/booscaaa/go-paginate/test.yaml?style=for-the-badge">
    <img alt="Go Version" src="https://img.shields.io/badge/go-1.21+-blue?style=for-the-badge">
  </p>
</p>

<br>

## Table of Contents

- [Why v4?](#why-v4)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Model Setup](#model-setup)
- [Global Configuration](#global-configuration)
- [Offset Pagination](#offset-pagination)
  - [Builder API](#builder-api)
  - [HTTP Binding](#http-binding)
  - [Page Response](#page-response)
- [Cursor Pagination](#cursor-pagination)
  - [How It Works](#how-it-works)
  - [Basic Cursor Usage](#basic-cursor-usage)
  - [Multi-Column Keyset (100% Stable)](#multi-column-keyset-100-stable)
  - [Cursor Response](#cursor-response)
  - [Frontend Integration](#frontend-integration)
- [Filtering Reference](#filtering-reference)
  - [AND Filters](#and-filters)
  - [OR Group Filters](#or-group-filters)
  - [Full-Text Search](#full-text-search)
  - [Raw WHERE Clauses](#raw-where-clauses)
- [Sorting](#sorting)
- [Joins](#joins)
- [Column Selection](#column-selection)
- [Schema Support](#schema-support)
- [Vacuum / Count Estimation](#vacuum--count-estimation)
- [FromJSON / FromMap / FromStruct](#fromjson--frommap--fromstruct)
- [Query String Binding](#query-string-binding)
- [Complete HTTP Handler Example](#complete-http-handler-example)
- [SQL Generation Reference](#sql-generation-reference)
- [API Reference](#api-reference)

---

## Why v4?

v4 is a **complete rewrite** focused on developer experience, generics, and production-grade cursor pagination.

| Feature | v3 | v4 |
|---|---|---|
| Generic response types | ✗ | ✅ `Page[T]`, `CursorPage[T]` |
| Cursor pagination | ✗ | ✅ single & multi-column keyset |
| Keyset seek method | ✗ | ✅ 100% stable with any sort |
| HATEOAS links | ✗ | ✅ built-in |
| Zero-boilerplate cursor | ✗ | ✅ one line in handler |
| HTTP query binding | ✅ | ✅ improved |
| 30+ filter types | ✅ | ✅ identical |
| OR grouping | ✅ | ✅ identical |
| Global config + env vars | ✅ | ✅ improved |

---

## Installation

```bash
go get github.com/booscaaa/go-paginate/v4
```

**Requirements**: Go 1.21+

---

## Quick Start

```go
package main

import (
    "net/http"
    "encoding/json"
    "github.com/booscaaa/go-paginate/v4/paginate"
)

type User struct {
    ID    int    `json:"id"    paginate:"users.id"`
    Name  string `json:"name"  paginate:"users.name"`
    Email string `json:"email" paginate:"users.email"`
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
    params, _ := paginate.BindQueryParamsToStruct(r.URL.Query())

    b := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        FromStruct(params)

    result, _ := b.Build()   // both queries at once

    // execute queries against your DB...
    var total int
    db.QueryRow(result.CountQuery, result.CountArgs...).Scan(&total)
    users := db.Query(result.Query, result.Args...)

    page := paginate.NewPage(users, total, params, r.URL)
    json.NewEncoder(w).Encode(page)
}
```

---

## Model Setup

Every model field needs two struct tags:

- **`json`** — the JSON key name used in query parameters (e.g. `?sort=created_at`)
- **`paginate`** — the actual SQL column reference (e.g. `table.column`)

```go
type Product struct {
    ID          int       `json:"id"           paginate:"p.id"`
    Name        string    `json:"name"         paginate:"p.name"`
    Price       float64   `json:"price"        paginate:"p.price"`
    CategoryID  int       `json:"category_id"  paginate:"p.category_id"`
    StockQty    int       `json:"stock_qty"    paginate:"p.stock_qty"`
    Active      bool      `json:"active"       paginate:"p.active"`
    CreatedAt   time.Time `json:"created_at"   paginate:"p.created_at"`
    UpdatedAt   time.Time `json:"updated_at"   paginate:"p.updated_at"`
    DeletedAt   *time.Time `json:"deleted_at"  paginate:"p.deleted_at"`
    // Fields from joined table
    CategoryName string   `json:"category_name" paginate:"c.name"`
}
```

> The `paginate` tag is resolved to the actual SQL column. If omitted, the `json` tag value is used as-is.

---

## Global Configuration

Configure once at application startup.

```go
func main() {
    // Programmatic configuration
    paginate.SetDefaultLimit(20)   // default items per page (default: 10)
    paginate.SetMaxLimit(200)      // maximum allowed limit (default: 100)
    paginate.SetDebugMode(true)    // log all generated SQL to slog

    // Custom logger
    paginate.SetLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

    // ...
}
```

**Environment variables** (loaded automatically at startup):

| Variable | Default | Description |
|---|---|---|
| `GO_PAGINATE_DEFAULT_LIMIT` | `10` | Default items per page |
| `GO_PAGINATE_MAX_LIMIT` | `100` | Maximum allowed limit |
| `GO_PAGINATE_DEBUG` | `false` | Enable SQL debug logging |

```bash
GO_PAGINATE_DEFAULT_LIMIT=25 GO_PAGINATE_MAX_LIMIT=500 ./myapp
```

---

## Offset Pagination

### Builder API

The fluent builder constructs SQL queries step by step.

```go
query, args, err := paginate.NewBuilder().
    Table("products").
    Schema("store").              // optional: generates FROM store.products
    Model(&Product{}).
    Page(2).
    Limit(25).
    Select("id", "name", "price"). // SELECT id, name, price (default: *)
    OrderBy("created_at", "DESC").
    OrderBy("id").                 // ASC is default
    Eq("active", true).
    WhereBetween("price", 10.0, 500.0).
    LeftJoin("categories c", "c.id = p.category_id").
    BuildSQL()
```

Generated SQL:
```sql
SELECT id, name, price
FROM store.products
LEFT JOIN categories c ON c.id = p.category_id
WHERE (p.active = $1) AND p.price BETWEEN $2 AND $3
ORDER BY p.created_at DESC, p.id ASC
LIMIT $4 OFFSET $5
```

#### Count only

```go
countQuery, countArgs, err := paginate.NewBuilder().
    Table("products").
    Model(&Product{}).
    Eq("active", true).
    BuildCountSQL()
// SELECT COUNT(p.id) FROM products WHERE (p.active = $1)
```

#### Both queries at once

```go
result, err := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Page(1).
    Limit(10).
    Eq("active", true).
    Build()
if err != nil {
    log.Fatal(err)
}

// result.Query      → paginated SELECT
// result.Args       → args for Query
// result.CountQuery → SELECT COUNT(...)
// result.CountArgs  → args for CountQuery
rows, _ := db.Query(result.Query, result.Args...)
var total int
db.QueryRow(result.CountQuery, result.CountArgs...).Scan(&total)
```

---

### HTTP Binding

Bind URL query parameters directly to pagination config in one call.

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    // r.URL.Query() = ?page=2&limit=20&sort=-created_at,name&eq[active]=true

    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    b := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        FromStruct(params)

    query, args, _ := b.BuildSQL()
    // ...
}
```

You can also bind from a raw query string:

```go
params, err := paginate.BindQueryStringToStruct("page=1&limit=10&sort=-name")
```

Or bind into your own struct:

```go
type MyParams struct {
    Page  int    `query:"page"`
    Limit int    `query:"limit"`
    Name  string `query:"name"`
}
var p MyParams
err := paginate.BindQueryParams(r.URL.Query(), &p)
```

---

### Page Response

`NewPage` builds a fully-featured HATEOAS response. All existing query params are preserved in the links; only `?page=N` is rewritten.

`page` and `perPage` are derived from `params` automatically — no need to call `CurrentPage()` / `CurrentLimit()`.

```go
page := paginate.NewPage(users, totalCount, params, r.URL)
```

**JSON output:**
```json
{
  "data": [...],
  "meta": {
    "current_page": 2,
    "per_page": 25,
    "total_items": 342,
    "total_pages": 14,
    "from": 26,
    "to": 50,
    "has_prev": true,
    "has_next": true
  },
  "links": {
    "self":  "/products?page=2&limit=25",
    "first": "/products?page=1&limit=25",
    "last":  "/products?page=14&limit=25",
    "prev":  "/products?page=1&limit=25",
    "next":  "/products?page=3&limit=25"
  }
}
```

---

## Cursor Pagination

Cursor pagination is **more efficient** than offset for large datasets — there is no `OFFSET` scan. Pages are stable even when data is inserted or deleted between requests.

### How It Works

1. Fetch **`limit + 1`** rows from the database
2. If `len(rows) > limit`, there is a next page — trim the slice back to `limit`
3. Encode the sort column values of the boundary rows into an opaque cursor token
4. Embed the token in `?cursor=<token>` links
5. On the next request, decode the token and inject a `WHERE` clause that seeks past the last seen row

### Basic Cursor Usage

**Single sort column** (simplest case):

```go
// repository
func ListUsers(ctx context.Context, params *paginate.PaginationParams) ([]User, error) {
    b := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        OrderBy("id").
        FromStruct(params)     // automatically applies cursor WHERE if params.Cursor is set

    query, args, _ := b.BuildSQL()
    return db.QueryUsers(ctx, query, args...) // fetch limit+1
}

// handler
func UsersHandler(w http.ResponseWriter, r *http.Request) {
    params, _ := paginate.BindQueryParamsToStruct(r.URL.Query())
    rawItems, _ := repo.ListUsers(r.Context(), params)
    page := paginate.NewCursorPage(rawItems, params, r.URL)
    json.NewEncoder(w).Encode(page)
}
```

**First request** (`GET /users?limit=10&sort=id`):
```sql
SELECT * FROM users ORDER BY users.id ASC LIMIT 11
```

**Second request** (`GET /users?limit=10&sort=id&cursor=<token>`):
```sql
SELECT * FROM users WHERE (users.id > $1) ORDER BY users.id ASC LIMIT 11
```

### Multi-Column Keyset (100% Stable)

With a single sort column, rows with duplicate values can be missed or repeated. The correct solution is the **seek method** with all sort columns.

```go
// repository
b := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    OrderBy("created_at", "DESC").
    OrderBy("id").             // tie-breaker — always add a unique column last
    FromStruct(params)

// handler — identical, zero extra code
page := paginate.NewCursorPage(rawItems, params, r.URL)
```

**Generated SQL for `?sort=-created_at,id&cursor=<token>`:**
```sql
SELECT * FROM users
WHERE (
    (users.created_at < $1)
    OR (users.created_at = $2 AND users.id > $3)
)
ORDER BY users.created_at DESC, users.id ASC
LIMIT $4
```

#### Operator matrix

| Pagination direction | Column sort | SQL operator |
|---|---|---|
| `after` (next page) | `ASC` | `>` |
| `after` (next page) | `DESC` | `<` |
| `before` (prev page) | `ASC` | `<` |
| `before` (prev page) | `DESC` | `>` |

#### Three sort columns example

```go
paginate.NewBuilder().
    Table("events").
    Model(&Event{}).
    OrderBy("year", "DESC").
    OrderBy("month", "DESC").
    OrderBy("id").
    FromStruct(params)
```

Generated keyset WHERE for `after`:
```sql
WHERE (
    (events.year < $1)
    OR (events.year = $2 AND events.month < $3)
    OR (events.year = $4 AND events.month = $5 AND events.id > $6)
)
```

#### Cursor is compatible with all filters

The cursor WHERE clause is added alongside all other filters. Nothing changes in the handler:

```go
// ?sort=-created_at,id&cursor=<token>&eq[active]=true&like[name]=john&gte[price]=10
b := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    OrderBy("created_at", "DESC").
    OrderBy("id").
    FromStruct(params)
```

```sql
SELECT * FROM users
WHERE (users.active = $1)
  AND (users.name::TEXT ILIKE $2)
  AND users.price >= $3
  AND ((users.created_at < $4) OR (users.created_at = $5 AND users.id > $6))
ORDER BY users.created_at DESC, users.id ASC
LIMIT $7
```

### Cursor Response

`NewCursorPage` accepts `limit+1` raw items and a `*PaginationParams`. It:

- Detects `hasNext` automatically (`len(rawItems) > limit`)
- Trims `data` to `limit` automatically
- Derives `hasPrev` from `params.Cursor != ""`
- Extracts boundary values from items via reflection on `json` tags
- Encodes multi-column tokens internally

```go
page := paginate.NewCursorPage(rawItems, params, r.URL)
```

**JSON output:**
```json
{
  "data": [...],
  "meta": {
    "per_page": 10,
    "has_next": true,
    "has_prev": true
  },
  "links": {
    "self": "/users?sort=-created_at%2Cid&limit=10",
    "next": "/users?sort=-created_at%2Cid&limit=10&cursor=eyJjb2xzIjpbImNyZWF0ZWRfYXQiLCJpZCJdLCJ2YWxzIjpbIjIwMjQtMDEtMTBUMTI6MDA6MDBaIiw0Ml0sImRpcnMiOlsiREVTQyIsIkFTQyJdLCJkaXIiOiJhZnRlciJ9",
    "prev": "/users?sort=-created_at%2Cid&limit=10&cursor=eyJjb2xzIjpbImNyZWF0ZWRfYXQiLCJpZCJdLCJ2YWxzIjpbIjIwMjQtMDEtMTBUMTI6MDA6MDBaIiwyMV0sImRpcnMiOlsiREVTQyIsIkFTQyJdLCJkaXIiOiJiZWZvcmUifQ=="
  }
}
```

The cursor token is **opaque** — the frontend treats it as a black box string and never needs to parse it.

### Frontend Integration

```js
async function fetchPage(cursor = null) {
  const params = new URLSearchParams(window.location.search)
  if (cursor) {
    params.set('cursor', cursor)
  } else {
    params.delete('cursor')
  }

  const res = await fetch(`/users?${params}`)
  const page = await res.json()

  renderTable(page.data)
  updateButtons(page.meta, page.links)
}

function updateButtons(meta, links) {
  const getCursor = (url) => url ? new URL(url).searchParams.get('cursor') : null

  document.getElementById('btn-prev').disabled = !meta.has_prev
  document.getElementById('btn-next').disabled = !meta.has_next

  document.getElementById('btn-prev').onclick = () => fetchPage(getCursor(links.prev))
  document.getElementById('btn-next').onclick = () => fetchPage(getCursor(links.next))
}

// load first page
fetchPage()
```

### Manual Cursor Encoding

For cases where you need to build cursor tokens manually (e.g., deep-linking to a specific position):

```go
// encode
token := paginate.EncodeCursor("id", 42, "after")

// decode
column, value, direction, err := paginate.DecodeCursor(token)
```

---

## Filtering Reference

All filters are driven by the `json` tag name of the model field and resolved to the `paginate` tag column in SQL.

### AND Filters

These filters are combined with `AND`.

#### Equality — `Eq`

Matches rows where the column equals **any** of the given values (implicit OR within the same field).

```go
// Builder
.Eq("status", "active")
.Eq("status", "active", "pending")   // status = 'active' OR status = 'pending'

// Query string
// ?eq[status]=active
// ?eq[status]=active&eq[status]=pending
```

#### Equality AND — `EqAnd`

All values must match (useful for array/tag fields):

```go
.EqAnd("role", "admin", "editor")
// role = 'admin' AND role = 'editor'
```

#### Greater / Less Than

```go
.WhereGreaterThan("age", 18)           // age > 18
.WhereGreaterThanOrEqual("price", 0)   // price >= 0
.WhereLessThan("stock_qty", 5)         // stock_qty < 5
.WhereLessThanOrEqual("price", 999.99) // price <= 999.99

// Query string
// ?gt[age]=18
// ?gte[price]=0
// ?lt[stock_qty]=5
// ?lte[price]=999.99
```

#### IN / NOT IN

```go
.In("category_id", 1, 2, 3)           // category_id IN (1, 2, 3)
.NotIn("status", "deleted", "banned")  // status NOT IN ('deleted', 'banned')

// Query string
// ?in[category_id]=1&in[category_id]=2&in[category_id]=3
// ?notin[status]=deleted&notin[status]=banned
```

#### BETWEEN

```go
.WhereBetween("price", 10.0, 500.0)   // price BETWEEN 10.0 AND 500.0

// Query string
// ?between[price]=10&between[price]=500
```

#### LIKE (ILIKE)

Case-insensitive substring match. Multiple values are OR'd within the same field.

```go
.WhereLike("name", "john")             // name ILIKE '%john%'
.WhereLike("name", "john", "jane")     // name ILIKE '%john%' OR name ILIKE '%jane%'

// Query string
// ?like[name]=john
```

#### LIKE AND

All patterns must match:

```go
.LikeAnd("name", "john", "doe")       // name ILIKE '%john%' AND name ILIKE '%doe%'

// Query string
// ?likeand[name]=john&likeand[name]=doe
```

#### IS NULL / IS NOT NULL

```go
.WhereIsNull("deleted_at")             // deleted_at IS NULL
.WhereIsNotNull("verified_at")         // verified_at IS NOT NULL

// Query string
// ?isnull=deleted_at
// ?isnotnull=verified_at
```

### OR Group Filters

All `*Or` variants are collected into a single `(... OR ...)` group that is AND'd with the rest of the WHERE clause. This lets you express "match any of these conditions" across multiple fields.

```go
paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    Eq("active", true).               // AND active = true
    LikeOr("name", "john").           // \
    EqOr("status", "vip", "premium"). //  > OR group
    GteOr("age", 21).                 // /
    BuildSQL()
```

```sql
WHERE (users.active = $1)
  AND (
    users.name::TEXT ILIKE $2
    OR users.status = $3
    OR users.status = $4
    OR users.age >= $5
  )
```

#### Complete OR group reference

| Method | Query string key | SQL |
|---|---|---|
| `LikeOr(field, values...)` | `likeor[field]=val` | `field ILIKE '%val%'` |
| `EqOr(field, values...)` | `eqor[field]=val` | `field = val` |
| `GteOr(field, value)` | `gteor[field]=val` | `field >= val` |
| `GtOr(field, value)` | `gtor[field]=val` | `field > val` |
| `LteOr(field, value)` | `lteor[field]=val` | `field <= val` |
| `LtOr(field, value)` | `ltor[field]=val` | `field < val` |
| `InOr(field, values...)` | `inor[field]=val` | `field IN (...)` |
| `NotInOr(field, values...)` | `notinor[field]=val` | `field NOT IN (...)` |
| `WhereIsNullOr(field)` | `isnullor=field` | `field IS NULL` |
| `WhereIsNotNullOr(field)` | `isnotnullor=field` | `field IS NOT NULL` |

### Full-Text Search

Search across multiple fields simultaneously:

```go
.Search("john doe", "name", "email", "bio")
```

```sql
WHERE (
    users.name::TEXT ILIKE '%john doe%'
    OR users.email::TEXT ILIKE '%john doe%'
    OR users.bio::TEXT ILIKE '%john doe%'
)
```

Query string:
```
?search=john+doe&search_fields=name,email,bio
```

### Raw WHERE Clauses

For conditions that don't fit the filter API:

```go
.Where("users.score > users.threshold")
.Where("users.expires_at > NOW()")
.Where("users.metadata->>'plan' = $1", "enterprise")
```

Multiple `.Where()` calls are joined with `AND` by default.

---

## Sorting

```go
// Single column, ASC (default)
.OrderBy("name")

// Single column, explicit direction
.OrderBy("created_at", "DESC")
.OrderByDesc("created_at")           // shorthand

// Multiple columns
.OrderBy("created_at", "DESC").
.OrderBy("id")
// ORDER BY users.created_at DESC, users.id ASC
```

**Query string — modern syntax** (preferred):

```
?sort=name             → ORDER BY name ASC
?sort=-created_at      → ORDER BY created_at DESC
?sort=-created_at,id   → ORDER BY created_at DESC, id ASC
```

**Query string — legacy syntax** (also supported):

```
?sort_columns=created_at,id&sort_directions=DESC,ASC
```

---

## Joins

```go
// Explicit JOIN string
.Join("LEFT JOIN categories c ON c.id = p.category_id")

// Convenience helpers
.LeftJoin("categories c", "c.id = p.category_id")
.InnerJoin("order_items oi", "oi.product_id = p.id")
.RightJoin("warehouses w", "w.id = p.warehouse_id")

// Multiple joins
paginate.NewBuilder().
    Table("orders o").
    Model(&Order{}).
    LeftJoin("users u", "u.id = o.user_id").
    LeftJoin("products p", "p.id = o.product_id").
    InnerJoin("statuses s", "s.id = o.status_id").
    Eq("status", "shipped").
    BuildSQL()
```

---

## Column Selection

```go
// Select specific columns
.Select("id", "name", "email", "created_at")
// SELECT id, name, email, created_at FROM ...

// Select with expressions
.Select("u.id", "u.name", "c.name AS category_name", "COUNT(*) AS total")

// Default: SELECT *
```

Query string:
```
?columns=id,name,email
```

---

## Schema Support

```go
paginate.NewBuilder().
    Schema("public").
    Table("users").
    Model(&User{}).
    BuildSQL()
// FROM public.users
```

---

## Vacuum / Count Estimation

For very large tables, PostgreSQL's `count_estimate` function is significantly faster than `COUNT(*)`. Enable it with `.WithVacuum()`:

```go
countQuery, countArgs, _ := paginate.NewBuilder().
    Table("events").
    Model(&Event{}).
    Eq("type", "click").
    WithVacuum().
    BuildCountSQL()
// SELECT count_estimate('SELECT COUNT(events.id) FROM events WHERE ...')
```

> Requires the `count_estimate` function to be installed in your PostgreSQL database.

---

## FromJSON / FromMap / FromStruct

Populate a builder from an external source:

```go
// From JSON string
b := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    FromJSON(`{"page":2,"limit":10,"sort":["-created_at","id"],"eq":{"active":[true]}}`)

// From map
b := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    FromMap(map[string]any{
        "page":  2,
        "limit": 10,
        "eq":    map[string]any{"active": []any{true}},
    })

// From struct (most common — used with HTTP binding)
params, _ := paginate.BindQueryParamsToStruct(r.URL.Query())
b := paginate.NewBuilder().
    Table("users").
    Model(&User{}).
    FromStruct(params)
```

All three methods support the complete filter/sort/cursor surface. Fields set on the builder *before* `FromStruct` can be overridden by params, and fields set *after* always take precedence.

---

## Query String Binding

`PaginationParams` is bound from URL query parameters. Every filter and sort option has a corresponding query string key.

### Complete Query String Reference

| Parameter | Example | Description |
|---|---|---|
| `page` | `?page=2` | Page number (offset pagination) |
| `limit` | `?limit=25` | Items per page |
| `items_per_page` | `?items_per_page=25` | Alias for limit |
| `sort` | `?sort=-created_at,id` | Sort columns (`-` prefix = DESC) |
| `sort_columns` | `?sort_columns=name,age` | Legacy sort columns |
| `sort_directions` | `?sort_directions=ASC,DESC` | Legacy sort directions |
| `columns` | `?columns=id,name` | SELECT specific columns |
| `search` | `?search=john` | Full-text search term |
| `search_fields` | `?search_fields=name,email` | Fields to search in |
| `vacuum` | `?vacuum=true` | Use count_estimate |
| `no_offset` | `?no_offset=true` | Skip OFFSET clause |
| `cursor` | `?cursor=<token>` | Cursor pagination token |
| `eq[field]` | `?eq[status]=active` | Equality (OR within field) |
| `eqand[field]` | `?eqand[role]=admin` | Equality AND |
| `eqor[field]` | `?eqor[status]=vip` | Equality in OR group |
| `like[field]` | `?like[name]=john` | ILIKE match |
| `likeand[field]` | `?likeand[name]=john` | ILIKE AND |
| `likeor[field]` | `?likeor[name]=john` | ILIKE in OR group |
| `gte[field]` | `?gte[age]=18` | >= |
| `gt[field]` | `?gt[price]=0` | > |
| `lte[field]` | `?lte[price]=999` | <= |
| `lt[field]` | `?lt[stock]=5` | < |
| `in[field]` | `?in[id]=1&in[id]=2` | IN |
| `notin[field]` | `?notin[status]=deleted` | NOT IN |
| `between[field]` | `?between[price]=10&between[price]=500` | BETWEEN |
| `isnull` | `?isnull=deleted_at` | IS NULL |
| `isnotnull` | `?isnotnull=verified_at` | IS NOT NULL |
| `gteor[field]` | `?gteor[age]=21` | >= in OR group |
| `gtor[field]` | `?gtor[score]=100` | > in OR group |
| `lteor[field]` | `?lteor[price]=50` | <= in OR group |
| `ltor[field]` | `?ltor[qty]=10` | < in OR group |
| `inor[field]` | `?inor[tag]=go` | IN in OR group |
| `notinor[field]` | `?notinor[status]=spam` | NOT IN in OR group |
| `isnullor` | `?isnullor=archived_at` | IS NULL in OR group |
| `isnotnullor` | `?isnotnullor=paid_at` | IS NOT NULL in OR group |

---

## Complete HTTP Handler Example

A production-ready handler combining offset pagination, cursor pagination, filters, joins, and response building.

```go
package handlers

import (
    "context"
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/booscaaa/go-paginate/v4/paginate"
)

type Product struct {
    ID           int     `json:"id"            paginate:"p.id"`
    Name         string  `json:"name"          paginate:"p.name"`
    Price        float64 `json:"price"         paginate:"p.price"`
    Stock        int     `json:"stock"         paginate:"p.stock"`
    Active       bool    `json:"active"        paginate:"p.active"`
    CategoryName string  `json:"category_name" paginate:"c.name"`
    CreatedAt    string  `json:"created_at"    paginate:"p.created_at"`
}

// --- Offset pagination ---

// GET /products?page=2&limit=20&sort=-created_at&eq[active]=true&gte[price]=10&like[name]=shirt
func ListProductsOffset(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        b := paginate.NewBuilder().
            Table("products p").
            Model(&Product{}).
            LeftJoin("categories c", "c.id = p.category_id").
            FromStruct(params)

        result, err := b.Build()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        products := scanProducts(db, result.Query, result.Args)
        var total int
        db.QueryRowContext(r.Context(), result.CountQuery, result.CountArgs...).Scan(&total)

        page := paginate.NewPage(products, total, params, r.URL)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(page)
    }
}

// --- Cursor pagination ---

// GET /products/cursor?limit=20&sort=-created_at,id&eq[active]=true&cursor=<token>
func ListProductsCursor(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        b := paginate.NewBuilder().
            Table("products p").
            Model(&Product{}).
            LeftJoin("categories c", "c.id = p.category_id").
            OrderBy("created_at", "DESC").
            OrderBy("id").            // tie-breaker
            FromStruct(params)        // cursor WHERE injected automatically

        query, args, err := b.BuildSQL()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Fetch limit+1 rows to probe hasNext
        rawItems := scanProducts(db, query, args)

        page := paginate.NewCursorPage(rawItems, params, r.URL)
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(page)
    }
}

func scanProducts(db *sql.DB, query string, args []any) []Product {
    rows, _ := db.Query(query, args...)
    defer rows.Close()
    var products []Product
    for rows.Next() {
        var p Product
        rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Active, &p.CategoryName, &p.CreatedAt)
        products = append(products, p)
    }
    return products
}
```

---

## SQL Generation Reference

### Placeholder style

v4 generates PostgreSQL-style `$1, $2, $3` placeholders. All values are passed as arguments — no string interpolation is ever performed.

### SELECT clause

```go
// Default
SELECT * FROM users ...

// With .Select(...)
SELECT u.id, u.name, c.name AS category_name FROM users u ...
```

### WHERE clause construction order

1. Full-text `Search` across fields
2. `Like` (AND group)
3. `LikeAnd` (AND group)
4. `Eq` (AND group, OR within field)
5. `EqAnd` (AND group)
6. `Gte`, `Gt`, `Lte`, `Lt` (AND group)
7. `In`, `NotIn` (AND group)
8. All `*Or` variants collected into one `(... OR ...)` block
9. `Between` (AND group)
10. `IsNull`, `IsNotNull` (AND group)
11. Raw `.Where()` clauses (combined with `WhereCombining`, default `AND`)
12. Cursor keyset clause (multi-column OR cascade, or single-column `>` / `<`)

### LIMIT / OFFSET

```
LIMIT $N OFFSET $M    -- offset pagination (default)
LIMIT $N              -- cursor pagination (NoOffset=true)
```

---

## API Reference

### Builder methods

| Method | Description |
|---|---|
| `NewBuilder()` | Create a new builder with global defaults |
| `.Table(name)` | Set the FROM table |
| `.Schema(name)` | Set the database schema |
| `.Model(struct)` | Set the model for tag resolution |
| `.Page(n)` | Set the page number (1-based) |
| `.Limit(n)` | Set items per page (capped at `MaxLimit`) |
| `.Select(cols...)` | SELECT specific columns |
| `.OrderBy(col, dir?)` | Add ORDER BY clause |
| `.OrderByDesc(col)` | Add ORDER BY col DESC |
| `.Join(clause)` | Add raw JOIN clause |
| `.LeftJoin(table, on)` | Add LEFT JOIN |
| `.InnerJoin(table, on)` | Add INNER JOIN |
| `.RightJoin(table, on)` | Add RIGHT JOIN |
| `.Search(term, fields...)` | Full-text search across fields |
| `.Eq(field, vals...)` | Equality (OR within field) |
| `.EqAnd(field, vals...)` | Equality AND |
| `.EqOr(field, vals...)` | Equality in OR group |
| `.In(field, vals...)` | IN |
| `.NotIn(field, vals...)` | NOT IN |
| `.InOr(field, vals...)` | IN in OR group |
| `.NotInOr(field, vals...)` | NOT IN in OR group |
| `.WhereLike(field, vals...)` | ILIKE (OR within field) |
| `.LikeAnd(field, vals...)` | ILIKE AND |
| `.LikeOr(field, vals...)` | ILIKE in OR group |
| `.WhereGreaterThan(field, val)` | > |
| `.WhereGreaterThanOrEqual(field, val)` | >= |
| `.WhereLessThan(field, val)` | < |
| `.WhereLessThanOrEqual(field, val)` | <= |
| `.GteOr(field, val)` | >= in OR group |
| `.GtOr(field, val)` | > in OR group |
| `.LteOr(field, val)` | <= in OR group |
| `.LtOr(field, val)` | < in OR group |
| `.WhereBetween(field, min, max)` | BETWEEN |
| `.WhereIsNull(field)` | IS NULL |
| `.WhereIsNotNull(field)` | IS NOT NULL |
| `.WhereIsNullOr(field)` | IS NULL in OR group |
| `.WhereIsNotNullOr(field)` | IS NOT NULL in OR group |
| `.Where(clause, args...)` | Raw WHERE clause |
| `.After(col, val)` | Single-column forward cursor |
| `.Before(col, val)` | Single-column backward cursor |
| `.WithoutOffset()` | Disable OFFSET |
| `.WithVacuum()` | Use count_estimate |
| `.FromJSON(json)` | Populate from JSON string |
| `.FromMap(map)` | Populate from map |
| `.FromStruct(struct)` | Populate from any struct (incl. PaginationParams) |
| `.BuildSQL()` | Return `(query, args, error)` — paginated SELECT only |
| `.BuildCountSQL()` | Return `(query, args, error)` — SELECT COUNT only |
| `.Build()` | Return `(*SQLResult, error)` — both queries at once |
| `.CurrentPage()` | Return current page number |
| `.CurrentLimit()` | Return current items per page |

### Response constructors

| Function | Description |
|---|---|
| `NewPage[T](data, total, params, url)` | Offset pagination response with HATEOAS |
| `NewCursorPage[T](rawItems, params, url)` | Cursor pagination response with HATEOAS |

### Binding functions

| Function | Description |
|---|---|
| `BindQueryParamsToStruct(url.Values)` | Bind URL values to `*PaginationParams` |
| `BindQueryStringToStruct(string)` | Bind raw query string to `*PaginationParams` |
| `BindQueryParams(url.Values, target)` | Bind URL values to any struct with `query` tags |
| `NewPaginationParams()` | Create `PaginationParams` with global defaults |

### Cursor functions

| Function | Description |
|---|---|
| `EncodeCursor(col, val, dir)` | Encode single-column cursor token |
| `DecodeCursor(token)` | Decode single-column cursor token |

### Configuration functions

| Function | Description |
|---|---|
| `SetDefaultLimit(n)` | Set default items per page |
| `SetMaxLimit(n)` | Set maximum allowed limit |
| `SetDebugMode(bool)` | Enable/disable SQL logging |
| `SetLogger(*slog.Logger)` | Set custom slog logger |
| `GetDefaultLimit()` | Get current default limit |
| `GetMaxLimit()` | Get current max limit |
| `IsDebugMode()` | Get debug mode status |

---

## Migration from v3

```bash
go get github.com/booscaaa/go-paginate/v4
```

| v3 | v4 |
|---|---|
| `paginate.Paginate(...)` | `paginate.NewBuilder().Table(...).Model(...).Build()` |
| Manual response struct | `paginate.NewPage[T](data, total, params, url)` |
| No cursor pagination | `paginate.NewCursorPage[T](rawItems, params, url)` |
| `import ".../v3/paginate"` | `import ".../v4/paginate"` |

All filter methods and query string keys are identical between v3 and v4.
