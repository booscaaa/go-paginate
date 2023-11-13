<p align="center">
  <h1 align="center">Go Paginate - Go package to generate query pagination</h1>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/booscaaa/go-paginate"><img alt="Reference" src="https://img.shields.io/badge/go-reference-purple?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/releases/latest"><img alt="Release" src="https://img.shields.io/github/v/release/booscaaa/go-paginate.svg?style=for-the-badge"></a>
    <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-red.svg?style=for-the-badge"></a>
    <a href="https://github.com/booscaaa/go-paginate/actions/workflows/test.yaml"><img alt="Test status" src="https://img.shields.io/github/workflow/status/booscaaa/go-paginate/Test?label=TESTS&style=for-the-badge"></a>
    <a href="https://codecov.io/gh/booscaaa/go-paginate"><img alt="Coverage" src="https://img.shields.io/codecov/c/github/booscaaa/go-paginate/master.svg?style=for-the-badge"></a>
  </p>
</p>

<br>

## Why?

This project is part of my personal portfolio, so, I'll be happy if you could provide me any feedback about the project, code, structure or anything that you can report that could make me a better developer!

Email-me: boscardinvinicius@gmail.com

Connect with me at [LinkedIn](https://www.linkedin.com/in/booscaaa/).

<br>

# Paginate Package Readme

## Overview

The `paginate` package provides a flexible and easy-to-use solution for paginated queries in Go. It allows you to construct paginated queries with various options, including sorting, filtering, and custom column selection.

## Usage

To use the `paginate` package, follow these steps:

1. **Import the package:**

   ```go
   import "github.com/booscaaa/go-paginate/v2/paginate"
   ```

2. **Create a struct to represent your database model.**

   Define a struct that mirrors your database model, with struct tags specifying the corresponding database columns.

   ```go
   type MyModel struct {
       ID           int       `json:"id" paginate:"my_table.id"`
       CreatedAt    time.Time `json:"created_at" paginate:"my_table.created_at"`
       Name         string    `json:"name" paginate:"my_table.name"`
   }
   ```

3. **Use the `PaginQuery` function to create paginated queries:**

   ```go
   // Example usage:
   params, err := paginate.PaginQuery(
       paginate.WithStruct(MyModel{}),
       paginate.WithTable("my_table"),
       paginate.WithColumn("my_table.*"),
       paginate.WithPage(2),
       paginate.WithItemsPerPage(10),
       paginate.WithSort([]string{"created_at"}, []string{"true"}),
       paginate.WithSearch("example"),
   )
   if err != nil {
       log.Fatal(err)
   }

   // Generate SQL and arguments
   sql, args := paginate.GenerateSQL(params)
   countSQL, countArgs := paginate.GenerateCountQuery(params)
   ```

   The above example will output SQL queries and arguments like:

   ```sql
   SELECT my_table.* FROM my_table WHERE (my_table.name::TEXT ILIKE $1) ORDER BY created_at DESC LIMIT $2 OFFSET $3
   ```

   SQL Arguments:

   ```
   [%example% 10 10]
   ```

   Count Query:

   ```sql
   SELECT COUNT(id) FROM my_table WHERE (my_table.name::TEXT ILIKE $1)
   ```

   Count Arguments:

   ```
   [%example%]
   ```

4. **Options and Customization:**

   You can customize your paginated query using various options such as `WithPage`, `WithItemsPerPage`, `WithSort`, `WithSearch`, `WithSearchFields`, `WithVacuum`, `WithColumn`, `WithJoin`, `WithWhereCombining`, and `WithWhereClause`. These options allow you to tailor your query to specific requirements.

   ```go
   // Example options:
   options := []paginate.Option{
       paginate.WithPage(2),
       paginate.WithItemsPerPage(20),
       paginate.WithSort([]string{"created_at"}, []string{"true"}),
       paginate.WithSearch("example"),
       paginate.WithSearchFields([]string{"name"}),
       paginate.WithVacuum(true),
       paginate.WithColumn("my_table.*"),
       paginate.WithJoin("INNER JOIN other_table ON my_table.id = other_table.my_table_id"),
       paginate.WithWhereClause("status = ?", "active"),
   }

   params, err := paginate.PaginQuery(options...)
   ```

5. **Run your query:**

   Once you've configured your paginated query, use the generated SQL and arguments to execute the query against your database.

## Options

### `WithNoOffset`

Disable OFFSET and LIMIT for pagination. Useful for scenarios where OFFSET is not performant.

### `WithMapArgs`

Pass a map of custom arguments to be used in the WHERE clause.

### `WithStruct`

Specify the database model struct to be used for generating SQL queries.

### `WithTable`

Specify the main table for the paginated query.

### `WithPage`

Set the page number for pagination.

### `WithItemsPerPage`

Set the number of items per page.

### `WithSearch`

Specify a search term to filter results.

### `WithSearchFields`

Specify fields to search within.

### `WithVacuum`

Enable or disable VACUUM optimization for the query.

### `WithColumn`

Add a custom column to the SELECT clause.

### `WithSort`

Specify sorting columns and directions.

### `WithJoin`

Add a custom JOIN clause.

### `WithWhereCombining`

Specify the combining operator for multiple WHERE clauses.

### `WithWhereClause`

Add a custom WHERE clause.

## Example

Check the provided example in the code for a comprehensive demonstration of the package's usage.

```go
// Example usage:
params, err := paginate.PaginQuery(
   // ... (options)
)
```

## Contribution

Feel free to contribute to the `paginate` package by creating issues, submitting pull requests, or providing feedback. Your contributions are highly appreciated!

## Contributing

You can send how many PR's do you want, I'll be glad to analyze and accept them! And if you have any question about the project...

Email-me: boscardinvinicius@gmail.com

Connect with me at [LinkedIn](https://www.linkedin.com/in/booscaaa/)

Thank you!

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/booscaaa/go-paginate/blob/master/LICENSE) file for details
