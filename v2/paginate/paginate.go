package paginate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// paginQueryParams contém os parâmetros para a consulta paginada
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
	Table          string
	Struct         interface{}
}

// Option é uma função que configura opções em paginQueryParams
type Option func(*paginQueryParams)

func WithStruct(s interface{}) Option {
	return func(params *paginQueryParams) {
		params.Struct = s
	}
}

func WithTable(table string) Option {
	return func(params *paginQueryParams) {
		params.Table = table
	}
}

// WithPage configura o campo Page
func WithPage(page int) Option {
	return func(params *paginQueryParams) {
		params.Page = page
	}
}

// WithItemsPerPage configura o campo ItemsPerPage
func WithItemsPerPage(itemsPerPage int) Option {
	return func(params *paginQueryParams) {
		params.ItemsPerPage = itemsPerPage
	}
}

// WithSearch configura o campo Search
func WithSearch(search string) Option {
	return func(params *paginQueryParams) {
		params.Search = search
	}
}

// WithSearchFields configura o campo SearchFields
func WithSearchFields(searchFields []string) Option {
	return func(params *paginQueryParams) {
		params.SearchFields = searchFields
	}
}

// WithVacuum configura o campo Vacuum
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

// WithSort configura os campos SortColumns e SortDirections
func WithSort(sortColumns []string, sortDirections []string) Option {
	return func(params *paginQueryParams) {
		params.SortColumns = sortColumns
		params.SortDirections = sortDirections
	}
}

// WithJoin configura o campo Joins
func WithJoin(join string) Option {
	return func(params *paginQueryParams) {
		params.Joins = append(params.Joins, join)
	}
}

// WithWhereCombining especifica o operador de combinação para as cláusulas WHERE
func WithWhereCombining(combining string) Option {
	return func(params *paginQueryParams) {
		params.WhereCombining = combining
	}
}

// WithWhereClause adiciona uma cláusula WHERE
func WithWhereClause(clause string, args ...interface{}) Option {
	return func(params *paginQueryParams) {
		params.WhereClauses = append(params.WhereClauses, clause)
		params.WhereArgs = append(params.WhereArgs, args...)
	}
}

func PaginQuery(options ...Option) (*paginQueryParams, error) {
	// Valores padrão
	params := &paginQueryParams{
		Page:           1,
		ItemsPerPage:   10,
		WhereCombining: "AND", // Combinação padrão é "AND"
	}

	// Aplicar opções
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
	// Inicializa uma lista vazia de cláusulas SQL
	clauses := []string{}
	args := []interface{}{}

	// Função auxiliar para obter o próximo número de argumento
	nextArg := func() int {
		argNum := len(args) + 1
		args = append(args, nil) // Adicione um espaço reservado para o próximo argumento
		return argNum
	}

	// Cláusula SELECT com colunas personalizadas
	selectClause := "SELECT "
	if len(params.Columns) > 0 {
		// columnNames := []string{}
		// for _, column := range params.Columns {
		// 	columnName := getFieldName(column, "json", "paginate", params.Struct)
		// 	if columnName != "" {
		// 		columnNames = append(columnNames, columnName)
		// 	}
		// }
		// if len(columnNames) > 0 {
		selectClause += strings.Join(params.Columns, ", ")
	} else {
		selectClause += "*"
	}
	clauses = append(clauses, selectClause)

	// Cláusula FROM com tabela principal
	clauses = append(clauses, fmt.Sprintf("FROM %s", params.Table))

	// Cláusulas JOIN personalizadas
	if len(params.Joins) > 0 {
		joinClause := strings.Join(params.Joins, " ")
		clauses = append(clauses, joinClause)
	}

	// Cláusula WHERE para pesquisa
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
			searchClause := fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR "))
			whereClauses = append(whereClauses, searchClause)
		}
	}

	// Adicionar cláusulas WHERE personalizadas
	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	if len(whereClauses) > 0 {
		whereClause := strings.Join(whereClauses, " AND ")
		clauses = append(clauses, "WHERE "+whereClause)
	}

	// Cláusula ORDER BY com suporte a múltiplas colunas e direções de ordenação
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
			sortClause := fmt.Sprintf("ORDER BY %s", strings.Join(sortClauses, ", "))
			clauses = append(clauses, sortClause)
		}
	}

	// Cláusula LIMIT e OFFSET para paginação
	offset := (params.Page - 1) * params.ItemsPerPage
	clauses = append(clauses, "LIMIT $"+fmt.Sprint(nextArg())+" OFFSET $"+fmt.Sprint(nextArg()))
	args[len(args)-2] = params.ItemsPerPage
	args[len(args)-1] = offset

	replacePlaceholders := func(query string, args []interface{}) (string, []interface{}) {
		// Use um contador para acompanhar a posição do próximo argumento
		// argCount := 1
		// Encontre o último número de argumento disponível antes do primeiro ?
		lastArg := 0
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				break
			} else if query[i] == '$' {
				// Se encontrarmos um placeholder existente, atualizamos o contador
				lastArg, _ = strconv.Atoi(string(query[i+1]))
			}
		}
		// Substitua todos os ? pelos placeholders $n, onde n é o último número encontrado e incrementado
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				query = query[:i] + "$" + strconv.Itoa(lastArg+1) + query[i+1:]
				lastArg++
			}
		}
		return query, args
	}

	// Junte todas as cláusulas em uma única consulta SQL
	query := strings.Join(clauses, " ")

	// Substitua os placeholders ? pelos placeholders $n
	query, args = replacePlaceholders(query, args)

	return query, args
}

func GenerateCountQuery(params *paginQueryParams) (string, []interface{}) {
	// Inicializa uma lista vazia de cláusulas SQL para a contagem
	clauses := []string{}
	args := []interface{}{}

	// Função auxiliar para obter o próximo número de argumento
	nextArg := func() int {
		argNum := len(args) + 1
		args = append(args, nil) // Adicione um espaço reservado para o próximo argumento
		return argNum
	}

	// Cláusula SELECT para contagem
	countSelectClause := "SELECT COUNT(id)"
	clauses = append(clauses, countSelectClause)

	// Cláusula FROM com tabela principal
	clauses = append(clauses, fmt.Sprintf("FROM %s", params.Table))

	// Cláusulas JOIN personalizadas
	if len(params.Joins) > 0 {
		joinClause := strings.Join(params.Joins, " ")
		clauses = append(clauses, joinClause)
	}

	// Cláusula WHERE para pesquisa
	whereClauses := []string{}

	if params.Search != "" && len(params.SearchFields) > 0 {
		searchConditions := []string{}
		for _, field := range params.SearchFields {
			columnName := getFieldName(field, "json", "paginate", params.Struct)
			if columnName != "" {
				searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE $%d", columnName, nextArg()))
				args[len(args)-1] = "%" + params.Search + "%"
			}
		}
		if len(searchConditions) > 0 {
			searchClause := fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR "))
			whereClauses = append(whereClauses, searchClause)
		}
	}

	// Adicionar cláusulas WHERE personalizadas
	if len(params.WhereClauses) > 0 {
		whereClauses = append(whereClauses, strings.Join(params.WhereClauses, fmt.Sprintf(" %s ", params.WhereCombining)))
		args = append(args, params.WhereArgs...)
	}

	if len(whereClauses) > 0 {
		whereClause := strings.Join(whereClauses, " AND ")
		clauses = append(clauses, "WHERE "+whereClause)
	}

	// Verifica se VACUUM deve ser aplicado
	if params.Vacuum {
		countQuery := strings.Join(clauses, " ")
		countQuery = "SELECT count_estimate('" + countQuery + "');"
		countQuery = strings.Replace(countQuery, "COUNT(id)", "1", -1)
		return countQuery, args
	}

	replacePlaceholders := func(query string, args []interface{}) (string, []interface{}) {
		// Use um contador para acompanhar a posição do próximo argumento
		// argCount := 1
		// Encontre o último número de argumento disponível antes do primeiro ?
		lastArg := 0
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				break
			} else if query[i] == '$' {
				// Se encontrarmos um placeholder existente, atualizamos o contador
				lastArg, _ = strconv.Atoi(string(query[i+1]))
			}
		}
		// Substitua todos os ? pelos placeholders $n, onde n é o último número encontrado e incrementado
		for i := 0; i < len(query); i++ {
			if query[i] == '?' {
				query = query[:i] + "$" + strconv.Itoa(lastArg+1) + query[i+1:]
				lastArg++
			}
		}
		return query, args
	}

	// Junte todas as cláusulas em uma única consulta SQL
	query := strings.Join(clauses, " ")

	// Substitua os placeholders ? pelos placeholders $n
	query, args = replacePlaceholders(query, args)

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
