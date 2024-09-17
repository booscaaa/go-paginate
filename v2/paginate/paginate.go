package paginate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// paginQueryParams contains the parameters for the paginated query
type paginQueryParams struct {
	Page           int
	ItemsPerPage   int
	Search         string
	SearchFields   []string
	Vacuum         bool
	Columns        []string
	Joins          []string
	SortColumns    []string
	SortDirections []string
	WhereClauses   []string
	WhereArgs      []interface{}
	WhereCombining string
	Schema         string // New field to store the schema name
	Table          string
	Struct         interface{}
	mapArgs        map[string]any
	noOffset       bool
}

// Option is a function that configures options in paginQueryParams
type Option func(*paginQueryParams)

func WithNoOffset(noOffset bool) Option {
	return func(params *paginQueryParams) {
		params.noOffset = noOffset
	}
}

func WithMapArgs(mapArgs map[string]any) Option {
	return func(params *paginQueryParams) {
		params.mapArgs = mapArgs
	}
}

func WithStruct(s interface{}) Option {
	return func(params *paginQueryParams) {
		params.Struct = s
	}
}

// WithSchema configures the schema field
func WithSchema(schema string) Option {
	return func(params *paginQueryParams) {
		params.Schema = schema
	}
}

func WithTable(table string) Option {
	return func(params *paginQueryParams) {
		params.Table = table
	}
}

func WithPage(page int) Option {
	return func(params *paginQueryParams) {
		params.Page = page
	}
}

func WithItemsPerPage(itemsPerPage int) Option {
	return func(params *paginQueryParams) {
		params.ItemsPerPage = itemsPerPage
	}
}

func WithSearch(search string) Option {
	return func(params *paginQueryParams) {
		params.Search = search
	}
}

func WithSearchFields(searchFields []string) Option {
	return func(params *paginQueryParams) {
		params.SearchFields = searchFields
	}
}

func WithVacuum(vacuum bool) Option {
	return func(params *paginQueryParams) {
		params.Vacuum = vacuum
	}
}

func WithColumn(column string) Option {
	return func(params *paginQueryParams) {
		params.Columns = append(params.Columns, column)
	}
}

func WithSort(sortColumns []string, sortDirections []string) Option {
	return func(params *paginQueryParams) {
		params.SortColumns = sortColumns
		params.SortDirections = sortDirections
	}
}

func WithJoin(join string) Option {
	return func(params *paginQueryParams) {
		params.Joins = append(params.Joins, join)
	}
}

func WithWhereCombining(combining string) Option {
	return func(params *paginQueryParams) {
		params.WhereCombining = combining
	}
}

func WithWhereClause(clause string, args ...interface{}) Option {
	return func(params *paginQueryParams) {
		params.WhereClauses = append(params.WhereClauses, clause)
		params.WhereArgs = append(params.WhereArgs, args...)
	}
}

func PaginQuery(options ...Option) (*paginQueryParams, error) {
	params := &paginQueryParams{
		Page:           1,
		ItemsPerPage:   10,
		WhereCombining: "AND", // Default combination is "AND"
		noOffset:       false,
	}

	// Apply options
	for _, option := range options {
		option(params)
	}

	if params.Table == "" {
		return nil, errors.New("principal table is required")
	}

	if params.Struct == nil {
		return nil, errors.New("struct is required")
	}

	return params, nil
}

func GenerateSQL(params *paginQueryParams) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}

	nextArg := func() int {
		argNum := len(args) + 1
		args = append(args, nil)
		return argNum
	}

	selectClause := "SELECT "
	if len(params.Columns) > 0 {
		selectClause += strings.Join(params.Columns, ", ")
	} else {
		selectClause += "*"
	}
	clauses = append(clauses, selectClause)

	// FROM clause with schema if provided
	if params.Schema != "" {
		clauses = append(clauses, fmt.Sprintf("FROM %s.%s", params.Schema, params.Table))
	} else {
		clauses = append(clauses, fmt.Sprintf("FROM %s", params.Table))
	}

	// JOIN clauses
	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	// WHERE clause for search
	whereClauses := []string{}

	if params.Search != "" && len(params.SearchFields) > 0 {
		searchConditions := []string{}
		for _, field := range params.SearchFields {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				searchConditions = append(searchConditions, fmt.Sprintf("%s::TEXT ILIKE $%d", columnName, nextArg()))
				args[len(args)-1] = "%" + params.Search + "%"
			}
		}
		if len(searchConditions) > 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR ")))
		}
	}

	// Additional WHERE clauses
	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
	}

	// ORDER BY clause
	if len(params.SortColumns) > 0 && len(params.SortDirections) == len(params.SortColumns) {
		sortClauses := []string{}
		for i, column := range params.SortColumns {
			columnName := getFieldName(column, "json", "paginate", params.Struct)
			if columnName != "" {
				direction := "ASC"
				if params.SortDirections[i] == "true" {
					direction = "DESC"
				}
				sortClauses = append(sortClauses, fmt.Sprintf("%s %s", columnName, direction))
			}
		}
		if len(sortClauses) > 0 {
			clauses = append(clauses, "ORDER BY "+strings.Join(sortClauses, ", "))
		}
	}

	// LIMIT and OFFSET for pagination
	offset := (params.Page - 1) * params.ItemsPerPage
	clauses = append(clauses, "LIMIT $"+fmt.Sprint(nextArg()))
	args[len(args)-1] = params.ItemsPerPage

	if !params.noOffset {
		clauses = append(clauses, "OFFSET $"+fmt.Sprint(nextArg()))
		args[len(args)-1] = offset
	}

	// Join all clauses into a single SQL query
	query := strings.Join(clauses, " ")

	// Replace placeholders and return
	query, args = replacePlaceholders(query, args)
	return query, args
}

func GenerateCountQuery(params *paginQueryParams) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}

	nextArg := func() int {
		argNum := len(args) + 1
		args = append(args, nil)
		return argNum
	}

	// SELECT COUNT clause
	countSelectClause := "SELECT COUNT(id)"
	idColumnName := getFieldName("id", "json", "paginate", params.Struct)
	if idColumnName != "" {
		countSelectClause = fmt.Sprintf("SELECT COUNT(%s)", idColumnName)
	}
	clauses = append(clauses, countSelectClause)

	// FROM clause with schema if provided
	if params.Schema != "" {
		clauses = append(clauses, fmt.Sprintf("FROM %s.%s", params.Schema, params.Table))
	} else {
		clauses = append(clauses, fmt.Sprintf("FROM %s", params.Table))
	}

	// JOIN clauses
	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	// WHERE clause
	whereClauses := []string{}

	if params.Search != "" && len(params.SearchFields) > 0 {
		searchConditions := []string{}
		for _, field := range params.SearchFields {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				searchConditions = append(searchConditions, fmt.Sprintf("%s::TEXT ILIKE $%d", columnName, nextArg()))
				args[len(args)-1] = "%" + params.Search + "%"
			}
		}
		if len(searchConditions) > 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR ")))
		}
	}

	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
	}

	query := strings.Join(clauses, " ")

	// Replace placeholders
	query, args = replacePlaceholders(query, args)

	if params.Vacuum {
		countQuery := "SELECT count_estimate('" + query + "');"
		countQuery = strings.Replace(countQuery, "COUNT(id)", "1", -1)
		re := regexp.MustCompile(`(\$[0-9]+)`)
		countQuery = re.ReplaceAllStringFunc(countQuery, func(match string) string {
			return "''" + match + "''"
		})
		return countQuery, args
	}

	return query, args
}

// Helper functions to replace placeholders and extract field names

func replacePlaceholders(query string, args []interface{}) (string, []interface{}) {
	lastArg := 0
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			break
		} else if query[i] == '$' {
			lastArg, _ = strconv.Atoi(string(query[i+1]))
		}
	}
	for i := 0; i < len(query); i++ {
		if query[i] == '?' {
			query = query[:i] + "$" + strconv.Itoa(lastArg+1) + query[i+1:]
			lastArg++
		}
	}
	return query, args
}

func getFieldName(tag, key, keyTarget string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0]
		if v == tag {
			return f.Tag.Get(keyTarget)
		}
	}
	return ""
}
