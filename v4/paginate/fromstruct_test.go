package paginate

import (
	"strings"
	"testing"
)

type fsModel struct {
	ID     int    `json:"id" paginate:"id"`
	Name   string `json:"name" paginate:"name"`
	Email  string `json:"email" paginate:"email"`
	Age    int    `json:"age" paginate:"age"`
	Status string `json:"status" paginate:"status"`
	Score  int    `json:"score" paginate:"score"`
	Role   string `json:"role" paginate:"role"`
	Tag    string `json:"tag" paginate:"tag"`
	Cat    string `json:"cat" paginate:"cat"`
	Del    *bool  `json:"deleted_at" paginate:"deleted_at"`
}

func buildFromStruct(params *PaginationParams) (string, []any, error) {
	return NewBuilder().
		Table("users").
		Model(fsModel{}).
		FromStruct(params).
		BuildSQL()
}

func TestFromStruct_NotIn(t *testing.T) {
	params := NewPaginationParams()
	params.NotIn = map[string][]any{"status": {"banned", "inactive"}}

	sql, args, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "NOT IN") {
		t.Errorf("expected NOT IN in SQL, got: %s", sql)
	}
	if len(args) < 2 {
		t.Errorf("expected 2+ args, got %v", args)
	}
}

func TestFromStruct_GteOr(t *testing.T) {
	params := NewPaginationParams()
	params.GteOr = map[string]any{"age": 18}
	params.LteOr = map[string]any{"score": 100}

	sql, _, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, ">=") {
		t.Errorf("expected >= in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, "<=") {
		t.Errorf("expected <= in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, " OR ") {
		t.Errorf("expected OR grouping in SQL, got: %s", sql)
	}
}

func TestFromStruct_GtOr(t *testing.T) {
	params := NewPaginationParams()
	params.GtOr = map[string]any{"age": 18}

	sql, _, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, " > ") {
		t.Errorf("expected > in SQL, got: %s", sql)
	}
}

func TestFromStruct_LtOr(t *testing.T) {
	params := NewPaginationParams()
	params.LtOr = map[string]any{"score": 50}

	sql, _, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, " < ") {
		t.Errorf("expected < in SQL, got: %s", sql)
	}
}

func TestFromStruct_InOr(t *testing.T) {
	params := NewPaginationParams()
	// Two OR conditions so " OR " appears between them
	params.InOr = map[string][]any{"role": {"admin", "editor"}}
	params.IsNullOr = []string{"deleted_at"}

	sql, args, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "IN") {
		t.Errorf("expected IN in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, " OR ") {
		t.Errorf("expected OR grouping in SQL, got: %s", sql)
	}
	if len(args) < 2 {
		t.Errorf("expected 2+ args, got %v", args)
	}
}

func TestFromStruct_NotInOr(t *testing.T) {
	params := NewPaginationParams()
	// Two OR conditions so " OR " appears between them
	params.NotInOr = map[string][]any{"tag": {"spam", "draft"}}
	params.IsNullOr = []string{"deleted_at"}

	sql, args, err := buildFromStruct(params)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "NOT IN") {
		t.Errorf("expected NOT IN in SQL, got: %s", sql)
	}
	if !strings.Contains(sql, " OR ") {
		t.Errorf("expected OR grouping in SQL, got: %s", sql)
	}
	if len(args) < 2 {
		t.Errorf("expected 2+ args, got %v", args)
	}
}

func TestFromStruct_Vacuum(t *testing.T) {
	params := NewPaginationParams()
	params.Vacuum = true

	b := NewBuilder().Table("users").Model(fsModel{}).FromStruct(params)
	_, _, err := b.BuildCountSQL()
	if err != nil {
		t.Fatal(err)
	}
	if !b.params.Vacuum {
		t.Error("expected Vacuum to be true")
	}
}

func TestFromStruct_NoOffset(t *testing.T) {
	params := NewPaginationParams()
	params.NoOffset = true

	b := NewBuilder().Table("users").Model(fsModel{}).FromStruct(params)
	if !b.params.NoOffset {
		t.Error("expected NoOffset to be true")
	}
}

func TestFromStruct_AllFilters(t *testing.T) {
	params := NewPaginationParams()
	params.Page = 2
	params.Limit = 10
	params.Like = map[string][]string{"name": {"john"}}
	params.LikeOr = map[string][]string{"email": {"gmail"}}
	params.LikeAnd = map[string][]string{"name": {"doe"}}
	params.Eq = map[string][]any{"status": {"active"}}
	params.EqOr = map[string][]any{"role": {"admin", "editor"}}
	params.EqAnd = map[string][]any{"cat": {"tech"}}
	params.Gte = map[string]any{"age": 18}
	params.Lte = map[string]any{"age": 65}
	params.In = map[string][]any{"role": {"admin"}}
	params.NotIn = map[string][]any{"status": {"banned"}}
	params.Between = map[string][2]any{"score": {0, 100}}
	params.IsNull = []string{"deleted_at"}
	params.IsNotNull = []string{"email"}
	params.GteOr = map[string]any{"score": 10}
	params.LteOr = map[string]any{"score": 90}
	params.InOr = map[string][]any{"tag": {"go", "rust"}}
	params.NotInOr = map[string][]any{"cat": {"spam"}}
	params.IsNullOr = []string{"deleted_at"}
	params.IsNotNullOr = []string{"email"}

	sql, _, err := buildFromStruct(params)
	if err != nil {
		t.Fatalf("BuildSQL error: %v", err)
	}

	checks := []string{"ILIKE", "= ", "IN", "NOT IN", ">=", "<=", "BETWEEN", "IS NULL", "IS NOT NULL", "OR"}
	for _, c := range checks {
		if !strings.Contains(sql, c) {
			t.Errorf("expected %q in SQL\nSQL: %s", c, sql)
		}
	}
}
