package main

import (
	"fmt"
	"log"

	"github.com/booscaaa/go-paginate/v3/paginate"
)

// User struct used for example.
type User struct {
	ID    int    `json:"id" paginate:"users.id"`
	Name  string `json:"name" paginate:"users.name"`
	Email string `json:"email" paginate:"users.email"`
	Age   int    `json:"age" paginate:"users.age"`
}

func main() {
	// Example 1: LikeOr - search for "vini" OR "joao" in the name field
	p1, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithLikeOr(map[string][]string{
			"name": {"vini", "joao"},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	query1, args1 := p1.GenerateSQL()
	fmt.Println("Example 1 - LikeOr:")
	fmt.Printf("Query: %s\n", query1)
	fmt.Printf("Args: %v\n\n", args1)

	// Example 2: LikeAnd - search for "john" AND "doe" in the name field
	p2, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithLikeAnd(map[string][]string{
			"name": {"john", "doe"},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	query2, args2 := p2.GenerateSQL()
	fmt.Println("Example 2 - LikeAnd:")
	fmt.Printf("Query: %s\n", query2)
	fmt.Printf("Args: %v\n\n", args2)

	// Example 3: EqOr - age equal to 25 OR 30 OR 35
	p3, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithEqOr(map[string][]any{
			"age": {25, 30, 35},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	query3, args3 := p3.GenerateSQL()
	fmt.Println("Example 3 - EqOr:")
	fmt.Printf("Query: %s\n", query3)
	fmt.Printf("Args: %v\n\n", args3)

	// Example 4: EqAnd - ID equal to 1 AND 2 (normally doesn't make sense, but it's possible)
	p4, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithEqAnd(map[string][]any{
			"id": {1, 2},
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	query4, args4 := p4.GenerateSQL()
	fmt.Println("Example 4 - EqAnd:")
	fmt.Printf("Query: %s\n", query4)
	fmt.Printf("Args: %v\n\n", args4)

	// Example 5: Comparison filters (Gte, Gt, Lte, Lt)
	p5, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithGte(map[string]any{"age": 18}), // age >= 18
		paginate.WithLte(map[string]any{"age": 65}), // age <= 65
		paginate.WithGt(map[string]any{"id": 0}),    // id > 0
		paginate.WithLt(map[string]any{"id": 1000}), // id < 1000
	)
	if err != nil {
		log.Fatal(err)
	}

	query5, args5 := p5.GenerateSQL()
	fmt.Println("Example 5 - Comparison Filters:")
	fmt.Printf("Query: %s\n", query5)
	fmt.Printf("Args: %v\n\n", args5)

	// Example 6: Combining multiple filters
	p6, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithLikeOr(map[string][]string{
			"name": {"vini", "joao"},
		}),
		paginate.WithEqOr(map[string][]any{
			"age": {25, 30},
		}),
		paginate.WithGte(map[string]any{"id": 1}),
		paginate.WithPage(2),
		paginate.WithItemsPerPage(5),
	)
	if err != nil {
		log.Fatal(err)
	}

	query6, args6 := p6.GenerateSQL()
	fmt.Println("Example 6 - Combined Filters:")
	fmt.Printf("Query: %s\n", query6)
	fmt.Printf("Args: %v\n\n", args6)

	// Example 7: Using JSON format as requested
	// {"likeor": {"nome": ["vini", "joao"]}}
	likeOrData := map[string][]string{
		"name": {"vini", "joao"},
	}

	p7, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithLikeOr(likeOrData),
	)
	if err != nil {
		log.Fatal(err)
	}

	query7, args7 := p7.GenerateSQL()
	count7, countArgs7 := p7.GenerateCountQuery()
	fmt.Println("Example 7 - JSON Format (likeor):")
	fmt.Printf("Query: %s\n", query7)
	fmt.Printf("Args: %v\n", args7)
	fmt.Printf("Count Query: %s\n", count7)
	fmt.Printf("Count Args: %v\n\n", countArgs7)
}
