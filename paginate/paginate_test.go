package paginate_test

import (
	"testing"

	"github.com/booscaaa/go-paginate/paginate"
)

type Test struct {
	Name     string `json:"name"     db:"name" paginate:"test.name"`
	LastName string `json:"lastName" db:"last_name" paginate:"test.last_name"`
}

func TestPaginate(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE 1=1 and ((test.name::TEXT ilike $1) ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE 1=1 and ((test.name::TEXT ilike $1) ) "

	pagin := paginate.Instance(Test{})
	query, queryCount := pagin.
		Query("SELECT t.* FROM test t").
		Sort([]string{"name", "lastName"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50).
		SearchBy("vinicius", "name").
		Select()

	if queryString != *query {
		t.Errorf("Wrong query")
		return
	}

	if queryCountString != *queryCount {
		t.Errorf("Wrong query count")
	}
}

func TestPaginateWithArgs(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE 1=1  and test.name = 'jhon' and ((test.last_name::TEXT ilike $1) ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE 1=1  and test.name = 'jhon' and ((test.last_name::TEXT ilike $1) ) "

	pagin := paginate.Instance(Test{})

	pagin.Query("SELECT t.* FROM test t").
		Sort([]string{"name", "lastName"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50)

	pagin.WhereArgs("and", "test.name = 'jhon'")
	pagin.SearchBy("vinicius", []string{"lastName"}...)
	query, queryCount := pagin.Select()

	if queryString != *query {
		t.Errorf("Wrong query")
	}

	if queryCountString != *queryCount {
		t.Errorf("Wrong query count")
	}
}
