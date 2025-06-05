package main

import (
	"fmt"
	"log"

	"github.com/booscaaa/go-paginate/v3/paginate"
)

// User represents a user model
type User struct {
	ID        int    `json:"id" paginate:"users.id"`
	Name      string `json:"name" paginate:"users.name"`
	Email     string `json:"email" paginate:"users.email"`
	Age       int    `json:"age" paginate:"users.age"`
	Status    string `json:"status" paginate:"users.status"`
	Salary    int    `json:"salary" paginate:"users.salary"`
	DeptID    int    `json:"dept_id" paginate:"users.dept_id"`
	DeptName  string `json:"dept_name" paginate:"users.dept_name"`
	CreatedAt string `json:"created_at" paginate:"users.created_at"`
}

func main() {
	fmt.Println("=== Exemplos da Nova API Fluente ===")
	fmt.Println()

	// Exemplo 1: Uso básico
	fmt.Println("1. Uso Básico:")
	basicExample()
	fmt.Println()

	// Exemplo 2: Filtros avançados
	fmt.Println("2. Filtros Avançados:")
	advancedFiltersExample()
	fmt.Println()

	// Exemplo 3: Joins
	fmt.Println("3. Joins:")
	joinsExample()
	fmt.Println()

	// Exemplo 4: A partir de JSON
	fmt.Println("4. A partir de JSON:")
	fromJSONExample()
	fmt.Println()

	// Exemplo 5: Comparação com API antiga
	fmt.Println("5. Comparação com API Antiga:")
	comparisonExample()
	fmt.Println()

	// Exemplo 6: Filtros complexos combinados
	fmt.Println("6. Filtros Complexos Combinados:")
	complexFiltersExample()
}

func basicExample() {
	// Nova API fluente - muito mais simples!
	sql, args, err := paginate.NewBuilder().
		Table("users").
		Model(&User{}).
		Page(2).
		Limit(20).
		Search("john", "name", "email").
		OrderBy("name").
		OrderByDesc("created_at").
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
}

func advancedFiltersExample() {
	// Filtros avançados com a nova API
	sql, args, err := paginate.NewBuilder().
		Table("users").
		Model(&User{}).
		WhereEquals("status", "active").
		WhereIn("dept_id", 1, 2, 3).
		WhereGreaterThan("age", 25).
		WhereLessThanOrEqual("salary", 100000).
		WhereBetween("created_at", "2023-01-01", "2023-12-31").
		SearchOr("name", "John", "Jane").
		SearchAnd("email", "@company.com").
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
}

func joinsExample() {
	// Joins simplificados
	sql, args, err := paginate.NewBuilder().
		Table("users u").
		Model(&User{}).
		Select("u.id", "u.name", "u.email", "d.name as dept_name").
		LeftJoin("departments d", "u.dept_id = d.id").
		InnerJoin("roles r", "u.role_id = r.id").
		Where("r.name = ?", "admin").
		WhereEquals("u.status", "active").
		OrderBy("u.name").
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
}

func fromJSONExample() {
	// Construir a partir de JSON (útil para APIs REST)
	jsonQuery := `{
		"page": 1,
		"limit": 10,
		"equals_or": {
			"status": ["active", "pending"]
		},
		"search_or": {
			"name": ["John", "Jane"],
			"email": ["@company.com", "@gmail.com"]
		},
		"equals_and": {
			"email": ["@company.com", "@gmail.com"]
		},
		"gte": {
			"age": 18
		},
		"lte": {
			"salary": 150000
		},
		"sort": ["-name", "created_at"]
	}`

	sql, args, err := paginate.NewBuilder().
		Table("users").
		Model(&User{}).
		InnerJoin("departments d", "u.dept_id = d.id").
		FromJSON(jsonQuery).
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("JSON Query: %s\n", jsonQuery)
	fmt.Printf("SQL: %s\n", sql)
	fmt.Printf("Args: %v\n", args)
}

func comparisonExample() {
	fmt.Println("=== API Antiga ===")
	// API antiga - mais verbosa
	oldParams, err := paginate.NewPaginator(
		paginate.WithTable("users"),
		paginate.WithStruct(&User{}),
		paginate.WithPage(1),
		paginate.WithItemsPerPage(10),
		paginate.WithSearch("john"),
		paginate.WithSearchFields([]string{"name", "email"}),
		paginate.WithEqualsOr(map[string][]any{
			"status": {"active", "pending"},
		}),
		paginate.WithGte(map[string]any{
			"age": 18,
		}),
	)

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	oldSQL, oldArgs := oldParams.GenerateSQL()
	fmt.Printf("SQL: %s\n", oldSQL)
	fmt.Printf("Args: %v\n", oldArgs)

	fmt.Println("\n=== Nova API Fluente ===")
	// Nova API - muito mais limpa!
	newSQL, newArgs, err := paginate.NewBuilder().
		Table("users").
		Model(&User{}).
		Page(1).
		Limit(10).
		Search("john", "name", "email").
		WhereIn("status", "active", "pending").
		WhereGreaterThanOrEqual("age", 18).
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("SQL: %s\n", newSQL)
	fmt.Printf("Args: %v\n", newArgs)

	// Ambas as APIs geram o mesmo resultado!
	fmt.Printf("\nResultados idênticos: %v\n", oldSQL == newSQL)
}

func complexFiltersExample() {
	// Exemplo de filtros muito complexos de forma simples
	sql, args, err := paginate.NewBuilder().
		Table("users u").
		Model(&User{}).
		Select("u.*", "d.name as department_name", "r.name as role_name").
		LeftJoin("departments d", "u.dept_id = d.id").
		LeftJoin("roles r", "u.role_id = r.id").
		// Busca em múltiplos campos
		SearchOr("name", "John", "Jane", "Bob").
		SearchAnd("email", "@company.com").
		// Filtros de igualdade
		WhereIn("u.status", "active", "pending").
		WhereIn("d.type", "engineering", "product").
		// Filtros de comparação
		WhereGreaterThanOrEqual("u.age", 21).
		WhereLessThan("u.age", 65).
		WhereGreaterThan("u.salary", 50000).
		WhereLessThanOrEqual("u.salary", 200000).
		// Filtros customizados
		Where("u.created_at >= ?", "2023-01-01").
		Where("u.last_login_at IS NOT NULL").
		// Ordenação
		OrderBy("d.name").
		OrderBy("u.name").
		OrderByDesc("u.salary").
		// Paginação
		Page(1).
		Limit(25).
		BuildSQL()

	if err != nil {
		log.Printf("Erro: %v", err)
		return
	}

	fmt.Printf("SQL Complexo: %s\n", sql)
	fmt.Printf("Args: %v\n", args)

	// Também podemos gerar a query de contagem
	countSQL, countArgs, err := paginate.NewBuilder().
		Table("users u").
		Model(&User{}).
		LeftJoin("departments d", "u.dept_id = d.id").
		LeftJoin("roles r", "u.role_id = r.id").
		SearchOr("name", "John", "Jane", "Bob").
		SearchAnd("email", "@company.com").
		WhereIn("u.status", "active", "pending").
		WhereIn("d.type", "engineering", "product").
		WhereGreaterThanOrEqual("u.age", 21).
		WhereLessThan("u.age", 65).
		WhereGreaterThan("u.salary", 50000).
		WhereLessThanOrEqual("u.salary", 200000).
		Where("u.created_at >= ?", "2023-01-01").
		Where("u.last_login_at IS NOT NULL").
		BuildCountSQL()

	if err != nil {
		log.Printf("Erro na query de contagem: %v", err)
		return
	}

	fmt.Printf("\nSQL de Contagem: %s\n", countSQL)
	fmt.Printf("Args de Contagem: %v\n", countArgs)
}
