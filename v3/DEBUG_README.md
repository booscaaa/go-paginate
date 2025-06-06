# Go-Paginate v3 - Debug Mode

This document describes the debug functionality implemented in go-paginate v3, which allows structured logging of all generated SQL queries.

## 🔧 Configuration

### Environment Variables

```bash
# Enable debug mode (prints generated SQL)
export GO_PAGINATE_DEBUG=true

# Set default page limit
export GO_PAGINATE_DEFAULT_LIMIT=25

# Set maximum page limit
export GO_PAGINATE_MAX_LIMIT=1000
```

### Global Configuration

```go
package main

import "github.com/booscaaa/go-paginate/v3/paginate"

func init() {
    // Set global configurations
    paginate.SetDefaultLimit(25)
    paginate.SetMaxLimit(1000)
    paginate.SetDebugMode(true)
}
```

## 📊 Structured Logs

When debug mode is enabled (`GO_PAGINATE_DEBUG=true` or `paginate.SetDebugMode(true)`), go-paginate will generate structured logs in JSON format for all created SQL queries.

### Log Format

```json
{
  "time": "2025-06-06T09:03:44.087649546-03:00",
  "level": "INFO",
  "msg": "Generated SQL query",
  "component": "go-paginate-sql",
  "operation": "BuildSQL",
  "query": "SELECT * FROM users WHERE name ILIKE $1 ORDER BY name ASC LIMIT $2 OFFSET $3",
  "args": ["john", 10, 0],
  "args_count": 3
}
```

### Log Fields

- **time**: Log timestamp
- **level**: Log level (INFO for SQL queries)
- **msg**: Descriptive message
- **component**: Component that generated the log (`go-paginate-sql`)
- **operation**: Operation that generated the query:
  - `BuildSQL`: Main pagination query
  - `BuildCountSQL`: Count query
  - `GenerateSQL`: Internally generated query
  - `GenerateCountQuery`: Internally generated count query
  - `GenerateCountQuery (Vacuum)`: Optimized count query
- **query**: The generated SQL query
- **args**: Array with query arguments
- **args_count**: Total number of arguments

## 🚀 Usage Example

```go
package main

import (
    "log/slog"
    "os"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

type User struct {
    ID    int    `json:"id" paginate:"id"`
    Name  string `json:"name" paginate:"name"`
    Email string `json:"email" paginate:"email"`
}

func main() {
    // Configure structured logging
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    slog.SetDefault(logger)
    
    // Enable debug mode
    paginate.SetDebugMode(true)
    
    // Build query
    sql, args, err := paginate.NewBuilder().
        Table("users").
        Model(User{}).
        Page(1).
        Limit(10).
        Search("john", "name", "email").
        OrderBy("name", "ASC").
        BuildSQL()
    
    if err != nil {
        panic(err)
    }
    
    // Logs will be automatically printed in JSON format
    // The query and arguments are also available for use
    println("SQL:", sql)
    println("Args:", args)
}
```

## 🔍 Operations that Generate Logs

### 1. BuildSQL()
Generates logs for the main pagination query:
```json
{
  "operation": "BuildSQL",
  "query": "SELECT * FROM users WHERE name ILIKE $1 LIMIT $2 OFFSET $3",
  "args": ["%john%", 10, 0]
}
```

### 2. BuildCountSQL()
Generates logs for the count query:
```json
{
  "operation": "BuildCountSQL",
  "query": "SELECT COUNT(id) FROM users WHERE name ILIKE $1",
  "args": ["%john%"]
}
```

### 3. GenerateSQL() (interno)
Called internally by BuildSQL():
```json
{
  "operation": "GenerateSQL",
  "query": "SELECT * FROM users WHERE name ILIKE $1 LIMIT $2 OFFSET $3",
  "args": ["%john%", 10, 0]
}
```

### 4. GenerateCountQuery() (interno)
Called internally by BuildCountSQL():
```json
{
  "operation": "GenerateCountQuery",
  "query": "SELECT COUNT(id) FROM users WHERE name ILIKE $1",
  "args": ["%john%"]
}
```

### 5. Vacuum Mode
When vacuum mode is enabled:
```json
{
  "operation": "GenerateCountQuery (Vacuum)",
  "query": "SELECT count_estimate('SELECT COUNT(1) FROM users WHERE name ILIKE ''$1''');",
  "args": ["%john%"]
}
```

## ⚙️ Advanced Configuration

### Custom Logger

```go
// Configure custom logger
customLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))

paginate.SetLogger(customLogger)
```

### Check Configuration Status

```go
// Check current configurations
fmt.Println("Debug Mode:", paginate.IsDebugMode())
fmt.Println("Default Limit:", paginate.GetDefaultLimit())
fmt.Println("Max Limit:", paginate.GetMaxLimit())
```

## 🛡️ Security

- Logs include query arguments, but these are parameterized and safe against SQL injection
- In production, consider disabling debug mode or configuring the appropriate log level
- Logs may contain sensitive data in arguments - configure appropriately in production environments

## 📝 Notes

- Debug mode uses the `INFO` level to ensure log visibility
- Each operation may generate multiple logs (internal + public)
- Logs are thread-safe and use Go's standard logger (`log/slog`)
- Configuration is global and affects all paginate instances