package paginate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
)

// PaginatorBuilder provides a fluent interface for building paginated queries
type PaginatorBuilder struct {
	params *QueryParams
	err    error
}

// NewBuilder creates a new PaginatorBuilder with default values from global config
func NewBuilder() *PaginatorBuilder {
	return &PaginatorBuilder{
		params: &QueryParams{
			Page:           1,
			ItemsPerPage:   GetDefaultLimit(), // Use global default
			Vacuum:         false,
			WhereCombining: "AND",
		},
	}
}

// Table sets the main table for the query
func (b *PaginatorBuilder) Table(table string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Table = table
	return b
}

// Schema sets the database schema
func (b *PaginatorBuilder) Schema(schema string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Schema = schema
	return b
}

// Model sets the struct model for the query
func (b *PaginatorBuilder) Model(model any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Struct = model
	return b
}

// Page sets the page number (1-based)
func (b *PaginatorBuilder) Page(page int) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if page < 1 {
		b.err = errors.New("page must be greater than 0")
		return b
	}
	b.params.Page = page
	return b
}

// Limit sets the number of items per page with global max limit validation
func (b *PaginatorBuilder) Limit(limit int) *PaginatorBuilder {
	logger := slog.With("component", "go-paginate-builder")

	if b.err != nil {
		return b
	}

	if limit < 1 {
		b.err = errors.New("limit must be greater than 0")
		logger.Error("Invalid limit value",
			"attempted_value", limit,
			"error", b.err)
		return b
	}

	// Check against global max limit
	maxLimit := GetMaxLimit()
	if limit > maxLimit {
		logger.Warn("Limit exceeds maximum allowed, using max limit",
			"requested_limit", limit,
			"max_limit", maxLimit)
		limit = maxLimit
	}

	b.params.ItemsPerPage = limit
	logger.Debug("Limit set", "limit", limit)
	return b
}

// Select sets the columns to select
func (b *PaginatorBuilder) Select(columns ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Columns = columns
	return b
}

// Search adds a simple search across specified fields
func (b *PaginatorBuilder) Search(term string, fields ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Search = term
	b.params.SearchFields = fields
	return b
}

// LikeOr adds OR search conditions for specific fields
func (b *PaginatorBuilder) LikeOr(field string, values ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.LikeOr == nil {
		b.params.LikeOr = make(map[string][]string)
	}
	b.params.LikeOr[field] = append(b.params.LikeOr[field], values...)
	return b
}

// SearchOr is deprecated, use LikeOr instead
func (b *PaginatorBuilder) SearchOr(field string, values ...string) *PaginatorBuilder {
	return b.LikeOr(field, values...)
}

// LikeAnd adds AND search conditions for specific fields
func (b *PaginatorBuilder) LikeAnd(field string, values ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.LikeAnd == nil {
		b.params.LikeAnd = make(map[string][]string)
	}
	b.params.LikeAnd[field] = append(b.params.LikeAnd[field], values...)
	return b
}

// SearchAnd is deprecated, use LikeAnd instead
func (b *PaginatorBuilder) SearchAnd(field string, values ...string) *PaginatorBuilder {
	return b.LikeAnd(field, values...)
}

// Where adds a custom WHERE clause
func (b *PaginatorBuilder) Where(clause string, args ...any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.WhereClauses = append(b.params.WhereClauses, clause)
	b.params.WhereArgs = append(b.params.WhereArgs, args...)
	return b
}

// WhereEquals adds equality conditions
func (b *PaginatorBuilder) WhereEquals(field string, value any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.EqAnd == nil {
		b.params.EqAnd = make(map[string][]any)
	}
	b.params.EqAnd[field] = append(b.params.EqAnd[field], value)
	return b
}

// EqOr adds OR equality conditions
func (b *PaginatorBuilder) EqOr(field string, values ...any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.EqOr == nil {
		b.params.EqOr = make(map[string][]any)
	}
	b.params.EqOr[field] = append(b.params.EqOr[field], values...)
	return b
}

// EqAnd adds AND equality conditions
func (b *PaginatorBuilder) EqAnd(field string, values ...any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.EqAnd == nil {
		b.params.EqAnd = make(map[string][]any)
	}
	b.params.EqAnd[field] = append(b.params.EqAnd[field], values...)
	return b
}

// WhereEqualsOr is deprecated, use EqOr instead
func (b *PaginatorBuilder) WhereEqualsOr(field string, values ...any) *PaginatorBuilder {
	return b.EqOr(field, values...)
}

// WhereIn adds IN conditions (alias for EqOr)
func (b *PaginatorBuilder) WhereIn(field string, values ...any) *PaginatorBuilder {
	return b.EqOr(field, values...)
}

// WhereGreaterThan adds greater than conditions
func (b *PaginatorBuilder) WhereGreaterThan(field string, value any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Gt == nil {
		b.params.Gt = make(map[string]any)
	}
	b.params.Gt[field] = value
	return b
}

// WhereGreaterThanOrEqual adds greater than or equal conditions
func (b *PaginatorBuilder) WhereGreaterThanOrEqual(field string, value any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Gte == nil {
		b.params.Gte = make(map[string]any)
	}
	b.params.Gte[field] = value
	return b
}

// WhereLessThan adds less than conditions
func (b *PaginatorBuilder) WhereLessThan(field string, value any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Lt == nil {
		b.params.Lt = make(map[string]any)
	}
	b.params.Lt[field] = value
	return b
}

// WhereLessThanOrEqual adds less than or equal conditions
func (b *PaginatorBuilder) WhereLessThanOrEqual(field string, value any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Lte == nil {
		b.params.Lte = make(map[string]any)
	}
	b.params.Lte[field] = value
	return b
}

// WhereBetween adds BETWEEN conditions
func (b *PaginatorBuilder) WhereBetween(field string, min, max any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Between == nil {
		b.params.Between = make(map[string][2]any)
	}
	b.params.Between[field] = [2]any{min, max}
	return b
}

// WhereNotIn adds NOT IN conditions
func (b *PaginatorBuilder) WhereNotIn(field string, values ...any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.NotIn == nil {
		b.params.NotIn = make(map[string][]any)
	}
	b.params.NotIn[field] = values
	return b
}

// WhereIsNull adds IS NULL conditions
func (b *PaginatorBuilder) WhereIsNull(field string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.IsNull = append(b.params.IsNull, field)
	return b
}

// WhereIsNotNull adds IS NOT NULL conditions
func (b *PaginatorBuilder) WhereIsNotNull(field string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.IsNotNull = append(b.params.IsNotNull, field)
	return b
}

// WhereLike adds LIKE conditions
func (b *PaginatorBuilder) WhereLike(field string, values ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.Like == nil {
		b.params.Like = make(map[string][]string)
	}
	b.params.Like[field] = values
	return b
}



// OrderBy adds sorting
func (b *PaginatorBuilder) OrderBy(column string, direction ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.SortColumns = append(b.params.SortColumns, column)
	dir := "ASC"
	if len(direction) > 0 && strings.ToUpper(direction[0]) == "DESC" {
		dir = "DESC"
	}
	b.params.SortDirections = append(b.params.SortDirections, dir)
	return b
}

// OrderByDesc adds descending sorting
func (b *PaginatorBuilder) OrderByDesc(column string) *PaginatorBuilder {
	return b.OrderBy(column, "DESC")
}

// Join adds a JOIN clause
func (b *PaginatorBuilder) Join(joinClause string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Joins = append(b.params.Joins, joinClause)
	return b
}

// LeftJoin adds a LEFT JOIN clause
func (b *PaginatorBuilder) LeftJoin(table, condition string) *PaginatorBuilder {
	return b.Join(fmt.Sprintf("LEFT JOIN %s ON %s", table, condition))
}

// InnerJoin adds an INNER JOIN clause
func (b *PaginatorBuilder) InnerJoin(table, condition string) *PaginatorBuilder {
	return b.Join(fmt.Sprintf("INNER JOIN %s ON %s", table, condition))
}

// RightJoin adds a RIGHT JOIN clause
func (b *PaginatorBuilder) RightJoin(table, condition string) *PaginatorBuilder {
	return b.Join(fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition))
}

// WithoutOffset disables OFFSET in the query (useful for cursor-based pagination)
func (b *PaginatorBuilder) WithoutOffset() *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.NoOffset = true
	return b
}

// WithVacuum enables vacuum mode
func (b *PaginatorBuilder) WithVacuum() *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	b.params.Vacuum = true
	return b
}

// FromJSON populates the builder from a JSON string
func (b *PaginatorBuilder) FromJSON(jsonStr string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		b.err = fmt.Errorf("invalid JSON: %w", err)
		return b
	}

	return b.fromMap(data)
}

// FromStruct populates the builder from a struct using reflection
func (b *PaginatorBuilder) FromStruct(structData any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}

	data, err := structToMap(structData)
	if err != nil {
		b.err = fmt.Errorf("failed to convert struct to map: %w", err)
		return b
	}

	return b.fromMap(data)
}

// FromMap populates the builder from a map
func (b *PaginatorBuilder) FromMap(data map[string]any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	return b.fromMap(data)
}

func (b *PaginatorBuilder) fromMap(data map[string]any) *PaginatorBuilder {
	// Handle page
	if page, ok := data["page"]; ok {
		if pageInt, err := toInt(page); err == nil {
			b.Page(pageInt)
		}
	}

	// Handle limit/items_per_page
	if limit, ok := data["limit"]; ok {
		if limitInt, err := toInt(limit); err == nil {
			b.Limit(limitInt)
		}
	} else if itemsPerPage, ok := data["items_per_page"]; ok {
		if limitInt, err := toInt(itemsPerPage); err == nil {
			b.Limit(limitInt)
		}
	}

	// Handle search
	if search, ok := data["search"]; ok {
		if searchStr, ok := search.(string); ok && searchStr != "" {
			if searchFields, ok := data["search_fields"]; ok {
				if fields := toStringSlice(searchFields); len(fields) > 0 {
					b.Search(searchStr, fields...)
				}
			} else {
				// If no search_fields specified, just set the search term
				b.params.Search = searchStr
			}
		}
	}

	// Handle likeor (both "likeor" and "like_or" for struct compatibility)
	for _, key := range []string{"likeor", "like_or"} {
		if likeOr, ok := data[key]; ok {
			if likeOrMap, ok := likeOr.(map[string]any); ok {
				for field, values := range likeOrMap {
					if valueSlice := toStringSlice(values); len(valueSlice) > 0 {
						b.LikeOr(field, valueSlice...)
					}
				}
			} else if likeOrMapStr, ok := likeOr.(map[string][]string); ok {
				// Handle direct map[string][]string from struct conversion
				for field, values := range likeOrMapStr {
					if len(values) > 0 {
						b.LikeOr(field, values...)
					}
				}
			}
			break // Only process the first match
		}
	}

	// Handle likeand (both "likeand" and "like_and" for struct compatibility)
	for _, key := range []string{"likeand", "like_and"} {
		if likeAnd, ok := data[key]; ok {
			if likeAndMap, ok := likeAnd.(map[string]any); ok {
				for field, values := range likeAndMap {
					if valueSlice := toStringSlice(values); len(valueSlice) > 0 {
						b.LikeAnd(field, valueSlice...)
					}
				}
			} else if likeAndMapStr, ok := likeAnd.(map[string][]string); ok {
				// Handle direct map[string][]string from struct conversion
				for field, values := range likeAndMapStr {
					if len(values) > 0 {
						b.LikeAnd(field, values...)
					}
				}
			}
			break // Only process the first match
		}
	}

	// Handle eqor (both "eqor" and "eq_or" for struct compatibility)
	for _, key := range []string{"eqor", "eq_or"} {
		if eqOr, ok := data[key]; ok {
			if eqOrMap, ok := eqOr.(map[string]any); ok {
				for field, values := range eqOrMap {
					if valueSlice := toInterfaceSlice(values); len(valueSlice) > 0 {
						b.EqOr(field, valueSlice...)
					}
				}
			} else if eqOrMapSlice, ok := eqOr.(map[string][]any); ok {
				// Handle direct map[string][]any from struct conversion
				for field, values := range eqOrMapSlice {
					if len(values) > 0 {
						b.EqOr(field, values...)
					}
				}
			}
			break // Only process the first match
		}
	}

	// Handle eqand (both "eqand" and "eq_and" for struct compatibility)
	for _, key := range []string{"eqand", "eq_and"} {
		if eqAnd, ok := data[key]; ok {
			if eqAndMap, ok := eqAnd.(map[string]any); ok {
				for field, values := range eqAndMap {
					if valueSlice := toInterfaceSlice(values); len(valueSlice) > 0 {
						b.EqAnd(field, valueSlice...)
					}
				}
			} else if eqAndMapSlice, ok := eqAnd.(map[string][]any); ok {
				// Handle direct map[string][]any from struct conversion
				for field, values := range eqAndMapSlice {
					if len(values) > 0 {
						b.EqAnd(field, values...)
					}
				}
			}
			break // Only process the first match
		}
	}

	// Handle comparison operators
	// Handle gte (both "gte" and "gte_" for struct compatibility)
	for _, key := range []string{"gte", "gte_"} {
		if gte, ok := data[key]; ok {
			if gteMap, ok := gte.(map[string]any); ok {
				for field, value := range gteMap {
					b.WhereGreaterThanOrEqual(field, value)
				}
			}
			break // Only process the first match
		}
	}

	// Handle gt (both "gt" and "gt_" for struct compatibility)
	for _, key := range []string{"gt", "gt_"} {
		if gt, ok := data[key]; ok {
			if gtMap, ok := gt.(map[string]any); ok {
				for field, value := range gtMap {
					b.WhereGreaterThan(field, value)
				}
			}
			break // Only process the first match
		}
	}

	// Handle lte (both "lte" and "lte_" for struct compatibility)
	for _, key := range []string{"lte", "lte_"} {
		if lte, ok := data[key]; ok {
			if lteMap, ok := lte.(map[string]any); ok {
				for field, value := range lteMap {
					b.WhereLessThanOrEqual(field, value)
				}
			}
			break // Only process the first match
		}
	}

	// Handle lt (both "lt" and "lt_" for struct compatibility)
	for _, key := range []string{"lt", "lt_"} {
		if lt, ok := data[key]; ok {
			if ltMap, ok := lt.(map[string]any); ok {
				for field, value := range ltMap {
					b.WhereLessThan(field, value)
				}
			}
			break // Only process the first match
		}
	}

	// Handle sorting - supports both single sort field and multiple sort fields
	// New sort pattern takes priority over legacy sort_columns/sort_directions
	if sort, ok := data["sort"]; ok {
		if sortSlice := toStringSlice(sort); len(sortSlice) > 0 {
			for _, sortField := range sortSlice {
				if strings.HasPrefix(sortField, "-") {
					b.OrderByDesc(strings.TrimPrefix(sortField, "-"))
				} else {
					b.OrderBy(sortField)
				}
			}
		}
	} else {
		// Handle legacy sort_columns and sort_directions for backward compatibility
		// Only process if new sort pattern is not present
		if sortColumns, ok := data["sort_columns"]; ok {
			if sortColumnsSlice := toStringSlice(sortColumns); len(sortColumnsSlice) > 0 {
				var sortDirectionsSlice []string
				if sortDirections, ok := data["sort_directions"]; ok {
					sortDirectionsSlice = toStringSlice(sortDirections)
				}

				// Apply sorting for each column
				for i, column := range sortColumnsSlice {
					direction := "ASC" // default direction
					if i < len(sortDirectionsSlice) {
						direction = sortDirectionsSlice[i]
					}

					if strings.ToUpper(direction) == "DESC" {
						b.OrderByDesc(column)
					} else {
						b.OrderBy(column)
					}
				}
			}
		}
	}

	return b
}

// Build creates the QueryParams and validates it
func (b *PaginatorBuilder) Build() (*QueryParams, error) {
	if b.err != nil {
		return nil, b.err
	}

	// Validation
	if b.params.Table == "" {
		return nil, errors.New("table is required")
	}

	if b.params.Struct == nil {
		return nil, errors.New("model struct is required")
	}

	return b.params, nil
}

// BuildSQL generates the SQL query directly
func (b *PaginatorBuilder) BuildSQL() (string, []any, error) {
	params, err := b.Build()
	if err != nil {
		return "", nil, err
	}

	sql, args := params.GenerateSQL()
	
	// Log SQL if debug mode is enabled
	logSQL("BuildSQL", sql, args)
	
	return sql, args, nil
}

// BuildCountSQL generates the count SQL query directly
func (b *PaginatorBuilder) BuildCountSQL() (string, []any, error) {
	params, err := b.Build()
	if err != nil {
		return "", nil, err
	}

	sql, args := params.GenerateCountQuery()
	
	// Log SQL if debug mode is enabled
	logSQL("BuildCountSQL", sql, args)
	
	return sql, args, nil
}

// Helper functions

func toInt(value any) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

func toStringSlice(value any) []string {
	switch v := value.(type) {
	case []string:
		return v
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	case string:
		return []string{v}
	default:
		return nil
	}
}

func toInterfaceSlice(value any) []any {
	switch v := value.(type) {
	case []any:
		return v
	case []string:
		result := make([]any, len(v))
		for i, str := range v {
			result[i] = str
		}
		return result
	default:
		// Try to use reflection for other slice types
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Slice {
			result := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = rv.Index(i).Interface()
			}
			return result
		}
		return []any{value}
	}
}

// structToMap converts a struct to a map[string]any using reflection
func structToMap(structData any) (map[string]any, error) {
	result := make(map[string]any)

	// Handle nil input
	if structData == nil {
		return result, nil
	}

	v := reflect.ValueOf(structData)
	t := reflect.TypeOf(structData)

	// Handle pointer to struct
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return result, nil
		}
		v = v.Elem()
		t = t.Elem()
	}

	// Ensure we have a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct or pointer to struct, got %T", structData)
	}

	// Iterate through struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get field name from json tag or field name
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			// Handle json tag with options like "field_name,omitempty"
			tagParts := strings.Split(jsonTag, ",")
			if tagParts[0] != "" && tagParts[0] != "-" {
				fieldName = tagParts[0]
			}
			// Skip fields marked with json:"-"
			if tagParts[0] == "-" {
				continue
			}
		}

		// Convert field name to snake_case for consistency
		fieldName = toSnakeCase(fieldName)

		// Get field value
		fieldValue := field.Interface()

		// Handle zero values based on field type
		if isZeroValue(field) {
			continue
		}

		result[fieldName] = fieldValue
	}

	return result, nil
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// isZeroValue checks if a reflect.Value is a zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}
