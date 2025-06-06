package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/booscaaa/go-paginate/v3/paginate"
)

func main() {
	fmt.Println("=== Example of Query Parameters Binding ===")

	// Example 1: Using BindQueryParamsToStruct with basic parameters
	fmt.Println("\n1. Bind de parâmetros básicos:")
	queryString1 := "page=2&limit=25&search=john&search_fields=name,email&vacuum=true"
	params1, err := paginate.BindQueryStringToStruct(queryString1)
	if err != nil {
		log.Fatalf("Error binding: %v", err)
	}

	fmt.Printf("Query String: %s\n", queryString1)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  Page: %d\n", params1.Page)
	fmt.Printf("  Limit: %d\n", params1.Limit)
	fmt.Printf("  Search: %s\n", params1.Search)
	fmt.Printf("  SearchFields: %v\n", params1.SearchFields)
	fmt.Printf("  Vacuum: %t\n", params1.Vacuum)

	// Example 2: Using complex parameters with arrays
	fmt.Println("\n2. Bind de parâmetros complexos:")
	queryString2 := "page=1&likeor[status]=active&likeor[status]=pending&eqor[age]=25&eqor[age]=30&gte[created_at]=2023-01-01&gt[score]=80"
	params2, err := paginate.BindQueryStringToStruct(queryString2)
	if err != nil {
		log.Fatalf("Error binding: %v", err)
	}

	fmt.Printf("Query String: %s\n", queryString2)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  LikeOr: %v\n", params2.LikeOr)
	fmt.Printf("  EqOr: %v\n", params2.EqOr)
	fmt.Printf("  Gte: %v\n", params2.Gte)
	fmt.Printf("  Gt: %v\n", params2.Gt)

	// Example 3: Using url.Values directly
	fmt.Println("\n3. Bind usando url.Values:")
	queryParams := url.Values{
		"page":            {"3"},
		"limit":           {"50"},
		"sort_columns":    {"name,created_at"},
		"sort_directions": {"ASC,DESC"},
		"likeand[name]":   {"admin"},
		"lte[updated_at]": {"2023-12-31"},
	}

	params3, err := paginate.BindQueryParamsToStruct(queryParams)
	if err != nil {
		log.Fatalf("Error binding: %v", err)
	}

	fmt.Printf("Query Params: %v\n", queryParams)
	fmt.Printf("Resultado:\n")
	fmt.Printf("  Page: %d\n", params3.Page)
	fmt.Printf("  Limit: %d\n", params3.Limit)
	fmt.Printf("  SortColumns: %v\n", params3.SortColumns)
	fmt.Printf("  SortDirections: %v\n", params3.SortDirections)
	fmt.Printf("  LikeAnd: %v\n", params3.LikeAnd)
	fmt.Printf("  Lte: %v\n", params3.Lte)

	// Example 4: Bind to custom struct
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
		log.Fatalf("Error binding custom: %v", err)
	}

	fmt.Printf("Custom Query Params: %v\n", customQueryParams)
	fmt.Printf("Resultado Customizado:\n")
	fmt.Printf("  Page: %d\n", customParams.Page)
	fmt.Printf("  Limit: %d\n", customParams.Limit)
	fmt.Printf("  Search: %s\n", customParams.Search)
	fmt.Printf("  Filters: %v\n", customParams.Filters)
	fmt.Printf("  Active: %t\n", customParams.Active)

	// Example 5: Simulating use in an HTTP handler
	fmt.Println("\n5. Example of usage in HTTP handler:")
	simulateHTTPHandler()
}

// simulateHTTPHandler simulates how to use bind in a real HTTP handler
func simulateHTTPHandler() {
	// Simulate an HTTP request URL with new sort pattern
	requestURL := "https://api.example.com/users?page=2&limit=20&search=john&search_fields=name,email&likeor[status]=active&likeor[status]=pending&gte[age]=18&sort=name&sort=-created_at"

	// Parse the URL
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	// Extract query parameters
	queryParams := parsedURL.Query()

	// Bind to pagination struct
	paginationParams, err := paginate.BindQueryParamsToStruct(queryParams)
	if err != nil {
		log.Fatalf("Error binding parameters: %v", err)
	}

	fmt.Printf("URL simulada: %s\n", requestURL)
	fmt.Printf("Parâmetros extraídos:\n")
	fmt.Printf("  Page: %d\n", paginationParams.Page)
	fmt.Printf("  Limit: %d\n", paginationParams.Limit)
	fmt.Printf("  Search: %s\n", paginationParams.Search)
	fmt.Printf("  SearchFields: %v\n", paginationParams.SearchFields)
	fmt.Printf("  LikeOr: %v\n", paginationParams.LikeOr)
	fmt.Printf("  Sort: %v\n", paginationParams.Sort)
	fmt.Printf("  SortColumns: %v\n", paginationParams.SortColumns)
	fmt.Printf("  SortDirections: %v\n", paginationParams.SortDirections)
	fmt.Printf("  Gte: %v\n", paginationParams.Gte)

	// Now you can use these parameters to build your database query
	fmt.Println("\n✅ Parameters ready for use in query construction!")

	// Additional example: Demonstrate how to use with FromStruct in builder
	fmt.Println("\n6. Example using FromStruct with new sort pattern:")
	demonstrateFromStructWithSort(paginationParams)
}

// demonstrateFromStructWithSort demonstrates how to use FromStruct with the new sort pattern
func demonstrateFromStructWithSort(params *paginate.PaginationParams) {
	// Define an example struct for the model
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
		log.Fatalf("Error generating SQL: %v", err)
	}

	fmt.Printf("SQL gerado: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
	fmt.Println("\n✅ Sort funcionando corretamente com FromStruct!")
}
