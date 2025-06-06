package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/booscaaa/go-paginate/v3/paginate"
)

// User represents a user entity
type User struct {
	ID    int    `json:"id" paginate:"id"`
	Name  string `json:"name" paginate:"name"`
	Email string `json:"email" paginate:"email"`
	Age   int    `json:"age" paginate:"age"`
}

func main() {
	// Configure structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	fmt.Println("=== Go-Paginate Debug Mode Example ===")
	fmt.Println()

	// Example 1: Using environment variables
	fmt.Println("1. Testing with environment variables:")
	os.Setenv("GO_PAGINATE_DEBUG", "true")
	os.Setenv("GO_PAGINATE_DEFAULT_LIMIT", "25")
	os.Setenv("GO_PAGINATE_MAX_LIMIT", "1000")

	// Reload configuration from environment
	paginate.SetDebugMode(true)
	paginate.SetDefaultLimit(25)
	paginate.SetMaxLimit(1000)

	fmt.Println("Debug mode:", paginate.IsDebugMode())
	fmt.Println("Default limit:", paginate.GetDefaultLimit())
	fmt.Println("Max limit:", paginate.GetMaxLimit())
	fmt.Println()

	// Example 2: Building a simple query
	fmt.Println("2. Building a simple query with debug logs:")
	builder := paginate.NewBuilder().
		Table("users").
		Model(User{}).
		Page(1).
		Limit(10).
		Search("john", "name", "email").
		OrderBy("name", "ASC")

	sql, args, err := builder.BuildSQL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated SQL: %s\n", sql)
	fmt.Printf("Arguments: %v\n", args)
	fmt.Println()

	// Example 3: Building a count query
	fmt.Println("3. Building a count query with debug logs:")
	countSQL, countArgs, err := builder.BuildCountSQL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Generated Count SQL: %s\n", countSQL)
	fmt.Printf("Count Arguments: %v\n", countArgs)
	fmt.Println()

	// Example 4: Complex query with filters
	fmt.Println("4. Building a complex query with filters:")
	complexBuilder := paginate.NewBuilder().
		Table("users").
		Schema("public").
		Model(User{}).
		Page(2).
		Limit(20).
		Select("id", "name", "email", "age").
		Search("developer", "name", "email").
		WhereGreaterThanOrEqual("age", 18).
		WhereLessThan("age", 65).
		OrderBy("name", "ASC").
		OrderBy("age", "DESC")

	complexSQL, complexArgs, err := complexBuilder.BuildSQL()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Complex SQL: %s\n", complexSQL)
	fmt.Printf("Complex Arguments: %v\n", complexArgs)
	fmt.Println()

	// Example 5: Testing with debug mode disabled
	fmt.Println("5. Testing with debug mode disabled:")
	paginate.SetDebugMode(false)
	fmt.Println("Debug mode:", paginate.IsDebugMode())

	// This should not produce debug logs
	silentSQL, silentArgs, err := paginate.NewBuilder().
		Table("users").
		Model(User{}).
		Page(1).
		Limit(5).
		BuildSQL()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Silent SQL (no debug logs): %s\n", silentSQL)
	fmt.Printf("Silent Arguments: %v\n", silentArgs)
	fmt.Println()

	// Example 6: Re-enable debug mode
	fmt.Println("6. Re-enabling debug mode:")
	paginate.SetDebugMode(true)

	// This should produce debug logs again
	finalSQL, finalArgs, err := paginate.NewBuilder().
		Table("products").
		Model(struct {
			ID    int    `json:"id" paginate:"id"`
			Name  string `json:"name" paginate:"name"`
			Price float64 `json:"price" paginate:"price"`
		}{}).
		Page(1).
		Limit(15).
		WhereGreaterThan("price", 10.0).
		OrderBy("price", "DESC").
		BuildSQL()

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Final SQL: %s\n", finalSQL)
	fmt.Printf("Final Arguments: %v\n", finalArgs)
	fmt.Println()

	fmt.Println("=== Debug Mode Example Complete ===")
}