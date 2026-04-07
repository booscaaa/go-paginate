# Release Notes — go-paginate v4.0.0

> **Released:** 2026-04-07  
> **Module:** `github.com/booscaaa/go-paginate/v4`  
> **Go requirement:** 1.21+

---

## Overview

v4 is a full rewrite of go-paginate focused on three goals:

1. **Cursor pagination with keyset seek** — stable, index-friendly pagination for high-volume tables.
2. **Generic response envelopes** — type-safe `Page[T]` and `CursorPage[T]` with HATEOAS links baked in.
3. **Global configuration** — default and max limits driven by environment variables or code, no more per-call boilerplate.

---

## What's New

### Cursor Pagination (Keyset Seek)

New `CursorPage[T]` type and `NewCursorPage` constructor provide offset-free, stable cursor pagination.

- Multi-column keyset WHERE clause generation — handles any number of sort columns with mixed ASC/DESC directions.
- Opaque base64 cursor tokens encode column names, values, sort directions, and direction (`after` / `before`). Tokens are safe to expose in public URLs.
- `EncodeCursor` / `DecodeCursor` for single-column usage (backward compatible).
- `NewCursorPage[T]` automatically:
  - Detects `hasNext` via the `limit+1` fetch pattern (no extra COUNT query).
  - Extracts cursor values from the last/first item via JSON tag reflection.
  - Builds `next` and `prev` HATEOAS links preserving all existing query params.

```go
// repository: fetch limit+1 rows
b := paginate.NewBuilder().
    Table("events").Model(&Event{}).
    OrderBy("created_at", "DESC").OrderBy("id").
    FromStruct(params)
rows, _ := repo.Fetch(b.BuildSQL())

// handler: one line
page := paginate.NewCursorPage(rows, params, r.URL)
// page.Links.Next / page.Links.Prev are ready-to-use URLs
```

`PaginationParams` now includes a `Cursor string` field that `BindQueryParamsToStruct` populates automatically from `?cursor=<token>`.

---

### Generic Page Envelope — `Page[T]`

`NewPage[T]` replaces manual pagination math with a single call:

```go
page := paginate.NewPage(items, totalCount, params, r.URL)
```

Response shape:

```json
{
  "data": [...],
  "meta": {
    "current_page": 2,
    "per_page": 20,
    "total_items": 340,
    "total_pages": 17,
    "from": 21,
    "to": 40,
    "has_prev": true,
    "has_next": true
  },
  "links": {
    "self":  "https://api.example.com/users?page=2",
    "first": "https://api.example.com/users?page=1",
    "last":  "https://api.example.com/users?page=17",
    "prev":  "https://api.example.com/users?page=1",
    "next":  "https://api.example.com/users?page=3"
  }
}
```

`prev` and `next` are `null` (not omitted) at the boundaries.

---

### Global Configuration

A singleton `GlobalConfig` is initialized at startup, readable via environment variables:

| Environment variable          | Default | Description                          |
|-------------------------------|---------|--------------------------------------|
| `GO_PAGINATE_DEFAULT_LIMIT`   | `10`    | Default items per page               |
| `GO_PAGINATE_MAX_LIMIT`       | `100`   | Hard cap — requests above this are clamped |
| `GO_PAGINATE_DEBUG`           | `false` | Log every generated SQL via `slog`   |

Programmatic overrides:

```go
paginate.SetDefaultLimit(25)
paginate.SetMaxLimit(500)
paginate.SetDebugMode(true)
paginate.SetLogger(myLogger) // custom slog.Logger
```

All getters (`GetDefaultLimit`, `GetMaxLimit`, `IsDebugMode`) are safe to call anywhere.

---

### Fluent Builder — `NewBuilder()`

Full rewrite of the builder with consistent method naming, OR-group support, and a `FromStruct` / `FromJSON` / `FromMap` bootstrap path.

**New / renamed methods:**

| v4 Method | Notes |
|---|---|
| `Eq(field, values...)` | OR-grouped equality (field = v1 OR field = v2) |
| `EqAnd(field, values...)` | AND-grouped equality |
| `EqOr(field, values...)` | Merges into global OR clause |
| `In(field, values...)` | `field IN (...)` |
| `NotIn(field, values...)` | `field NOT IN (...)` |
| `InOr / NotInOr` | IN/NOT IN in the global OR group |
| `WhereLike(field, values...)` | ILIKE OR within field |
| `LikeOr / LikeAnd` | Cross-field ILIKE in OR / AND groups |
| `GteOr / GtOr / LteOr / LtOr` | Range comparisons in the global OR group |
| `WhereIsNull / WhereIsNullOr` | NULL checks (AND and OR groups) |
| `WhereIsNotNull / WhereIsNotNullOr` | NOT NULL checks |
| `WhereBetween(field, min, max)` | `BETWEEN` clause |
| `After(column, value)` | Forward cursor (single-column) |
| `Before(column, value)` | Backward cursor (single-column) |
| `WithoutOffset()` | Disables OFFSET (for cursor mode) |
| `WithVacuum()` | Uses PostgreSQL `count_estimate` |
| `LeftJoin / InnerJoin / RightJoin` | Typed JOIN helpers |
| `FromStruct(s)` | Bootstraps builder from `PaginationParams` or any struct |
| `FromJSON(str)` | Bootstraps builder from a JSON string |
| `FromMap(m)` | Bootstraps builder from `map[string]any` |
| `Build()` | Returns `(*QueryParams, error)` |
| `BuildSQL()` | Returns `(query string, args []any, err error)` |
| `SQLResult` | Struct returned by `BuildSQL` for destructuring |

---

### Query Param Binding

`BindQueryParamsToStruct(url.Values)` and `BindQueryStringToStruct(string)` parse the full `PaginationParams` from an HTTP query string, including all filter maps using bracket syntax:

```
GET /users?page=2&limit=20&eq[status]=active&gte[age]=18&sort=-created_at,id
```

`NewPaginationParams()` returns a zero-value struct with all maps pre-initialized (no nil-map panics).

---

### Multi-Column Keyset WHERE Generation

`buildCursorWhereMulti` produces the correct seek-method predicate for any number of columns. For `ORDER BY created_at DESC, id ASC` with an "after" cursor it generates:

```sql
((created_at < $1) OR (created_at = $2 AND id > $3))
```

This is index-friendly and handles ties correctly.

---

## Breaking Changes from v3

| Area | v3 | v4 |
|---|---|---|
| Module path | `github.com/booscaaa/go-paginate/v3` | `github.com/booscaaa/go-paginate/v4` |
| Response envelope | manual struct | generic `Page[T]` / `CursorPage[T]` |
| Cursor pagination | not available | `NewCursorPage[T]` + token system |
| Global limits | not available | `SetDefaultLimit` / `SetMaxLimit` / env vars |
| `SearchOr` | method name | renamed to `LikeOr` (`SearchOr` kept as alias) |
| `SearchAnd` | method name | renamed to `LikeAnd` (`SearchAnd` kept as alias) |
| `WhereEqualsOr` | method name | renamed to `EqOr` (alias kept) |
| Build output | `(string, []any)` | `BuildSQL()` returns `(string, []any, error)` |

---

## Migration Guide

```go
// v3
import "github.com/booscaaa/go-paginate/v3/paginate"

// v4
import "github.com/booscaaa/go-paginate/v4/paginate"
```

Replace `Build()` calls:

```go
// v3
query, args := builder.Build()

// v4
query, args, err := builder.BuildSQL()
```

Wrap results:

```go
// v4 — offset pagination
page := paginate.NewPage(items, total, params, r.URL)

// v4 — cursor pagination
page := paginate.NewCursorPage(rawItems, params, r.URL) // fetch limit+1
```

---

## Installation

```bash
go get github.com/booscaaa/go-paginate/v4@v4.0.0
```
