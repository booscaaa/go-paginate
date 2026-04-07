<p align="center">
  <img src="https://raw.githubusercontent.com/booscaaa/go-paginate/master/assets/icon.png" alt="Go Paginate Logo" width="200">
</p>

<p align="center">
  <h1 align="center">Go Paginate — The Ultimate Go Pagination Library</h1>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/booscaaa/go-paginate/v4"><img alt="Reference" src="https://img.shields.io/badge/go-reference-purple?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/booscaaa/go-paginate.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge"></a>
    <img alt="GitHub Workflow Status" src="https://img.shields.io/github/actions/workflow/status/booscaaa/go-paginate/test.yaml?style=for-the-badge">
    <a href="https://codecov.io/gh/booscaaa/go-paginate"><img alt="Coverage" src="https://img.shields.io/codecov/c/github/booscaaa/go-paginate/master.svg?style=for-the-badge"></a>
    <img alt="Go Version" src="https://img.shields.io/badge/go-1.21+-blue?style=for-the-badge">
  </p>
</p>

<br>

## Versions

| Version | Status | Documentation | Install |
|---|---|---|---|
| **v4** ✨ | **Latest — recommended** | [v4/README.md](v4/README.md) | `go get github.com/booscaaa/go-paginate/v4` |
| v3 | Stable | [v3/README.md](v3/README.md) | `go get github.com/booscaaa/go-paginate/v3` |
| v2 | Legacy | — | `go get github.com/booscaaa/go-paginate/v2` |

---

## End-to-End Example

A complete walkthrough of both pagination modes: from the HTTP request all the way to the JSON response.

### Model

```go
type User struct {
    ID        int       `json:"id"         paginate:"users.id"`
    Name      string    `json:"name"       paginate:"users.name"`
    Email     string    `json:"email"      paginate:"users.email"`
    Role      string    `json:"role"       paginate:"users.role"`
    Active    bool      `json:"active"     paginate:"users.active"`
    CreatedAt time.Time `json:"created_at" paginate:"users.created_at"`
}
```

---

### Offset Pagination

#### Request

```
GET /users?page=2&limit=3&sort=-created_at&eq[role]=admin&eq[role]=editor&like[name]=john
```

#### Handler

```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
    // 1. Bind all query params (page, limit, sort, filters) in one call
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 2. Build queries — FromStruct maps every param automatically
    result, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        FromStruct(params).
        Build()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 3. Execute against the database
    rows, _ := db.QueryContext(r.Context(), result.Query, result.Args...)
    users   := scanUsers(rows)

    var total int
    db.QueryRowContext(r.Context(), result.CountQuery, result.CountArgs...).Scan(&total)

    // 4. Build the response — page/limit derived from params automatically
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(paginate.NewPage(users, total, params, r.URL))
}
```

#### Generated SQL

```sql
-- result.Query
SELECT * FROM users
WHERE (users.role = $1 OR users.role = $2)
  AND (users.name::TEXT ILIKE $3)
ORDER BY users.created_at DESC
LIMIT $4 OFFSET $5
-- args: ["admin", "editor", "%john%", 3, 3]

-- result.CountQuery
SELECT COUNT(users.id) FROM users
WHERE (users.role = $1 OR users.role = $2)
  AND (users.name::TEXT ILIKE $3)
-- args: ["admin", "editor", "%john%"]
```

#### JSON Response

```json
{
  "data": [
    { "id": 7,  "name": "John Smith",  "email": "john.smith@example.com",  "role": "admin",  "active": true, "created_at": "2024-03-10T09:00:00Z" },
    { "id": 4,  "name": "Johnny Cash", "email": "johnny.cash@example.com", "role": "editor", "active": true, "created_at": "2024-02-28T14:30:00Z" },
    { "id": 2,  "name": "John Doe",    "email": "john.doe@example.com",    "role": "admin",  "active": true, "created_at": "2024-01-15T11:00:00Z" }
  ],
  "meta": {
    "current_page": 2,
    "per_page":     3,
    "total_items":  14,
    "total_pages":  5,
    "from":         4,
    "to":           6,
    "has_prev":     true,
    "has_next":     true
  },
  "links": {
    "self":  "/users?eq%5Brole%5D=admin&eq%5Brole%5D=editor&like%5Bname%5D=john&limit=3&page=2&sort=-created_at",
    "first": "/users?eq%5Brole%5D=admin&eq%5Brole%5D=editor&like%5Bname%5D=john&limit=3&page=1&sort=-created_at",
    "last":  "/users?eq%5Brole%5D=admin&eq%5Brole%5D=editor&like%5Bname%5D=john&limit=3&page=5&sort=-created_at",
    "prev":  "/users?eq%5Brole%5D=admin&eq%5Brole%5D=editor&like%5Bname%5D=john&limit=3&page=1&sort=-created_at",
    "next":  "/users?eq%5Brole%5D=admin&eq%5Brole%5D=editor&like%5Bname%5D=john&limit=3&page=3&sort=-created_at"
  }
}
```

> All active filters (`eq[role]`, `like[name]`, `sort`, `limit`) are preserved in every link automatically.

---

### Cursor Pagination

#### First Request

```
GET /users/feed?limit=3&sort=-created_at,id&eq[active]=true
```

#### Handler

```go
func ListUsersFeed(w http.ResponseWriter, r *http.Request) {
    // 1. Bind — cursor token is included automatically when present
    params, err := paginate.BindQueryParamsToStruct(r.URL.Query())
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 2. Build — FromStruct decodes the cursor token and injects the
    //    keyset WHERE clause automatically. No special handling needed.
    query, args, err := paginate.NewBuilder().
        Table("users").
        Model(&User{}).
        OrderBy("created_at", "DESC").
        OrderBy("id").              // tie-breaker — always add a unique column last
        FromStruct(params).
        BuildSQL()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 3. Fetch limit+1 rows — the extra row is used to detect hasNext
    rows, _   := db.QueryContext(r.Context(), query, args...)
    rawItems  := scanUsers(rows) // may contain limit+1 items

    // 4. NewCursorPage handles everything: trim, hasNext, hasPrev, token encoding
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(paginate.NewCursorPage(rawItems, params, r.URL))
}
```

#### Generated SQL — first page (no cursor)

```sql
SELECT * FROM users
WHERE (users.active = $1)
ORDER BY users.created_at DESC, users.id ASC
LIMIT $2
-- args: [true, 4]   (limit+1 = 3+1)
```

#### JSON Response — first page

```json
{
  "data": [
    { "id": 9, "name": "Alice",   "email": "alice@example.com",   "role": "user", "active": true, "created_at": "2024-04-01T10:00:00Z" },
    { "id": 7, "name": "Bob",     "email": "bob@example.com",     "role": "user", "active": true, "created_at": "2024-03-20T08:00:00Z" },
    { "id": 4, "name": "Charlie", "email": "charlie@example.com", "role": "user", "active": true, "created_at": "2024-02-15T16:00:00Z" }
  ],
  "meta": {
    "per_page": 3,
    "has_next": true,
    "has_prev": false
  },
  "links": {
    "self": "/users/feed?eq%5Bactive%5D=true&limit=3&sort=-created_at%2Cid",
    "next": "/users/feed?eq%5Bactive%5D=true&limit=3&sort=-created_at%2Cid&cursor=eyJjb2xzIjpbImNyZWF0ZWRfYXQiLCJpZCJdLCJ2YWxzIjpbIjIwMjQtMDItMTVUMTY6MDA6MDBaIiw0XSwiZGlycyI6WyJERVNDIiwiQVNDIl0sImRpciI6ImFmdGVyIn0=",
    "prev": null
  }
}
```

#### Second Request — following the `next` link

```
GET /users/feed?eq[active]=true&limit=3&sort=-created_at,id&cursor=eyJjb2xzIjpbImNyZWF0Z...
```

#### Generated SQL — second page (cursor decoded automatically)

```sql
SELECT * FROM users
WHERE (users.active = $1)
  AND (
    (users.created_at < $2)
    OR (users.created_at = $3 AND users.id > $4)
  )
ORDER BY users.created_at DESC, users.id ASC
LIMIT $5
-- args: [true, "2024-02-15T16:00:00Z", "2024-02-15T16:00:00Z", 4, 4]
```

> The keyset `WHERE` clause is built automatically from the cursor token — no `OFFSET` scan, 100% stable regardless of concurrent inserts or deletes.

#### JSON Response — second page

```json
{
  "data": [
    { "id": 3, "name": "Diana", "email": "diana@example.com", "role": "user", "active": true, "created_at": "2024-01-30T12:00:00Z" },
    { "id": 1, "name": "Eve",   "email": "eve@example.com",   "role": "user", "active": true, "created_at": "2024-01-10T09:00:00Z" }
  ],
  "meta": {
    "per_page": 3,
    "has_next": false,
    "has_prev": true
  },
  "links": {
    "self": "/users/feed?eq%5Bactive%5D=true&limit=3&sort=-created_at%2Cid&cursor=eyJjb2xzIjpb...",
    "next": null,
    "prev": "/users/feed?eq%5Bactive%5D=true&limit=3&sort=-created_at%2Cid&cursor=eyJjb2xzIjpbImNyZWF0ZWRfYXQiLCJpZCJdLCJ2YWxzIjpbIjIwMjQtMDEtMzBUMTI6MDA6MDBaIiwzXSwiZGlycyI6WyJERVNDIiwiQVNDIl0sImRpciI6ImJlZm9yZSJ9"
  }
}
```

---

## Why v4?

v4 is a complete rewrite with first-class **cursor pagination**, **generic response types**, and zero boilerplate.

**New in v4:**

- `Page[T]` and `CursorPage[T]` — generic response envelopes with HATEOAS links
- Cursor pagination with **keyset seek method** — 100% stable with any multi-column sort
- `NewCursorPage` derives everything from `PaginationParams` via reflection — zero boilerplate in handlers
- `Build()` returns both SELECT and COUNT queries at once
- All existing filters, sorts, joins, and query-string binding work identically

---

## Installation

```bash
go get github.com/booscaaa/go-paginate/v4
```

📖 **[Full v4 Documentation](v4/README.md)**

---

## Key Features

- **Fluent Builder API** — chainable, readable query construction
- **30+ filter types** — `Eq`, `Like`, `In`, `Between`, `IsNull`, `Gte`, `Gt`, and all `*Or` variants for granular OR grouping
- **Cursor pagination** — keyset seek method, stable with any sort configuration
- **Automatic HTTP binding** — `BindQueryParamsToStruct` converts URL query params in one call
- **HATEOAS links** — `self`, `first`, `last`, `prev`, `next` built automatically
- **JOIN support** — `LeftJoin`, `InnerJoin`, `RightJoin`
- **Multi-column sorting** — `?sort=-created_at,name,id`
- **Schema support** — `FROM schema.table`
- **Vacuum / count estimation** — fast row counts on large PostgreSQL tables
- **SQL injection safe** — parameterized queries, no string interpolation
- **Global config + env vars** — `GO_PAGINATE_DEFAULT_LIMIT`, `GO_PAGINATE_MAX_LIMIT`, `GO_PAGINATE_DEBUG`
- **Debug mode** — logs every generated SQL query via `slog`
