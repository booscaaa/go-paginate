# Go-Paginate v3 - Debug Mode

Este documento descreve a funcionalidade de debug implementada no go-paginate v3, que permite logging estruturado de todas as queries SQL geradas.

## 🔧 Configuração

### Variáveis de Ambiente

```bash
# Habilitar modo debug (imprime SQL gerado)
export GO_PAGINATE_DEBUG=true

# Set default page limit
export GO_PAGINATE_DEFAULT_LIMIT=25

# Definir limite máximo de página
export GO_PAGINATE_MAX_LIMIT=1000
```

### Configuração Global

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

## 📊 Logs Estruturados

Quando o modo debug está habilitado (`GO_PAGINATE_DEBUG=true` ou `paginate.SetDebugMode(true)`), o go-paginate irá gerar logs estruturados em formato JSON para todas as queries SQL criadas.

### Formato dos Logs

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

### Campos dos Logs

- **time**: Timestamp do log
- **level**: Nível do log (INFO para queries SQL)
- **msg**: Mensagem descritiva
- **component**: Componente que gerou o log (`go-paginate-sql`)
- **operation**: Operação que gerou a query:
  - `BuildSQL`: Query principal de paginação
  - `BuildCountSQL`: Query de contagem
  - `GenerateSQL`: Query gerada internamente
  - `GenerateCountQuery`: Query de contagem gerada internamente
  - `GenerateCountQuery (Vacuum)`: Query de contagem otimizada
- **query**: A query SQL gerada
- **args**: Array com os argumentos da query
- **args_count**: Número total de argumentos

## 🚀 Exemplo de Uso

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
    // Configurar logging estruturado
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    slog.SetDefault(logger)
    
    // Habilitar modo debug
    paginate.SetDebugMode(true)
    
    // Construir query
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

## 🔍 Operações que Geram Logs

### 1. BuildSQL()
Gera logs para a query principal de paginação:
```json
{
  "operation": "BuildSQL",
  "query": "SELECT * FROM users WHERE name ILIKE $1 LIMIT $2 OFFSET $3",
  "args": ["%john%", 10, 0]
}
```

### 2. BuildCountSQL()
Gera logs para a query de contagem:
```json
{
  "operation": "BuildCountSQL",
  "query": "SELECT COUNT(id) FROM users WHERE name ILIKE $1",
  "args": ["%john%"]
}
```

### 3. GenerateSQL() (interno)
Chamado internamente pelo BuildSQL():
```json
{
  "operation": "GenerateSQL",
  "query": "SELECT * FROM users WHERE name ILIKE $1 LIMIT $2 OFFSET $3",
  "args": ["%john%", 10, 0]
}
```

### 4. GenerateCountQuery() (interno)
Chamado internamente pelo BuildCountSQL():
```json
{
  "operation": "GenerateCountQuery",
  "query": "SELECT COUNT(id) FROM users WHERE name ILIKE $1",
  "args": ["%john%"]
}
```

### 5. Vacuum Mode
Quando o modo vacuum está habilitado:
```json
{
  "operation": "GenerateCountQuery (Vacuum)",
  "query": "SELECT count_estimate('SELECT COUNT(1) FROM users WHERE name ILIKE ''$1''');",
  "args": ["%john%"]
}
```

## ⚙️ Configuração Avançada

### Logger Customizado

```go
// Configurar logger customizado
customLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))

paginate.SetLogger(customLogger)
```

### Verificar Status da Configuração

```go
// Check current configurations
fmt.Println("Debug Mode:", paginate.IsDebugMode())
fmt.Println("Default Limit:", paginate.GetDefaultLimit())
fmt.Println("Max Limit:", paginate.GetMaxLimit())
```

## 🛡️ Segurança

- Os logs incluem os argumentos da query, mas estes são parametrizados e seguros contra SQL injection
- Em produção, considere desabilitar o modo debug ou configurar o nível de log apropriado
- Os logs podem conter dados sensíveis nos argumentos - configure adequadamente em ambientes de produção

## 📝 Notas

- O modo debug utiliza o nível `INFO` para garantir visibilidade dos logs
- Cada operação pode gerar múltiplos logs (interno + público)
- Os logs são thread-safe e utilizam o logger padrão do Go (`log/slog`)
- A configuração é global e afeta todas as instâncias do paginate