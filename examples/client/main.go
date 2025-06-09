package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/booscaaa/go-paginate/v3/client"
)

func main() {
	fmt.Println("=== Go Paginate v3 Client Examples ===")
	fmt.Println()

	// Example 1: Basic pagination
	fmt.Println("1. 📄 Basic Pagination:")
	basicPaginationExample()
	fmt.Println()

	// Example 2: Search and filtering
	fmt.Println("2. 🔍 Search and Filtering:")
	searchAndFilterExample()
	fmt.Println()

	// Example 3: Complex filtering
	fmt.Println("3. 🎯 Complex Filtering:")
	complexFilteringExample()
	fmt.Println()

	// Example 4: Building from existing URL
	fmt.Println("4. 🔗 Building from Existing URL:")
	buildFromURLExample()
	fmt.Println()

	// Example 5: Client cloning and reuse
	fmt.Println("5. 📋 Client Cloning and Reuse:")
	clientCloningExample()
	fmt.Println()

	// Example 6: HTTP client integration
	fmt.Println("6. 🌐 HTTP Client Integration:")
	httpClientExample()
	fmt.Println()

	// Example 7: Query string only
	fmt.Println("7. 📝 Query String Only:")
	queryStringOnlyExample()
}

func basicPaginationExample() {
	// Create a new client for a users API endpoint
	client := client.New("https://api.example.com/users")
	
	// Build URL with basic pagination
	url := client.
		Page(2).
		Limit(25).
		BuildURL()
	
	fmt.Printf("  URL: %s\n", url)
	fmt.Printf("  Query String: %s\n", client.BuildQueryString())
}

func searchAndFilterExample() {
	// Create client and add search parameters
	client := client.New("https://api.example.com/users")
	
	url := client.
		Page(1).
		Limit(10).
		Search("john").
		SearchFields("name", "email", "username").
		Sort("name", "-created_at").
		Columns("id", "name", "email", "created_at").
		BuildURL()
	
	fmt.Printf("  URL: %s\n", url)
}

func complexFilteringExample() {
	// Create client with complex filters
	client := client.New("https://api.example.com/users")
	
	url := client.
		Page(1).
		Limit(20).
		// LIKE filters
		LikeOr("status", "active", "pending").
		LikeAnd("name", "john", "doe").
		// Equality filters
		EqOr("age", 25, 30, 35).
		Eq("department", "IT").
		// Comparison filters
		Gte("created_at", "2023-01-01").
		Lt("score", 100).
		// IN filters
		In("role", "admin", "manager", "user").
		NotIn("status", "deleted", "banned").
		// BETWEEN filter
		Between("salary", 50000, 150000).
		// NULL filters
		IsNotNull("email").
		IsNull("deleted_at").
		// Special options
		Vacuum(true).
		BuildURL()
	
	fmt.Printf("  URL: %s\n", url)
	fmt.Printf("  Query String Length: %d characters\n", len(client.BuildQueryString()))
}

func buildFromURLExample() {
	// Start with an existing URL that has some parameters
	existingURL := "https://api.example.com/users?page=1&limit=10&search=existing"
	
	client, err := client.NewFromURL(existingURL)
	if err != nil {
		log.Printf("Error creating client from URL: %v", err)
		return
	}
	
	// Add more parameters to the existing ones
	newURL := client.
		Page(3). // This will override the existing page=1
		Sort("-created_at").
		Eq("status", "active").
		BuildURL()
	
	fmt.Printf("  Original URL: %s\n", existingURL)
	fmt.Printf("  Modified URL: %s\n", newURL)
}

func clientCloningExample() {
	// Create a base client with common parameters
	baseClient := client.New("https://api.example.com/users")
	baseClient.
		Limit(25).
		Columns("id", "name", "email", "created_at").
		Sort("name")
	
	// Clone for active users
	activeUsersClient := baseClient.Clone()
	activeUsersURL := activeUsersClient.
		Page(1).
		Eq("status", "active").
		BuildURL()
	
	// Clone for inactive users
	inactiveUsersClient := baseClient.Clone()
	inactiveUsersURL := inactiveUsersClient.
		Page(1).
		Eq("status", "inactive").
		BuildURL()
	
	// Clone for admin search
	adminSearchClient := baseClient.Clone()
	adminSearchURL := adminSearchClient.
		Page(1).
		Search("admin").
		SearchFields("name", "email").
		Eq("role", "admin").
		BuildURL()
	
	fmt.Printf("  Active Users: %s\n", activeUsersURL)
	fmt.Printf("  Inactive Users: %s\n", inactiveUsersURL)
	fmt.Printf("  Admin Search: %s\n", adminSearchURL)
}

func httpClientExample() {
	// Create a client for making HTTP requests
	httpClient := &http.Client{}
	
	// Build URL using go-paginate client
	paginateClient := client.New("https://jsonplaceholder.typicode.com/users")
	url := paginateClient.
		Page(1).
		Limit(5).
		BuildURL()
	
	fmt.Printf("  Making request to: %s\n", url)
	
	// Make the HTTP request
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Printf("  Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	fmt.Printf("  Response Status: %s\n", resp.Status)
	fmt.Printf("  Content-Type: %s\n", resp.Header.Get("Content-Type"))
	
	// You can also get just the query parameters for use with other HTTP libraries
	params := paginateClient.GetParams()
	fmt.Printf("  Query Parameters: %v\n", params)
}

func queryStringOnlyExample() {
	// Sometimes you only need the query string, not the full URL
	client := client.New("") // Empty base URL
	
	queryString := client.
		Page(2).
		Limit(50).
		Search("golang").
		SearchFields("title", "description").
		Sort("-created_at").
		Eq("published", true).
		Gte("views", 100).
		BuildQueryString()
	
	fmt.Printf("  Query String: %s\n", queryString)
	
	// You can use this query string with any base URL
	fullURL1 := fmt.Sprintf("https://api1.example.com/posts?%s", queryString)
	fullURL2 := fmt.Sprintf("https://api2.example.com/articles?%s", queryString)
	
	fmt.Printf("  Full URL 1: %s\n", fullURL1)
	fmt.Printf("  Full URL 2: %s\n", fullURL2)
	
	// Parse query string back to url.Values if needed
	parsedParams, err := url.ParseQuery(queryString)
	if err != nil {
		fmt.Printf("  Error parsing query string: %v\n", err)
		return
	}
	
	fmt.Printf("  Parsed page parameter: %s\n", parsedParams.Get("page"))
	fmt.Printf("  Parsed limit parameter: %s\n", parsedParams.Get("limit"))
}