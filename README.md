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

## Functionalities

- Generate query and query count pagination for postgresql

<br>

## Starting

### Installation

```sh
go get github.com/booscaaa/go-paginate
```

<br>

### Usage

```go
import "github.com/booscaaa/go-paginate/paginate"

query, queryCount, err := paginate.
		Paginate("SELECT id, name FROM user").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(1).
		RowsPerPage(50).
		SearchBy("vinicius").
		Query()

if err != nil {
  //handler error
}

// else use the query and queryCount to paginate with pq, pgx or other
```

## Contributing

You can send how many PR's do you want, I'll be glad to analyze and accept them! And if you have any question about the project...

Email-me: boscardinvinicius@gmail.com

Connect with me at [LinkedIn](https://www.linkedin.com/in/booscaaa/)

Thank you!

## License

This project is licensed under the MIT License - see the [LICENSE.md](https://github.com/booscaaa/go-paginate/blob/master/LICENSE) file for details
