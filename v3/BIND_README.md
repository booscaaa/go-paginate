# Query Parameters Bind

This functionality allows binding URL query parameters to pagination structs in a simple and efficient way.

## Funcionalidades

- ✅ Bind de parâmetros básicos (page, limit, search, etc.)
- ✅ Suporte a arrays e slices
- ✅ Parâmetros complexos com sintaxe de array (`likeor[field]`, `eqor[field]`, etc.)
- ✅ Conversão automática de tipos (int, bool, string)
- ✅ Suporte a structs customizadas
- ✅ Validação de tipos
- ✅ Valores padrão

## Uso Básico

### 1. Bind para PaginationParams (struct padrão)

```go
package main

import (
    "fmt"
    "log"
    "github.com/booscaaa/go-paginate/v3/paginate"
)

func main() {
    // A partir de uma query string
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

### 2. Bind usando url.Values

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

## Parâmetros Suportados

### Parâmetros Básicos

| Parameter         | Type     | Description                 | Example                        |
| ----------------- | -------- | --------------------------- | ------------------------------ |
| `page`            | int      | Número da página            | `page=2`                       |
| `limit`           | int      | Itens por página            | `limit=25`                     |
| `items_per_page`  | int      | Alias para limit            | `items_per_page=25`            |
| `search`          | string   | Search term                 | `search=john`                  |
| `search_fields`   | []string | Fields for search           | `search_fields=name,email`     |
| `sort_columns`    | []string | Columns for sorting         | `sort_columns=name,created_at` |
| `sort_directions` | []string | Direções de ordenação       | `sort_directions=ASC,DESC`     |
| `columns`         | []string | Colunas para seleção        | `columns=id,name,email`        |
| `vacuum`          | bool     | Usar estimativa de contagem | `vacuum=true`                  |
| `no_offset`       | bool     | Desabilitar OFFSET          | `no_offset=false`              |

### Parâmetros Complexos (Sintaxe de Array)

| Parâmetro           | Tipo                | Descrição           | Exemplo                                              |
| ------------------- | ------------------- | ------------------- | ---------------------------------------------------- |
| `likeor[field]`     | map[string][]string | Busca OR por campo  | `likeor[status]=active&likeor[status]=pending`       |
| `likeand[field]`    | map[string][]string | Busca AND por campo | `likeand[name]=john`                                 |
| `eqor[field]`       | map[string][]any    | Igualdade OR        | `eqor[age]=25&eqor[age]=30`                          |
| `eqand[field]`      | map[string][]any    | Igualdade AND       | `eqand[role]=admin`                                  |
| `gte[field]`        | map[string]any      | Maior ou igual      | `gte[age]=18`                                        |
| `gt[field]`         | map[string]any      | Maior que           | `gt[score]=80`                                       |
| `lte[field]`        | map[string]any      | Menor ou igual      | `lte[price]=100.50`                                  |
| `lt[field]`         | map[string]any      | Menor que           | `lt[date]=2023-12-31`                                |

## Exemplos Avançados

### 1. Parâmetros Complexos

```go
queryString := "page=1&likeor[status]=active&likeor[status]=pending&eqor[age]=25&eqor[age]=30&gte[created_at]=2023-01-01"
params, err := paginate.BindQueryStringToStruct(queryString)

// Resultado:
// params.LikeOr["status"] = ["active", "pending"]
// params.EqOr["age"] = [25, 30]
// params.Gte["created_at"] = "2023-01-01"
```

### 2. Struct Customizada

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

### 3. Uso em Handler HTTP

```go
func usersHandler(w http.ResponseWriter, r *http.Request) {
    // Extrair query parameters da request
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

## Conversão de Tipos

A biblioteca faz conversão automática de tipos:

- **Strings**: Usadas diretamente
- **Integers**: Convertidos com `strconv.Atoi()`
- **Booleans**: Convertidos com `strconv.ParseBool()`
- **Floats**: Convertidos com `strconv.ParseFloat()`
- **Slices**: Múltiplos valores ou valores separados por vírgula

## Tratamento de Erros

- Valores inválidos são ignorados (mantém valor padrão)
- Tipos incompatíveis são ignorados
- Erros de parsing da query string são retornados
- Targets inválidos (não-ponteiro ou não-struct) retornam erro

## Valores Padrão

A struct `PaginationParams` tem valores padrão:

```go
params := &PaginationParams{
    Page:         1,  // default page
    Limit:        10, // default limit
    ItemsPerPage: 10, // default items per page
}
```

## Executar Exemplo

Para ver a funcionalidade em ação:

```bash
go run example_bind.go
```

## Executar Testes

```bash
go test -v ./paginate -run TestBind
```
