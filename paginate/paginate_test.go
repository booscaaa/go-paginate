package paginate_test

import (
	"testing"

	"github.com/booscaaa/go-paginate/paginate"
)

func TestPaginate(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE 1=1  and v.status = 1 and ((t.id::TEXT ilike '%vinicius%') ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE 1=1  and v.status = 1 and ((t.id::TEXT ilike '%vinicius%') ) "
	query, queryCount, err := paginate.
		Paginate("SELECT t.* FROM test t").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50).
		SearchBy("vinicius", "t.id").
		ManageStatusBy("v.status").
		Query()

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	t.Log(*query)
	t.Log(*queryCount)

	if queryString != *query {
		t.Errorf("Wrong query")
		return
	}

	if queryCountString != *queryCount {
		t.Errorf("Wrong query count")
	}
}

func TestPaginateWithArgs(t *testing.T) {
	queryString := "SELECT t.* FROM test t WHERE t.name = 'jhon' and v.status = 1 and ((t.id::TEXT ilike '%vinicius%') ) ORDER BY name DESC, last_name ASC  LIMIT 50 OFFSET 100;"
	queryCountString := "SELECT COUNT(t.id) FROM test t WHERE t.name = 'jhon' and v.status = 1 and ((t.id::TEXT ilike '%vinicius%') ) "
	query, queryCount, err := paginate.
		Paginate("SELECT t.* FROM test t").
		Sort([]string{"name", "last_name"}).
		Desc([]string{"true", "false"}).
		Page(3).
		RowsPerPage(50).
		SearchBy("vinicius", "t.id").
		ManageStatusBy("v.status").
		WhereArgs("t.name = 'jhon'").
		Query()

	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	t.Log(*query)
	t.Log(*queryCount)

	if queryString != *query {
		t.Errorf("Wrong query")
	}

	if queryCountString != *queryCount {
		t.Errorf("Wrong query count")
	}
}
