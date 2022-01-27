package paginate

import (
	"fmt"
	"strconv"
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
	status       int
	showStatus   bool
	statusField  string
}

func Paginate(
	query string,
) Pagination {
	pagination := Pagination{
		query:        query,
		where:        " WHERE 1=1 ",
		sort:         []string{},
		descending:   []string{},
		page:         1,
		itemsPerPage: 10,
		search:       "",
		status:       0,
		showStatus:   false,
		statusField:  "",
	}

	return pagination
}

func (pagination Pagination) WhereArgs(whereArgs string) Pagination {
	pagination.where = " WHERE " + whereArgs
	return pagination
}

func (pagination Pagination) Desc(desc []string) Pagination {
	pagination.descending = desc
	return pagination
}

func (pagination Pagination) Sort(sort []string) Pagination {
	pagination.sort = sort
	return pagination
}

func (pagination Pagination) Page(page int) Pagination {
	pagination.page = page
	return pagination
}

func (pagination Pagination) RowsPerPage(rows int) Pagination {
	pagination.itemsPerPage = rows
	return pagination
}

func (pagination Pagination) SearchBy(search string) Pagination {
	pagination.search = search
	return pagination
}

func (pagination Pagination) ManageStatusBy(statusField string) Pagination {
	pagination.showStatus = true
	pagination.statusField = statusField
	pagination.status = 1
	return pagination
}

func (pagination Pagination) Query() (*string, *string, error) {
	query := pagination.query
	countQuery := generateQueryCount(query, "SELECT", "FROM")
	query += pagination.where
	countQuery += pagination.where

	offset := (pagination.page * pagination.itemsPerPage) - pagination.itemsPerPage
	searchFields := getSearchFieldsBetween(query, "SELECT", "FROM")

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

	if pagination.showStatus {
		if pagination.status != 0 {
			query += " and " + pagination.statusField + " = " + strconv.Itoa(pagination.status) + " "
			countQuery += " and " + pagination.statusField + " = " + strconv.Itoa(pagination.status) + " "
		}
	}

	if pagination.search != "" {
		for i, p := range searchFields {

			if i == 0 {
				countQuery += "and ((" + p + "::TEXT ilike '%" + pagination.search + "%') "
				query += "and ((" + p + "::TEXT ilike '%" + pagination.search + "%') "
			} else {
				countQuery += "or (" + p + "::TEXT ilike '%" + pagination.search + "%') "
				query += "or (" + p + "::TEXT ilike '%" + pagination.search + "%') "
			}

		}
		countQuery += ") "
		query += ") "
	}

	if len(pagination.sort) > 0 && pagination.sort[0] != "" {
		query += `ORDER BY `

		for s, sort := range pagination.sort {
			if s == len(pagination.sort)-1 {
				query += sort + " " + descs[s] + ` `
			} else {
				query += sort + " " + descs[s] + `, `
			}
		}
	}

	if pagination.itemsPerPage > -1 {
		query += fmt.Sprintf(" LIMIT %v OFFSET %v;", pagination.itemsPerPage, offset)
	}

	return &query, &countQuery, nil
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

func generateQueryCount(str string, start string, end string) (result string) {
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
		}
	}

	return strings.ReplaceAll(str, columns, " COUNT("+fieldWhithID+") ")
}
