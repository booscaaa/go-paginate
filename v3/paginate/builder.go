package paginate

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// PaginatorBuilder provides a fluent interface for building paginated queries
type PaginatorBuilder struct {
	params *QueryParams
	err    error
}

// NewBuilder creates a new PaginatorBuilder instance
func NewBuilder() *PaginatorBuilder {
	return &PaginatorBuilder{
		params: &QueryParams{
			Page:           1,
			ItemsPerPage:   10,
			WhereCombining: "AND",
			NoOffset:       false,
			SearchOr:       make(map[string][]string),
			SearchAnd:      make(map[string][]string),
			EqualsOr:       make(map[string][]any),
			EqualsAnd:      make(map[string][]any),
			Gte:            make(map[string]any),
			Gt:             make(map[string]any),
			Lte:            make(map[string]any),
			Lt:             make(map[string]any),
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

// Limit sets the number of items per page
func (b *PaginatorBuilder) Limit(limit int) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if limit < 1 {
		b.err = errors.New("limit must be greater than 0")
		return b
	}
	b.params.ItemsPerPage = limit
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

// SearchOr adds OR search conditions for specific fields
func (b *PaginatorBuilder) SearchOr(field string, values ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.SearchOr == nil {
		b.params.SearchOr = make(map[string][]string)
	}
	b.params.SearchOr[field] = append(b.params.SearchOr[field], values...)
	return b
}

// SearchAnd adds AND search conditions for specific fields
func (b *PaginatorBuilder) SearchAnd(field string, values ...string) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.SearchAnd == nil {
		b.params.SearchAnd = make(map[string][]string)
	}
	b.params.SearchAnd[field] = append(b.params.SearchAnd[field], values...)
	return b
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
	if b.params.EqualsAnd == nil {
		b.params.EqualsAnd = make(map[string][]any)
	}
	b.params.EqualsAnd[field] = append(b.params.EqualsAnd[field], value)
	return b
}

// WhereEqualsOr adds OR equality conditions
func (b *PaginatorBuilder) WhereEqualsOr(field string, values ...any) *PaginatorBuilder {
	if b.err != nil {
		return b
	}
	if b.params.EqualsOr == nil {
		b.params.EqualsOr = make(map[string][]any)
	}
	b.params.EqualsOr[field] = append(b.params.EqualsOr[field], values...)
	return b
}

// WhereIn adds IN conditions (alias for WhereEqualsOr)
func (b *PaginatorBuilder) WhereIn(field string, values ...any) *PaginatorBuilder {
	return b.WhereEqualsOr(field, values...)
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
	return b.WhereGreaterThanOrEqual(field, min).WhereLessThanOrEqual(field, max)
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

	// Handle search_or
	if searchOr, ok := data["search_or"]; ok {
		if searchOrMap, ok := searchOr.(map[string]any); ok {
			for field, values := range searchOrMap {
				if valueSlice := toStringSlice(values); len(valueSlice) > 0 {
					b.SearchOr(field, valueSlice...)
				}
			}
		} else if searchOrMapStr, ok := searchOr.(map[string][]string); ok {
			// Handle direct map[string][]string from struct conversion
			for field, values := range searchOrMapStr {
				if len(values) > 0 {
					b.SearchOr(field, values...)
				}
			}
		}
	}

	// Handle search_and
	if searchAnd, ok := data["search_and"]; ok {
		if searchAndMap, ok := searchAnd.(map[string]any); ok {
			for field, values := range searchAndMap {
				if valueSlice := toStringSlice(values); len(valueSlice) > 0 {
					b.SearchAnd(field, valueSlice...)
				}
			}
		} else if searchAndMapStr, ok := searchAnd.(map[string][]string); ok {
			// Handle direct map[string][]string from struct conversion
			for field, values := range searchAndMapStr {
				if len(values) > 0 {
					b.SearchAnd(field, values...)
				}
			}
		}
	}

	// Handle equals_or
	if equalsOr, ok := data["equals_or"]; ok {
		if equalsOrMap, ok := equalsOr.(map[string]any); ok {
			for field, values := range equalsOrMap {
				if valueSlice := toInterfaceSlice(values); len(valueSlice) > 0 {
					b.WhereEqualsOr(field, valueSlice...)
				}
			}
		} else if equalsOrMapSlice, ok := equalsOr.(map[string][]any); ok {
			// Handle direct map[string][]any from struct conversion
			for field, values := range equalsOrMapSlice {
				if len(values) > 0 {
					b.WhereEqualsOr(field, values...)
				}
			}
		}
	}

	// Handle equals_and
	if equalsAnd, ok := data["equals_and"]; ok {
		if equalsAndMap, ok := equalsAnd.(map[string]any); ok {
			for field, values := range equalsAndMap {
				if valueSlice := toInterfaceSlice(values); len(valueSlice) > 0 {
					for _, value := range valueSlice {
						b.WhereEquals(field, value)
					}
				}
			}
		} else if equalsAndMapSlice, ok := equalsAnd.(map[string][]any); ok {
			// Handle direct map[string][]any from struct conversion
			for field, values := range equalsAndMapSlice {
				if len(values) > 0 {
					for _, value := range values {
						b.WhereEquals(field, value)
					}
				}
			}
		}
	}

	// Handle comparison operators
	if gte, ok := data["gte"]; ok {
		if gteMap, ok := gte.(map[string]any); ok {
			for field, value := range gteMap {
				b.WhereGreaterThanOrEqual(field, value)
			}
		}
	}

	if gt, ok := data["gt"]; ok {
		if gtMap, ok := gt.(map[string]any); ok {
			for field, value := range gtMap {
				b.WhereGreaterThan(field, value)
			}
		}
	}

	if lte, ok := data["lte"]; ok {
		if lteMap, ok := lte.(map[string]any); ok {
			for field, value := range lteMap {
				b.WhereLessThanOrEqual(field, value)
			}
		}
	}

	if lt, ok := data["lt"]; ok {
		if ltMap, ok := lt.(map[string]any); ok {
			for field, value := range ltMap {
				b.WhereLessThan(field, value)
			}
		}
	}

	// Handle sorting
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
	return sql, args, nil
}

// BuildCountSQL generates the count SQL query directly
func (b *PaginatorBuilder) BuildCountSQL() (string, []any, error) {
	params, err := b.Build()
	if err != nil {
		return "", nil, err
	}

	sql, args := params.GenerateCountQuery()
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
