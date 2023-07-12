package paginate

import (
	"fmt"
	"reflect"
	"strings"
)

type Pagination struct {
	query        string
	where        string
	sort         []string
	descending   []string
	page         int
	itemsPerPage int
	search       string
	statusField  string
	searchFields []string
	structType   interface{}
	withVacuum   bool
}

func Instance(structType interface{}) Pagination {
	pagination := Pagination{
		query:        "",
		where:        " WHERE 1=1 ",
		sort:         []string{},
		descending:   []string{},
		page:         1,
		itemsPerPage: 10,
		search:       "",
		statusField:  "",
		searchFields: []string{},
		structType:   structType,
		withVacuum:   false,
	}

	return pagination
}

func (pagination *Pagination) Query(query string) *Pagination {
	pagination.query = query
	return pagination
}

func (pagination *Pagination) WhereArgs(operation, whereArgs string) *Pagination {
	pagination.where += fmt.Sprintf(" %s %s ", operation, whereArgs)

	return pagination
}

func (pagination *Pagination) Desc(desc []string) *Pagination {
	pagination.descending = desc
	return pagination
}

func (pagination *Pagination) Sort(sort []string) *Pagination {
	pagination.sort = sort
	return pagination
}

func (pagination *Pagination) Page(page int) *Pagination {
	pagination.page = page
	return pagination
}

func (pagination *Pagination) RowsPerPage(rows int) *Pagination {
	pagination.itemsPerPage = rows
	return pagination
}

func (pagination *Pagination) SearchBy(search string, fields ...string) *Pagination {
	pagination.search = search
	pagination.searchFields = fields
	return pagination
}

func (pagination *Pagination) WithVacuum() *Pagination {
	pagination.withVacuum = true
	return pagination
}

func (pagination Pagination) Select() (*string, *string) {
	query := pagination.query
	countQuery := generateQueryCount(query, "SELECT", "FROM", pagination.withVacuum)
	query += pagination.where
	countQuery += pagination.where

	offset := (pagination.page * pagination.itemsPerPage) - pagination.itemsPerPage

	if len(pagination.searchFields) == 0 {
		pagination.searchFields = getSearchFieldsBetween(query, "SELECT", "FROM")
	}

	var descs []string

	if len(pagination.descending) == 0 {
		for range pagination.sort {
			descs = append(descs, "ASC")
		}
	} else {
		for _, desc := range pagination.descending {
			if desc == "true" {
				descs = append(descs, "DESC")
			} else {
				descs = append(descs, "ASC")
			}
		}
	}

	if pagination.search != "" && len(pagination.searchFields) > 0 {
		for i, p := range pagination.searchFields {
			if p != "" {
				p = getFieldName(p, "json", pagination.structType, "paginate")
				if i == 0 {
					countQuery += "and ((" + p + "::TEXT ilike $1) "
					query += "and ((" + p + "::TEXT ilike $1) "
				} else {
					countQuery += "or (" + p + "::TEXT ilike $1) "
					query += "or (" + p + "::TEXT ilike $1) "
				}
			}
		}

		if len(pagination.searchFields) > 0 {
			if pagination.searchFields[0] != "" {
				countQuery += ") "
				query += ") "
			}
		}
	}

	if len(pagination.sort) > 0 && pagination.sort[0] != "" {
		query += `ORDER BY `

		for s, sort := range pagination.sort {
			if s == len(pagination.sort)-1 {
				query += getFieldName(sort, "json", pagination.structType, "db") + " " + descs[s] + ` `
			} else {
				query += getFieldName(sort, "json", pagination.structType, "db") + " " + descs[s] + `, `
			}
		}
	}

	if pagination.itemsPerPage > -1 {
		query += fmt.Sprintf(" LIMIT %v OFFSET %v;", pagination.itemsPerPage, offset)
	}

	if pagination.withVacuum {
		countQuery = "SELECT count_estimate($" + strings.ReplaceAll(
			countQuery,
			"COUNT(dl.id)",
			"1",
		) + "$);"

		countQuery = strings.ReplaceAll(
			countQuery,
			"'",
			"''",
		)

		countQuery = strings.ReplaceAll(
			countQuery,
			"$",
			"'",
		)
	}

	return &query, &countQuery
}

func getSearchFieldsBetween(str string, start string, end string) (result []string) {
	a := strings.SplitAfterN(str, start, 2)
	b := strings.SplitAfterN(a[len(a)-1], end, 2)
	fields := strings.Split(strings.Replace((b[0][0:len(b[0])-len(end)]), " ", "", -1), ",")

	searchFields := []string{}
	for _, field := range fields {
		if !strings.Contains(field, "*") {
			searchFields = append(searchFields, field)
		}
	}

	return searchFields
}

func generateQueryCount(str string, start string, end string, vacuum bool) (result string) {
	a := strings.SplitAfterN(str, start, 2)
	b := strings.SplitAfterN(a[len(a)-1], end, 2)
	columns := b[0][0 : len(b[0])-len(end)]

	fields := strings.Split(strings.Replace((b[0][0:len(b[0])-len(end)]), " ", "", -1), ",")

	fieldWhithID := "id"
	for _, field := range fields {
		if !strings.Contains(field, ".*") {
			if strings.Contains(field, "id") {
				fieldWhithID = field
				break
			}
		} else if strings.Contains(field, ".*") {
			fieldWhithID = strings.ReplaceAll(field, ".*", ".id")
			break
		}
	}

	if vacuum {
		return strings.ReplaceAll(str, columns, " 1 ")
	}

	return strings.ReplaceAll(str, columns, " COUNT("+fieldWhithID+") ")
}

func getFieldName(tag, key string, s interface{}, keyTarget string) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == tag {
			return f.Tag.Get(keyTarget)
		}
	}
	return ""
}
