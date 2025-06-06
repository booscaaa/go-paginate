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
	// Exemplo 1: LikeOr - busca por "vini" OU "joao" no campo nome
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
	fmt.Println("Exemplo 1 - LikeOr:")
	fmt.Printf("Query: %s\n", query1)
	fmt.Printf("Args: %v\n\n", args1)

	// Exemplo 2: LikeAnd - busca por "john" E "doe" no campo nome
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
	fmt.Println("Exemplo 2 - LikeAnd:")
	fmt.Printf("Query: %s\n", query2)
	fmt.Printf("Args: %v\n\n", args2)

	// Exemplo 3: EqOr - idade igual a 25 OU 30 OU 35
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
	fmt.Println("Exemplo 3 - EqOr:")
	fmt.Printf("Query: %s\n", query3)
	fmt.Printf("Args: %v\n\n", args3)

	// Exemplo 4: EqAnd - ID igual a 1 E 2 (normalmente não faz sentido, mas é possível)
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
	fmt.Println("Exemplo 4 - EqAnd:")
	fmt.Printf("Query: %s\n", query4)
	fmt.Printf("Args: %v\n\n", args4)

	// Exemplo 5: Filtros de comparação (Gte, Gt, Lte, Lt)
	p5, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(User{}),
		paginate.WithGte(map[string]any{"age": 18}), // idade >= 18
		paginate.WithLte(map[string]any{"age": 65}), // idade <= 65
		paginate.WithGt(map[string]any{"id": 0}),    // id > 0
		paginate.WithLt(map[string]any{"id": 1000}), // id < 1000
	)
	if err != nil {
		log.Fatal(err)
	}

	query5, args5 := p5.GenerateSQL()
	fmt.Println("Exemplo 5 - Filtros de Comparação:")
	fmt.Printf("Query: %s\n", query5)
	fmt.Printf("Args: %v\n\n", args5)

	// Exemplo 6: Combinando múltiplos filtros
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
	fmt.Println("Exemplo 6 - Filtros Combinados:")
	fmt.Printf("Query: %s\n", query6)
	fmt.Printf("Args: %v\n\n", args6)

	// Exemplo 7: Usando o formato JSON como solicitado
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
	fmt.Println("Exemplo 7 - Formato JSON (likeor):")
	fmt.Printf("Query: %s\n", query7)
	fmt.Printf("Args: %v\n", args7)
	fmt.Printf("Count Query: %s\n", count7)
	fmt.Printf("Count Args: %v\n\n", countArgs7)
}
