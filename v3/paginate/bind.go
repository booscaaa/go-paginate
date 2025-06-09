package paginate

import (
	"fmt"
	"log/slog"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// PaginationParams represents pagination parameters that can be extracted from query params
type PaginationParams struct {
	Page           int                 `query:"page"`
	Limit          int                 `query:"limit"`
	ItemsPerPage   int                 `query:"items_per_page"`
	Search         string              `query:"search"`
	SearchFields   []string            `query:"search_fields"`
	Sort           []string            `query:"sort"`
	SortColumns    []string            `query:"sort_columns"`
	SortDirections []string            `query:"sort_directions"`
	Columns        []string            `query:"columns"`
	Vacuum         bool                `query:"vacuum"`
	NoOffset       bool                `query:"no_offset"`
	Like           map[string][]string `query:"like"`
	LikeOr         map[string][]string `query:"likeor"`
	LikeAnd        map[string][]string `query:"likeand"`
	Eq             map[string][]any    `query:"eq"`
	EqOr           map[string][]any    `query:"eqor"`
	EqAnd          map[string][]any    `query:"eqand"`
	Gte            map[string]any      `query:"gte"`
	Gt             map[string]any      `query:"gt"`
	Lte            map[string]any      `query:"lte"`
	Lt             map[string]any      `query:"lt"`
	In             map[string][]any    `query:"in"`
	NotIn          map[string][]any    `query:"notin"`
	Between        map[string][2]any   `query:"between"`
	IsNull         []string            `query:"isnull"`
	IsNotNull      []string            `query:"isnotnull"`
}

// BindQueryParams binds url.Values to a pagination struct
func BindQueryParams(queryParams url.Values, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to struct")
	}

	v = v.Elem()
	t := v.Type()

	// Inicializar maps se necessário
	initializeMaps(v, t)

	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)
		queryTag := fieldType.Tag.Get("query")

		if queryTag == "" || !field.CanSet() {
			continue
		}

		// Tratar campos especiais com sintaxe de array
		if isMapField(fieldType.Type) {
			bindMapField(queryParams, field, queryTag)
			continue
		}

		// Tratar campos normais
		values, exists := queryParams[queryTag]
		if !exists || len(values) == 0 {
			continue
		}

		if err := setFieldValue(field, values); err != nil {
			return fmt.Errorf("error setting field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

// initializeMaps inicializa os campos de map se eles forem nil
func initializeMaps(v reflect.Value, t reflect.Type) {
	for i := range v.NumField() {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanSet() {
			continue
		}

		switch fieldType.Type.Kind() {
		case reflect.Map:
			if field.IsNil() {
				field.Set(reflect.MakeMap(fieldType.Type))
			}
		}
	}
}

// isMapField verifica se o campo é um map
func isMapField(t reflect.Type) bool {
	return t.Kind() == reflect.Map
}

// bindMapField faz bind de parâmetros com sintaxe de array para campos de map
func bindMapField(queryParams url.Values, field reflect.Value, queryTag string) {
	if !field.CanSet() || field.Kind() != reflect.Map {
		return
	}

	mapType := field.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()

	// Procurar por parâmetros com formato: queryTag[key]=value
	prefix := queryTag + "["
	for paramName, values := range queryParams {
		if strings.HasPrefix(paramName, prefix) && strings.HasSuffix(paramName, "]") {
			// Extrair a chave do parâmetro
			key := paramName[len(prefix) : len(paramName)-1]
			if key == "" {
				continue
			}

			// Converter a chave para o tipo correto
			keyValue := reflect.ValueOf(key)
			if keyType.Kind() != reflect.String {
				continue // Por enquanto, só suportamos chaves string
			}

			// Converter os valores para o tipo correto
			var mapValue reflect.Value
			switch valueType.Kind() {
			case reflect.Array:
				// Para arrays fixos como [2]any (usado em Between)
				if valueType.Len() == 2 && len(values) >= 2 {
					arrayValue := reflect.New(valueType).Elem()
					for i := 0; i < 2 && i < len(values); i++ {
						value := values[i]
						var elem reflect.Value
						if valueType.Elem().Kind() == reflect.Interface {
							// Try to convert to number first, then boolean, otherwise keep as string
							if intVal, err := strconv.Atoi(value); err == nil {
								elem = reflect.ValueOf(intVal)
							} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
								elem = reflect.ValueOf(floatVal)
							} else if boolVal, err := strconv.ParseBool(value); err == nil {
								elem = reflect.ValueOf(boolVal)
							} else {
								elem = reflect.ValueOf(value)
							}
						} else {
							elem = reflect.ValueOf(value)
						}
						arrayValue.Index(i).Set(elem)
					}
					mapValue = arrayValue
				}
			case reflect.Slice:
				// Para []string ou []any
				sliceType := valueType.Elem()
				slice := reflect.MakeSlice(valueType, 0, len(values))
				for _, value := range values {
					var elem reflect.Value
					if sliceType.Kind() == reflect.Interface {
						// Try to convert to number first, then boolean, otherwise keep as string
						if intVal, err := strconv.Atoi(value); err == nil {
							elem = reflect.ValueOf(intVal)
						} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
							elem = reflect.ValueOf(floatVal)
						} else if boolVal, err := strconv.ParseBool(value); err == nil {
							elem = reflect.ValueOf(boolVal)
						} else {
							elem = reflect.ValueOf(value)
						}
					} else {
						elem = reflect.ValueOf(value)
					}
					slice = reflect.Append(slice, elem)
				}
				mapValue = slice
			case reflect.Interface:
				// Para any, usar o primeiro valor
				if len(values) > 0 {
					value := values[0]
					// Try to convert to number first, then boolean, otherwise keep as string
					if intVal, err := strconv.Atoi(value); err == nil {
						mapValue = reflect.ValueOf(intVal)
					} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
						mapValue = reflect.ValueOf(floatVal)
					} else if boolVal, err := strconv.ParseBool(value); err == nil {
						mapValue = reflect.ValueOf(boolVal)
					} else {
						mapValue = reflect.ValueOf(value)
					}
				}
			default:
				continue // Tipo não suportado
			}

			if mapValue.IsValid() {
				field.SetMapIndex(keyValue, mapValue)
			}
		}
	}
}

// setFieldValue define o valor de um campo baseado nos valores dos query params
func setFieldValue(field reflect.Value, values []string) error {
	if len(values) == 0 {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(values[0])
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(values[0], 10, 64); err == nil {
			field.SetInt(intVal)
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(values[0]); err == nil {
			field.SetBool(boolVal)
		}
	case reflect.Slice:
		// Para slices, usar todos os valores ou dividir por vírgula se for um único valor
		var finalValues []string
		if len(values) == 1 && strings.Contains(values[0], ",") {
			finalValues = strings.Split(values[0], ",")
			// Remover espaços em branco
			for i, v := range finalValues {
				finalValues[i] = strings.TrimSpace(v)
			}
		} else {
			finalValues = values
		}

		slice := reflect.MakeSlice(field.Type(), len(finalValues), len(finalValues))
		for i, value := range finalValues {
			slice.Index(i).SetString(value)
		}
		field.Set(slice)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// BindQueryParamsToStruct é uma função de conveniência que cria uma nova instância de PaginationParams
// e faz bind dos query params para ela
func BindQueryParamsToStruct(queryParams url.Values) (*PaginationParams, error) {
	params := &PaginationParams{
		Page:         1,  // valor padrão
		Limit:        10, // valor padrão
		ItemsPerPage: 10, // valor padrão
	}

	err := BindQueryParams(queryParams, params)
	if err != nil {
		return nil, err
	}

	// If ItemsPerPage was set but Limit wasn't, use ItemsPerPage as Limit
	if params.ItemsPerPage != 10 && params.Limit == 10 {
		params.Limit = params.ItemsPerPage
	}

	return params, nil
}

// BindQueryStringToStruct faz bind de uma query string para PaginationParams
func BindQueryStringToStruct(queryString string) (*PaginationParams, error) {
	queryParams, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("error parsing query string: %w", err)
	}

	return BindQueryParamsToStruct(queryParams)
}

// NewPaginationParams cria uma nova instância com valores padrão globais
func NewPaginationParams() *PaginationParams {
	return &PaginationParams{
		Page:      1,
		Limit:     GetDefaultLimit(), // Use global default
		Like:      make(map[string][]string),
		LikeOr:    make(map[string][]string),
		LikeAnd:   make(map[string][]string),
		Eq:        make(map[string][]any),
		EqOr:      make(map[string][]any),
		EqAnd:     make(map[string][]any),
		Gte:       make(map[string]any),
		Gt:        make(map[string]any),
		Lte:       make(map[string]any),
		Lt:        make(map[string]any),
		In:        make(map[string][]any),
		NotIn:     make(map[string][]any),
		Between:   make(map[string][2]any),
		IsNull:    make([]string, 0),
		IsNotNull: make([]string, 0),
	}
}

// setDefaultValues define valores padrão usando configuração global
func setDefaultValues(params *PaginationParams) {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = GetDefaultLimit() // Use global default
	}

	// Apply global max limit
	maxLimit := GetMaxLimit()
	if params.Limit > maxLimit {
		logger := slog.With("component", "go-paginate-bind")
		logger.Warn("Limit exceeds maximum, applying global max limit",
			"requested_limit", params.Limit,
			"max_limit", maxLimit)
		params.Limit = maxLimit
	}

	// If ItemsPerPage was set but Limit wasn't, use ItemsPerPage as Limit
	if params.ItemsPerPage != GetDefaultLimit() && params.Limit == GetDefaultLimit() {
		params.Limit = params.ItemsPerPage
		// Apply max limit validation again
		if params.Limit > maxLimit {
			params.Limit = maxLimit
		}
	}
}
