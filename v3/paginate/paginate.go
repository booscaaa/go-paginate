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
	WhereArgs      []any
	WhereCombining string
	Schema         string
	Table          string
	Struct         any
	MapArgs        map[string]any
	NoOffset       bool
	// New filter fields
	LikeOr  map[string][]string
	LikeAnd map[string][]string
	EqOr    map[string][]any
	EqAnd   map[string][]any
	Gte     map[string]any
	Gt      map[string]any
	Lte     map[string]any
	Lt      map[string]any
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
func WithMapArgs(mapArgs map[string]any) Option {
	return func(params *QueryParams) {
		params.MapArgs = mapArgs
	}
}

// WithStruct sets the Struct option.
func WithStruct(s any) Option {
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
func WithWhereClause(clause string, args ...any) Option {
	return func(params *QueryParams) {
		params.WhereClauses = append(params.WhereClauses, clause)
		params.WhereArgs = append(params.WhereArgs, args...)
	}
}

// WithLikeOr sets the LikeOr filter.
func WithLikeOr(likeOr map[string][]string) Option {
	return func(params *QueryParams) {
		params.LikeOr = likeOr
	}
}

// WithLikeAnd sets the LikeAnd filter.
func WithLikeAnd(likeAnd map[string][]string) Option {
	return func(params *QueryParams) {
		params.LikeAnd = likeAnd
	}
}

// WithEqOr sets the EqOr filter.
func WithEqOr(eqOr map[string][]any) Option {
	return func(params *QueryParams) {
		params.EqOr = eqOr
	}
}

// WithEqAnd sets the EqAnd filter.
func WithEqAnd(eqAnd map[string][]any) Option {
	return func(params *QueryParams) {
		params.EqAnd = eqAnd
	}
}

// WithSearchOr is deprecated, use WithLikeOr instead.
func WithSearchOr(searchOr map[string][]string) Option {
	return WithLikeOr(searchOr)
}

// WithSearchAnd is deprecated, use WithLikeAnd instead.
func WithSearchAnd(searchAnd map[string][]string) Option {
	return WithLikeAnd(searchAnd)
}

// WithEqualsOr is deprecated, use WithEqOr instead.
func WithEqualsOr(equalsOr map[string][]any) Option {
	return WithEqOr(equalsOr)
}

// WithEqualsAnd is deprecated, use WithEqAnd instead.
func WithEqualsAnd(equalsAnd map[string][]any) Option {
	return WithEqAnd(equalsAnd)
}

// WithGte sets the Gte (greater than or equal) filter.
func WithGte(gte map[string]any) Option {
	return func(params *QueryParams) {
		params.Gte = gte
	}
}

// WithGt sets the Gt (greater than) filter.
func WithGt(gt map[string]any) Option {
	return func(params *QueryParams) {
		params.Gt = gt
	}
}

// WithLte sets the Lte (less than or equal) filter.
func WithLte(lte map[string]any) Option {
	return func(params *QueryParams) {
		params.Lte = lte
	}
}

// WithLt sets the Lt (less than) filter.
func WithLt(lt map[string]any) Option {
	return func(params *QueryParams) {
		params.Lt = lt
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
func (params *QueryParams) GenerateSQL() (string, []any) {
	var clauses []string
	var args []any

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
func (params *QueryParams) GenerateCountQuery() (string, []any) {
	var clauses []string
	var args []any

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
func (params *QueryParams) buildWhereClauses() ([]string, []any) {
	var whereClauses []string
	var args []any

	// Search conditions (legacy)
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

	// LikeOr conditions
	if len(params.LikeOr) > 0 {
		for field, values := range params.LikeOr {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" && len(values) > 0 {
				var searchConditions []string
				for _, value := range values {
					searchConditions = append(searchConditions, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
					args = append(args, "%"+value+"%")
				}
				whereClauses = append(whereClauses, "("+strings.Join(searchConditions, " OR ")+")")
			}
		}
	}

	// LikeAnd conditions
	if len(params.LikeAnd) > 0 {
		for field, values := range params.LikeAnd {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" && len(values) > 0 {
				var searchConditions []string
				for _, value := range values {
					searchConditions = append(searchConditions, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
					args = append(args, "%"+value+"%")
				}
				whereClauses = append(whereClauses, "("+strings.Join(searchConditions, " AND ")+")")
			}
		}
	}

	// EqOr conditions
	if len(params.EqOr) > 0 {
		for field, values := range params.EqOr {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" && len(values) > 0 {
				var equalsConditions []string
				for _, value := range values {
					equalsConditions = append(equalsConditions, fmt.Sprintf("%s = ?", columnName))
					args = append(args, value)
				}
				whereClauses = append(whereClauses, "("+strings.Join(equalsConditions, " OR ")+")")
			}
		}
	}

	// EqAnd conditions
	if len(params.EqAnd) > 0 {
		for field, values := range params.EqAnd {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" && len(values) > 0 {
				var equalsConditions []string
				for _, value := range values {
					equalsConditions = append(equalsConditions, fmt.Sprintf("%s = ?", columnName))
					args = append(args, value)
				}
				whereClauses = append(whereClauses, "("+strings.Join(equalsConditions, " AND ")+")")
			}
		}
	}

	// Gte conditions
	if len(params.Gte) > 0 {
		for field, value := range params.Gte {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				whereClauses = append(whereClauses, fmt.Sprintf("%s >= ?", columnName))
				args = append(args, value)
			}
		}
	}

	// Gt conditions
	if len(params.Gt) > 0 {
		for field, value := range params.Gt {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				whereClauses = append(whereClauses, fmt.Sprintf("%s > ?", columnName))
				args = append(args, value)
			}
		}
	}

	// Lte conditions
	if len(params.Lte) > 0 {
		for field, value := range params.Lte {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				whereClauses = append(whereClauses, fmt.Sprintf("%s <= ?", columnName))
				args = append(args, value)
			}
		}
	}

	// Lt conditions
	if len(params.Lt) > 0 {
		for field, value := range params.Lt {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				whereClauses = append(whereClauses, fmt.Sprintf("%s < ?", columnName))
				args = append(args, value)
			}
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
		return ""
	}

	var sortClauses []string
	for i, column := range params.SortColumns {
		columnName := getFieldName(column, "json", "paginate", params.Struct)
		if columnName != "" {
			direction := "ASC"
			if strings.ToUpper(params.SortDirections[i]) == "DESC" {
				direction = "DESC"
			}
			sortClauses = append(sortClauses, fmt.Sprintf("%s %s", columnName, direction))
		}
	}

	if len(sortClauses) > 0 {
		return "ORDER BY " + strings.Join(sortClauses, ", ")
	}
	return ""
}

// buildLimitOffsetClause constructs the LIMIT and OFFSET clauses.
func (params *QueryParams) buildLimitOffsetClause() (string, []any) {
	var clauses []string
	var args []any

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
func replacePlaceholders(query string, args []any) (string, []any) {
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
func getFieldName(tag, key, keyTarget string, s any) string {
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
