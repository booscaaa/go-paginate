package paginate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// QueryParams contains the parameters for the paginated query.
// It is used internally by the Builder and returned by Build().
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
	Like           map[string][]string
	LikeOr         map[string][]string
	LikeAnd        map[string][]string
	Eq             map[string][]any
	EqOr           map[string][]any
	EqAnd          map[string][]any
	Gte            map[string]any
	Gt             map[string]any
	Lte            map[string]any
	Lt             map[string]any
	In             map[string][]any
	NotIn          map[string][]any
	Between        map[string][2]any
	IsNull         []string
	IsNotNull      []string
	GteOr          map[string]any
	GtOr           map[string]any
	LteOr          map[string]any
	LtOr           map[string]any
	InOr           map[string][]any
	NotInOr        map[string][]any
	IsNullOr       []string
	IsNotNullOr    []string
}

// GenerateSQL generates the paginated SQL query and its arguments.
func (params *QueryParams) GenerateSQL() (string, []any) {
	var clauses []string
	var args []any

	selectClause := "SELECT "
	if len(params.Columns) > 0 {
		selectClause += strings.Join(params.Columns, ", ")
	} else {
		selectClause += "*"
	}
	clauses = append(clauses, selectClause)

	fromClause := fmt.Sprintf("FROM %s", params.Table)
	if params.Schema != "" {
		fromClause = fmt.Sprintf("FROM %s.%s", params.Schema, params.Table)
	}
	clauses = append(clauses, fromClause)

	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	whereClauses, whereArgs := params.buildWhereClauses()
	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
		args = append(args, whereArgs...)
	}

	orderClause := params.buildOrderClause()
	if orderClause != "" {
		clauses = append(clauses, orderClause)
	}

	limitOffsetClause, limitOffsetArgs := params.buildLimitOffsetClause()
	clauses = append(clauses, limitOffsetClause)
	args = append(args, limitOffsetArgs...)

	query := strings.Join(clauses, " ")
	query, args = replacePlaceholders(query, args)
	logSQL("GenerateSQL", query, args)

	return query, args
}

// GenerateCountQuery generates the SQL query for counting total records.
func (params *QueryParams) GenerateCountQuery() (string, []any) {
	var clauses []string
	var args []any

	countSelectClause := "SELECT COUNT(id)"
	idColumnName := getFieldName("id", "json", "paginate", params.Struct)
	if idColumnName != "" {
		countSelectClause = fmt.Sprintf("SELECT COUNT(%s)", idColumnName)
	}
	clauses = append(clauses, countSelectClause)

	fromClause := fmt.Sprintf("FROM %s", params.Table)
	if params.Schema != "" {
		fromClause = fmt.Sprintf("FROM %s.%s", params.Schema, params.Table)
	}
	clauses = append(clauses, fromClause)

	if len(params.Joins) > 0 {
		clauses = append(clauses, strings.Join(params.Joins, " "))
	}

	whereClauses, whereArgs := params.buildWhereClauses()
	if len(whereClauses) > 0 {
		clauses = append(clauses, "WHERE "+strings.Join(whereClauses, " AND "))
		args = append(args, whereArgs...)
	}

	query := strings.Join(clauses, " ")
	query, args = replacePlaceholders(query, args)

	if params.Vacuum {
		countQuery := "SELECT count_estimate('" + query + "');"
		countQuery = strings.ReplaceAll(countQuery, "COUNT(id)", "1")
		re := regexp.MustCompile(`(\$[0-9]+)`)
		countQuery = re.ReplaceAllStringFunc(countQuery, func(match string) string {
			return "''" + match + "''"
		})
		logSQL("GenerateCountQuery (Vacuum)", countQuery, args)
		return countQuery, args
	}

	logSQL("GenerateCountQuery", query, args)
	return query, args
}

func (params *QueryParams) buildWhereClauses() ([]string, []any) {
	var whereClauses []string
	var args []any
	var orClauses []string

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

	for field, values := range params.Like {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			var conds []string
			for _, value := range values {
				conds = append(conds, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
				args = append(args, "%"+value+"%")
			}
			whereClauses = append(whereClauses, "("+strings.Join(conds, " OR ")+")")
		}
	}

	for field, values := range params.LikeAnd {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			var conds []string
			for _, value := range values {
				conds = append(conds, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
				args = append(args, "%"+value+"%")
			}
			whereClauses = append(whereClauses, "("+strings.Join(conds, " AND ")+")")
		}
	}

	for field, values := range params.Eq {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			var conds []string
			for _, value := range values {
				conds = append(conds, fmt.Sprintf("%s = ?", columnName))
				args = append(args, value)
			}
			whereClauses = append(whereClauses, "("+strings.Join(conds, " OR ")+")")
		}
	}

	for field, values := range params.EqAnd {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			var conds []string
			for _, value := range values {
				conds = append(conds, fmt.Sprintf("%s = ?", columnName))
				args = append(args, value)
			}
			whereClauses = append(whereClauses, "("+strings.Join(conds, " AND ")+")")
		}
	}

	for field, value := range params.Gte {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s >= ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.Gt {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s > ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.Lte {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s <= ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.Lt {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s < ?", columnName))
			args = append(args, value)
		}
	}

	for field, values := range params.In {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s IN (%s)", columnName, strings.Join(placeholders, ", ")))
		}
	}

	for field, values := range params.NotIn {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			whereClauses = append(whereClauses, fmt.Sprintf("%s NOT IN (%s)", columnName, strings.Join(placeholders, ", ")))
		}
	}

	// OR group
	for field, values := range params.LikeOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			for _, value := range values {
				orClauses = append(orClauses, fmt.Sprintf("%s::TEXT ILIKE ?", columnName))
				args = append(args, "%"+value+"%")
			}
		}
	}

	for field, values := range params.EqOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			for _, value := range values {
				orClauses = append(orClauses, fmt.Sprintf("%s = ?", columnName))
				args = append(args, value)
			}
		}
	}

	for field, value := range params.GteOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s >= ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.GtOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s > ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.LteOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s <= ?", columnName))
			args = append(args, value)
		}
	}

	for field, value := range params.LtOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s < ?", columnName))
			args = append(args, value)
		}
	}

	for field, values := range params.InOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			orClauses = append(orClauses, fmt.Sprintf("%s IN (%s)", columnName, strings.Join(placeholders, ", ")))
		}
	}

	for field, values := range params.NotInOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" && len(values) > 0 {
			placeholders := make([]string, len(values))
			for i := range values {
				placeholders[i] = "?"
				args = append(args, values[i])
			}
			orClauses = append(orClauses, fmt.Sprintf("%s NOT IN (%s)", columnName, strings.Join(placeholders, ", ")))
		}
	}

	for _, field := range params.IsNullOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s IS NULL", columnName))
		}
	}

	for _, field := range params.IsNotNullOr {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			orClauses = append(orClauses, fmt.Sprintf("%s IS NOT NULL", columnName))
		}
	}

	if len(orClauses) > 0 {
		whereClauses = append(whereClauses, "("+strings.Join(orClauses, " OR ")+")")
	}

	for field, values := range params.Between {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN ? AND ?", columnName))
			args = append(args, values[0], values[1])
		}
	}

	for _, field := range params.IsNull {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s IS NULL", columnName))
		}
	}

	for _, field := range params.IsNotNull {
		columnName := getFieldName(field, "json", "paginate", params.Struct)
		if columnName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s IS NOT NULL", columnName))
		}
	}

	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	return whereClauses, args
}

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

func getFieldName(tag, key, keyTarget string, s any) string {
	rt := reflect.TypeOf(s)
	if rt == nil {
		return tag
	}
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return tag
	}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tagValue := strings.Split(field.Tag.Get(key), ",")[0]
		if tagValue == tag {
			paginateTag := field.Tag.Get(keyTarget)
			if paginateTag != "" {
				return paginateTag
			}
			return tagValue
		}
	}
	return tag
}
