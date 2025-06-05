package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/booscaaa/go-paginate/v3/paginate"
)

func main() {
	fmt.Println("=== Exemplo de Bind de Query Parameters ===")

	// Exemplo 1: Usando BindQueryParamsToStruct com parâmetros básicos
	fmt.Println("\n1. Bind de parâmetros básicos:")
	queryString1 := "page=2&limit=25&search=john&search_fields=name,email&vacuum=true"
	params1, err := paginate.BindQueryStringToStruct(queryString1)
	if err != nil {
		log.Fatalf("Erro ao fazer bind: %v", err)
	}

	fmt.Printf("Query String: %s\n", queryString1)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  Page: %d\n", params1.Page)
	fmt.Printf("  Limit: %d\n", params1.Limit)
	fmt.Printf("  Search: %s\n", params1.Search)
	fmt.Printf("  SearchFields: %v\n", params1.SearchFields)
	fmt.Printf("  Vacuum: %t\n", params1.Vacuum)

	// Exemplo 2: Usando parâmetros complexos com arrays
	fmt.Println("\n2. Bind de parâmetros complexos:")
	queryString2 := "page=1&search_or[status]=active&search_or[status]=pending&equals_or[age]=25&equals_or[age]=30&gte[created_at]=2023-01-01&gt[score]=80"
	params2, err := paginate.BindQueryStringToStruct(queryString2)
	if err != nil {
		log.Fatalf("Erro ao fazer bind: %v", err)
	}

	fmt.Printf("Query String: %s\n", queryString2)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  SearchOr: %v\n", params2.SearchOr)
	fmt.Printf("  EqualsOr: %v\n", params2.EqualsOr)
	fmt.Printf("  Gte: %v\n", params2.Gte)
	fmt.Printf("  Gt: %v\n", params2.Gt)

	// Exemplo 3: Usando url.Values diretamente
	fmt.Println("\n3. Bind usando url.Values:")
	queryParams := url.Values{
		"page":             {"3"},
		"limit":            {"50"},
		"sort_columns":     {"name,created_at"},
		"sort_directions":  {"ASC,DESC"},
		"search_and[name]": {"admin"},
		"lte[updated_at]":  {"2023-12-31"},
	}

	params3, err := paginate.BindQueryParamsToStruct(queryParams)
	if err != nil {
		log.Fatalf("Erro ao fazer bind: %v", err)
	}

	fmt.Printf("Query Params: %v\n", queryParams)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  Page: %d\n", params3.Page)
	fmt.Printf("  Limit: %d\n", params3.Limit)
	fmt.Printf("  SortColumns: %v\n", params3.SortColumns)
	fmt.Printf("  SortDirections: %v\n", params3.SortDirections)
	fmt.Printf("  SearchAnd: %v\n", params3.SearchAnd)
	fmt.Printf("  Lte: %v\n", params3.Lte)

	// Exemplo 4: Bind para struct customizada
	fmt.Println("\n4. Bind para struct customizada:")
	type CustomPaginationParams struct {
		Page    int      `query:"page"`
		Limit   int      `query:"limit"`
		Search  string   `query:"search"`
		Filters []string `query:"filters"`
		Active  bool     `query:"active"`
	}

	customQueryParams := url.Values{
		"page":    {"4"},
		"limit":   {"100"},
		"search":  {"custom search"},
		"filters": {"filter1,filter2,filter3"},
		"active":  {"true"},
	}

	customParams := &CustomPaginationParams{}
	err = paginate.BindQueryParams(customQueryParams, customParams)
	if err != nil {
		log.Fatalf("Erro ao fazer bind customizado: %v", err)
	}

	fmt.Printf("Custom Query Params: %v\n", customQueryParams)
	fmt.Printf("Resultado Customizado:\n")
	fmt.Printf("  Page: %d\n", customParams.Page)
	fmt.Printf("  Limit: %d\n", customParams.Limit)
	fmt.Printf("  Search: %s\n", customParams.Search)
	fmt.Printf("  Filters: %v\n", customParams.Filters)
	fmt.Printf("  Active: %t\n", customParams.Active)

	// Exemplo 5: Simulando uso em um handler HTTP
	fmt.Println("\n5. Exemplo de uso em handler HTTP:")
	simulateHTTPHandler()
}

// simulateHTTPHandler simula como usar o bind em um handler HTTP real
func simulateHTTPHandler() {
	// Simular uma URL de request HTTP com novo padrão de sort
	requestURL := "https://api.example.com/users?page=2&limit=20&search=john&search_fields=name,email&search_or[status]=active&search_or[status]=pending&gte[age]=18&sort=name&sort=-created_at"

	// Parse da URL
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		log.Fatalf("Erro ao fazer parse da URL: %v", err)
	}

	// Extrair query parameters
	queryParams := parsedURL.Query()

	// Fazer bind para struct de paginação
	paginationParams, err := paginate.BindQueryParamsToStruct(queryParams)
	if err != nil {
		log.Fatalf("Erro ao fazer bind dos parâmetros: %v", err)
	}

	fmt.Printf("URL simulada: %s\n", requestURL)
	fmt.Printf("Parâmetros extraídos:\n")
	fmt.Printf("  Page: %d\n", paginationParams.Page)
	fmt.Printf("  Limit: %d\n", paginationParams.Limit)
	fmt.Printf("  Search: %s\n", paginationParams.Search)
	fmt.Printf("  SearchFields: %v\n", paginationParams.SearchFields)
	fmt.Printf("  SearchOr: %v\n", paginationParams.SearchOr)
	fmt.Printf("  Sort: %v\n", paginationParams.Sort)
	fmt.Printf("  SortColumns: %v\n", paginationParams.SortColumns)
	fmt.Printf("  SortDirections: %v\n", paginationParams.SortDirections)
	fmt.Printf("  Gte: %v\n", paginationParams.Gte)

	// Agora você pode usar esses parâmetros para construir sua query de banco de dados
	fmt.Println("\n✅ Parâmetros prontos para uso na construção da query!")

	// Exemplo adicional: Demonstrar como usar com FromStruct no builder
	fmt.Println("\n6. Exemplo usando FromStruct com novo padrão de sort:")
	demonstrateFromStructWithSort(paginationParams)
}

// demonstrateFromStructWithSort demonstra como usar FromStruct com o novo padrão de sort
func demonstrateFromStructWithSort(params *paginate.PaginationParams) {
	// Definir uma struct de exemplo para o modelo
	type User struct {
		ID        int    `json:"id" paginate:"id"`
		Name      string `json:"name" paginate:"name"`
		Email     string `json:"email" paginate:"email"`
		Status    string `json:"status" paginate:"status"`
		Age       int    `json:"age" paginate:"age"`
		CreatedAt string `json:"created_at" paginate:"created_at"`
	}

	// Criar builder e usar FromStruct
	builder := paginate.NewBuilder().
		Table("users").
		Model(User{}).
		FromStruct(params)

	// Gerar SQL
	sql, args, err := builder.BuildSQL()
	if err != nil {
		log.Fatalf("Erro ao gerar SQL: %v", err)
	}

	fmt.Printf("SQL gerado: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
	fmt.Println("\n✅ Sort funcionando corretamente com FromStruct!")
}
