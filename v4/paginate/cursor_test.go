package paginate

import (
	"net/url"
	"strings"
	"testing"
	"time"
)

// cursorTestModel is a simple struct used across cursor tests.
type cursorTestModel struct {
	ID        int       `json:"id"         paginate:"users.id"`
	Name      string    `json:"name"       paginate:"users.name"`
	Status    string    `json:"status"     paginate:"users.status"`
	CreatedAt time.Time `json:"created_at" paginate:"users.created_at"`
}

// ---------------------------------------------------------------------------
// EncodeCursor / DecodeCursor (single-column backward compat)
// ---------------------------------------------------------------------------

func TestEncodeCursorDecodeCursor_RoundTrip(t *testing.T) {
	col, val, dir := "id", 42, "after"
	token := EncodeCursor(col, val, dir)

	gotCol, gotVal, gotDir, err := DecodeCursor(token)
	if err != nil {
		t.Fatalf("DecodeCursor returned error: %v", err)
	}
	if gotCol != col {
		t.Errorf("column: want %q, got %q", col, gotCol)
	}
	if gotVal.(float64) != float64(val) {
		t.Errorf("value: want %v, got %v", val, gotVal)
	}
	if gotDir != dir {
		t.Errorf("direction: want %q, got %q", dir, gotDir)
	}
}

func TestDecodeCursor_InvalidToken(t *testing.T) {
	_, _, _, err := DecodeCursor("not-valid-base64!!")
	if err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

// ---------------------------------------------------------------------------
// After / Before — single-column backward compat SQL generation
// ---------------------------------------------------------------------------

func TestAfter_GeneratesGreaterThanClause(t *testing.T) {
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		CursorColumn:    "id",
		CursorValue:     42,
		CursorDirection: "after",
		NoOffset:        true,
		SortColumns:     []string{"id"},
		SortDirections:  []string{"ASC"},
	}

	query, args := params.GenerateSQL()

	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected WHERE users.id > $N, got: %s", query)
	}
	if !strings.Contains(query, "LIMIT $") {
		t.Errorf("expected LIMIT clause, got: %s", query)
	}
	if strings.Contains(query, "OFFSET") {
		t.Errorf("expected no OFFSET for cursor pagination, got: %s", query)
	}
	found := false
	for _, a := range args {
		if a == 42 {
			found = true
		}
	}
	if !found {
		t.Errorf("cursor value 42 not found in args: %v", args)
	}
}

func TestBefore_GeneratesLessThanClause(t *testing.T) {
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		CursorColumn:    "id",
		CursorValue:     10,
		CursorDirection: "before",
		NoOffset:        true,
		SortColumns:     []string{"id"},
		SortDirections:  []string{"DESC"},
	}

	query, _ := params.GenerateSQL()

	if !strings.Contains(query, "users.id < $") {
		t.Errorf("expected WHERE users.id < $N, got: %s", query)
	}
}

// ---------------------------------------------------------------------------
// Builder After / Before methods
// ---------------------------------------------------------------------------

func TestBuilder_After(t *testing.T) {
	query, args, err := NewBuilder().
		Table("users").
		Model(&cursorTestModel{}).
		Limit(5).
		After("id", 100).
		OrderBy("id").
		BuildSQL()

	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected cursor WHERE clause, got: %s", query)
	}
	if strings.Contains(query, "OFFSET") {
		t.Errorf("expected no OFFSET, got: %s", query)
	}
	found := false
	for _, a := range args {
		if a == 100 {
			found = true
		}
	}
	if !found {
		t.Errorf("cursor value 100 not in args: %v", args)
	}
}

func TestBuilder_Before(t *testing.T) {
	query, _, err := NewBuilder().
		Table("users").
		Model(&cursorTestModel{}).
		Limit(5).
		Before("id", 50).
		OrderByDesc("id").
		BuildSQL()

	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}
	if !strings.Contains(query, "users.id < $") {
		t.Errorf("expected cursor WHERE clause, got: %s", query)
	}
}

// ---------------------------------------------------------------------------
// Single-column cursor compatible with filters
// ---------------------------------------------------------------------------

func TestCursor_CompatibleWithFilters(t *testing.T) {
	query, args, err := NewBuilder().
		Table("users").
		Model(&cursorTestModel{}).
		Limit(10).
		After("id", 200).
		OrderBy("id").
		Eq("status", "active").
		WhereLike("name", "john").
		BuildSQL()

	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("missing cursor clause in: %s", query)
	}
	if !strings.Contains(query, "users.status") {
		t.Errorf("missing status filter in: %s", query)
	}
	if !strings.Contains(query, "users.name") {
		t.Errorf("missing name LIKE filter in: %s", query)
	}
	_ = args
}

// ---------------------------------------------------------------------------
// Multi-column keyset SQL generation
// ---------------------------------------------------------------------------

func TestKeyset_MultiSort_AfterSQL(t *testing.T) {
	// cols=["created_at","id"], sortDirs=["DESC","ASC"], dir="after"
	// expects: ((users.created_at < $1) OR (users.created_at = $2 AND users.id > $3))
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		NoOffset:        true,
		CursorColumns:   []string{"created_at", "id"},
		CursorValues:    []any{"2024-01-01", 42},
		CursorSortDirs:  []string{"DESC", "ASC"},
		CursorDirection: "after",
		SortColumns:     []string{"created_at", "id"},
		SortDirections:  []string{"DESC", "ASC"},
	}

	query, args := params.GenerateSQL()

	if !strings.Contains(query, "users.created_at < $") {
		t.Errorf("expected DESC column to use '<' for after, got: %s", query)
	}
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected ASC column to use '>' for after, got: %s", query)
	}
	if strings.Contains(query, "OFFSET") {
		t.Errorf("expected no OFFSET, got: %s", query)
	}
	// args should contain the cursor values (may repeat for equality parts)
	if len(args) < 2 {
		t.Errorf("expected at least 2 args, got %d: %v", len(args), args)
	}
}

func TestKeyset_MultiSort_BeforeSQL(t *testing.T) {
	// before reverses operators: DESC→> and ASC→<
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		NoOffset:        true,
		CursorColumns:   []string{"created_at", "id"},
		CursorValues:    []any{"2024-01-01", 42},
		CursorSortDirs:  []string{"DESC", "ASC"},
		CursorDirection: "before",
		SortColumns:     []string{"created_at", "id"},
		SortDirections:  []string{"DESC", "ASC"},
	}

	query, _ := params.GenerateSQL()

	if !strings.Contains(query, "users.created_at > $") {
		t.Errorf("expected DESC column to use '>' for before, got: %s", query)
	}
	if !strings.Contains(query, "users.id < $") {
		t.Errorf("expected ASC column to use '<' for before, got: %s", query)
	}
}

func TestKeyset_SingleColumn_AscAfter(t *testing.T) {
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		NoOffset:        true,
		CursorColumns:   []string{"id"},
		CursorValues:    []any{5},
		CursorSortDirs:  []string{"ASC"},
		CursorDirection: "after",
	}
	query, _ := params.GenerateSQL()
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected users.id > $N, got: %s", query)
	}
}

func TestKeyset_SingleColumn_DescAfter(t *testing.T) {
	params := &QueryParams{
		Table:           "users",
		Struct:          &cursorTestModel{},
		ItemsPerPage:    10,
		Page:            1,
		NoOffset:        true,
		CursorColumns:   []string{"id"},
		CursorValues:    []any{5},
		CursorSortDirs:  []string{"DESC"},
		CursorDirection: "after",
	}
	query, _ := params.GenerateSQL()
	if !strings.Contains(query, "users.id < $") {
		t.Errorf("expected users.id < $N for DESC+after, got: %s", query)
	}
}

// ---------------------------------------------------------------------------
// FromStruct with multi-column token (via NewCursorPage round-trip)
// ---------------------------------------------------------------------------

func TestFromStruct_Cursor_SingleColumn(t *testing.T) {
	token := EncodeCursor("id", 99, "after")

	p := &PaginationParams{Cursor: token, Limit: 10}

	query, _, err := NewBuilder().
		Table("users").
		Model(&cursorTestModel{}).
		OrderBy("id").
		FromStruct(p).
		BuildSQL()

	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected cursor WHERE after FromStruct, got: %s", query)
	}
	if strings.Contains(query, "OFFSET") {
		t.Errorf("expected no OFFSET, got: %s", query)
	}
}

func TestFromStruct_Cursor_MultiColumn(t *testing.T) {
	// Simulate a token produced by NewCursorPage (multi-column)
	token := encodeTokenMulti(
		[]string{"created_at", "id"},
		[]any{"2024-01-01", 42},
		[]string{"DESC", "ASC"},
		"after",
	)

	p := &PaginationParams{Cursor: token, Limit: 10}

	query, _, err := NewBuilder().
		Table("users").
		Model(&cursorTestModel{}).
		OrderBy("created_at", "DESC").
		OrderBy("id").
		FromStruct(p).
		BuildSQL()

	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}
	if !strings.Contains(query, "users.created_at < $") {
		t.Errorf("expected keyset clause for created_at DESC+after, got: %s", query)
	}
	if !strings.Contains(query, "users.id > $") {
		t.Errorf("expected keyset clause for id ASC+after, got: %s", query)
	}
}

// ---------------------------------------------------------------------------
// NewCursorPage — new signature (rawItems, params, baseURL)
// ---------------------------------------------------------------------------

func TestNewCursorPage_Links(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users?status=active")

	// 11 items for limit=10 → hasNext=true
	rawItems := make([]cursorTestModel, 11)
	for i := range rawItems {
		rawItems[i] = cursorTestModel{ID: 10 + i, CreatedAt: time.Now()}
	}

	params := &PaginationParams{
		Limit: 10,
		Sort:  []string{"-created_at", "id"},
	}

	page := NewCursorPage(rawItems, params, base)

	if !page.Meta.HasNext {
		t.Error("expected HasNext true")
	}
	if page.Meta.HasPrev {
		t.Error("expected HasPrev false (no cursor in params)")
	}
	if len(page.Data) != 10 {
		t.Errorf("expected data trimmed to 10, got %d", len(page.Data))
	}
	if page.Links.Next == nil {
		t.Fatal("expected Next link, got nil")
	}
	if page.Links.Prev != nil {
		t.Errorf("expected nil Prev (first page), got %s", *page.Links.Prev)
	}
	if !strings.Contains(*page.Links.Next, "cursor=") {
		t.Errorf("Next link missing cursor param: %s", *page.Links.Next)
	}
	if !strings.Contains(*page.Links.Next, "status=active") {
		t.Errorf("Next link lost existing query params: %s", *page.Links.Next)
	}
	if strings.Contains(*page.Links.Next, "page=") {
		t.Errorf("Next link should not have page param: %s", *page.Links.Next)
	}
}

func TestNewCursorPage_HasPrev_WhenCursorSet(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users")

	rawItems := make([]cursorTestModel, 5) // less than limit → no next
	for i := range rawItems {
		rawItems[i] = cursorTestModel{ID: i + 1}
	}

	// Simulate second page (cursor was provided)
	token := EncodeCursor("id", 0, "after")
	params := &PaginationParams{
		Limit:  10,
		Sort:   []string{"id"},
		Cursor: token,
	}

	page := NewCursorPage(rawItems, params, base)

	if page.Meta.HasNext {
		t.Error("expected HasNext false (fewer items than limit)")
	}
	if !page.Meta.HasPrev {
		t.Error("expected HasPrev true (cursor was set)")
	}
	if page.Links.Prev == nil {
		t.Fatal("expected Prev link, got nil")
	}
}

func TestNewCursorPage_NoLinks_WhenNoSort(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users")
	rawItems := make([]cursorTestModel, 11)
	params := &PaginationParams{Limit: 10} // no sort → no cursor links

	page := NewCursorPage(rawItems, params, base)

	if page.Links.Next != nil {
		t.Errorf("expected nil Next when no sort configured, got %s", *page.Links.Next)
	}
}

func TestNewCursorPage_ExactLimit_NoNext(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users")
	rawItems := []cursorTestModel{{ID: 1}} // 1 item, limit=10 → no next
	params := &PaginationParams{Limit: 10, Sort: []string{"id"}}

	page := NewCursorPage(rawItems, params, base)

	if page.Links.Next != nil {
		t.Errorf("expected nil Next, got %s", *page.Links.Next)
	}
}

// ---------------------------------------------------------------------------
// BindQueryParamsToStruct with cursor
// ---------------------------------------------------------------------------

func TestBindQueryParams_Cursor(t *testing.T) {
	token := EncodeCursor("id", 55, "after")
	values := url.Values{
		"cursor": {token},
		"limit":  {"10"},
	}

	params, err := BindQueryParamsToStruct(values)
	if err != nil {
		t.Fatalf("BindQueryParamsToStruct error: %v", err)
	}
	if params.Cursor != token {
		t.Errorf("cursor not bound: want %q, got %q", token, params.Cursor)
	}
}

// ---------------------------------------------------------------------------
// extractSortInfo
// ---------------------------------------------------------------------------

func TestExtractSortInfo_Sort(t *testing.T) {
	p := &PaginationParams{Sort: []string{"-created_at", "id"}}
	cols, dirs := extractSortInfo(p)
	if len(cols) != 2 || cols[0] != "created_at" || cols[1] != "id" {
		t.Errorf("unexpected cols: %v", cols)
	}
	if dirs[0] != "DESC" || dirs[1] != "ASC" {
		t.Errorf("unexpected dirs: %v", dirs)
	}
}

func TestExtractSortInfo_SortColumns(t *testing.T) {
	p := &PaginationParams{
		SortColumns:    []string{"name", "id"},
		SortDirections: []string{"ASC", "DESC"},
	}
	cols, dirs := extractSortInfo(p)
	if len(cols) != 2 || cols[0] != "name" {
		t.Errorf("unexpected cols: %v", cols)
	}
	if dirs[1] != "DESC" {
		t.Errorf("unexpected dirs: %v", dirs)
	}
}
