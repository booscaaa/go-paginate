package paginate

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// QueryParams contains the parameters for the paginated query.
type QueryParams struct {
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
	Schema         string
	Table          string
	Struct         interface{}
	MapArgs        map[string]interface{}
	NoOffset       bool
}

// Option is a function that configures options in QueryParams.
type Option func(*QueryParams)

// WithNoOffset sets the NoOffset option.
func WithNoOffset(noOffset bool) Option {
	return func(params *QueryParams) {
		params.NoOffset = noOffset
	}
}

// WithMapArgs sets the MapArgs option.
func WithMapArgs(mapArgs map[string]interface{}) Option {
	return func(params *QueryParams) {
		params.MapArgs = mapArgs
	}
}

// WithStruct sets the Struct option.
func WithStruct(s interface{}) Option {
	return func(params *QueryParams) {
		params.Struct = s
	}
}

// WithSchema sets the Schema option.
func WithSchema(schema string) Option {
	return func(params *QueryParams) {
		params.Schema = schema
	}
}

// WithTable sets the Table option.
func WithTable(table string) Option {
	return func(params *QueryParams) {
		params.Table = table
	}
}

// WithPage sets the Page option.
func WithPage(page int) Option {
	return func(params *QueryParams) {
		params.Page = page
	}
}

// WithItemsPerPage sets the ItemsPerPage option.
func WithItemsPerPage(itemsPerPage int) Option {
	return func(params *QueryParams) {
		params.ItemsPerPage = itemsPerPage
	}
}

// WithSearch sets the Search option.
func WithSearch(search string) Option {
	return func(params *QueryParams) {
		params.Search = search
	}
}

// WithSearchFields sets the SearchFields option.
func WithSearchFields(searchFields []string) Option {
	return func(params *QueryParams) {
		params.SearchFields = searchFields
	}
}

// WithVacuum sets the Vacuum option.
func WithVacuum(vacuum bool) Option {
	return func(params *QueryParams) {
		params.Vacuum = vacuum
	}
}

// WithColumn adds a column to the Columns option.
func WithColumn(column string) Option {
	return func(params *QueryParams) {
		params.Columns = append(params.Columns, column)
	}
}

// WithSort sets the SortColumns and SortDirections options.
func WithSort(sortColumns, sortDirections []string) Option {
	return func(params *QueryParams) {
		params.SortColumns = sortColumns
		params.SortDirections = sortDirections
	}
}

// WithJoin adds a join clause to the Joins option.
func WithJoin(join string) Option {
	return func(params *QueryParams) {
		params.Joins = append(params.Joins, join)
	}
}

// WithWhereCombining sets the WhereCombining option.
func WithWhereCombining(combining string) Option {
	return func(params *QueryParams) {
		params.WhereCombining = combining
	}
}

// WithWhereClause adds a where clause and its arguments.
func WithWhereClause(clause string, args ...interface{}) Option {
	return func(params *QueryParams) {
		params.WhereClauses = append(params.WhereClauses, clause)
		params.WhereArgs = append(params.WhereArgs, args...)
	}
}

// NewPaginator creates a new QueryParams instance with the given options.
func NewPaginator(options ...Option) (*QueryParams, error) {
	params := &QueryParams{
		Page:           1,
		ItemsPerPage:   10,
		WhereCombining: "AND",
		NoOffset:       false,
	}

	// Apply options
	for _, option := range options {
		option(params)
	}

	// Validation
	if params.Table == "" {
		return nil, errors.New("principal table is required")
	}

	if params.Struct == nil {
		return nil, errors.New("struct is required")
	}

	return params, nil
}

// GenerateSQL generates the paginated SQL query and its arguments.
func (params *QueryParams) GenerateSQL() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	// SELECT clause
	selectClause := "SELECT "
	if len(params.Columns) > 0 {
		selectClause += strings.Join(params.Columns, ", ")
	} else {
		selectClause += "*"
	}
	clauses = append(clauses, selectClause)

	// FROM clause
	fromClause := fmt.Sprintf("FROM %s", params.Table)
	if params.Schema != "" {
		fromClause = fmt.Sprintf("FROM %s.%s", params.Schema, params.Table)
	}
	clauses = append(clauses, fromClause)

	// JOIN clauses
	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	// WHERE clause
	whereClauses, whereArgs := params.buildWhereClauses()
	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
		args = append(args, whereArgs...)
	}

	// ORDER BY clause
	orderClause := params.buildOrderClause()
	if orderClause != "" {
		clauses = append(clauses, orderClause)
	}

	// LIMIT and OFFSET
	limitOffsetClause, limitOffsetArgs := params.buildLimitOffsetClause()
	clauses = append(clauses, limitOffsetClause)
	args = append(args, limitOffsetArgs...)

	// Combine all clauses
	query := strings.Join(clauses, " ")

	// Replace placeholders
	query, args = replacePlaceholders(query, args)
	return query, args
}

// GenerateCountQuery generates the SQL query for counting total records.
func (params *QueryParams) GenerateCountQuery() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	// SELECT COUNT clause
	countSelectClause := "SELECT COUNT(id)"
	idColumnName := getFieldName("id", "json", "paginate", params.Struct)
	if idColumnName != "" {
		countSelectClause = fmt.Sprintf("SELECT COUNT(%s)", idColumnName)
	}
	clauses = append(clauses, countSelectClause)

	// FROM clause
	fromClause := fmt.Sprintf("FROM %s", params.Table)
	if params.Schema != "" {
		fromClause = fmt.Sprintf("FROM %s.%s", params.Schema, params.Table)
	}
	clauses = append(clauses, fromClause)

	// JOIN clauses
	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	// WHERE clause
	whereClauses, whereArgs := params.buildWhereClauses()
	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
		args = append(args, whereArgs...)
	}

	// Combine all clauses
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

// buildWhereClauses constructs the WHERE clauses and arguments.
func (params *QueryParams) buildWhereClauses() ([]string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	// Search conditions
	if params.Search != "" && len(params.SearchFields) > 0 {
		var searchConditions []string
		for _, field := range params.SearchFields {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				searchConditions = append(searchConditions, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
				args = append(args, "%"+params.Search+"%")
			}
		}
		if len(searchConditions) > 0 {
			whereClauses = append(whereClauses, "("+strings.Join(searchConditions, " OR ")+")")
		}
	}

	// Additional WHERE clauses
	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	return whereClauses, args
}

// buildOrderClause constructs the ORDER BY clause.
func (params *QueryParams) buildOrderClause() string {

	if len(params.SortColumns) == 0 || len(params.SortDirections) != len(params.SortColumns) {
		fmt.Println(params.SortColumns)
		fmt.Println(params.SortDirections)
		return ""
	}

	var sortClauses []string
	for i, column := range params.SortColumns {
		columnName := getFieldName(column, "json", "paginate", params.Struct)
		if columnName != "" {
			direction := "ASC"
			if strings.ToLower(params.SortDirections[i]) == "true" {
				direction = "DESC"
			}
			sortClauses = append(sortClauses, fmt.Sprintf("%s %s", columnName, direction))
		}
	}

	fmt.Println(sortClauses)
	if len(sortClauses) > 0 {
		return "ORDER BY " + strings.Join(sortClauses, ", ")
	}
	return ""
}

// buildLimitOffsetClause constructs the LIMIT and OFFSET clauses.
func (params *QueryParams) buildLimitOffsetClause() (string, []interface{}) {
	var clauses []string
	var args []interface{}

	clauses = append(clauses, "LIMIT ?")
	args = append(args, params.ItemsPerPage)

	if !params.NoOffset {
		offset := (params.Page - 1) * params.ItemsPerPage
		clauses = append(clauses, "OFFSET ?")
		args = append(args, offset)
	}

	return strings.Join(clauses, " "), args
}

// Helper functions

// replacePlaceholders replaces '?' with positional placeholders like '$1', '$2', etc.
func replacePlaceholders(query string, args []interface{}) (string, []interface{}) {
	var newQuery strings.Builder
	argIndex := 1
	for _, char := range query {
		if char == '?' {
			newQuery.WriteString(fmt.Sprintf("$%d", argIndex))
			argIndex++
		} else {
			newQuery.WriteRune(char)
		}
	}
	return newQuery.String(), args
}

// getFieldName retrieves the column name from struct tags based on the given key.
func getFieldName(tag, key, keyTarget string, s interface{}) string {
	rt := reflect.TypeOf(s)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		panic("struct type required")
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tagValue := strings.Split(field.Tag.Get(key), ",")[0]
		if tagValue == tag {
			return field.Tag.Get(keyTarget)
		}
	}
	return ""
}
