package paginate_test

import (
	"encoding/json"
	"net/url"
	"testing"

	"github.com/booscaaa/go-paginate/v4/paginate"
)

func TestNewPage_MiddlePage(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users?limit=10&sort=name")
	items := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	page := paginate.NewPage(items, 95, 3, 10, base)

	if page.Meta.CurrentPage != 3 {
		t.Errorf("expected current_page 3, got %d", page.Meta.CurrentPage)
	}
	if page.Meta.TotalPages != 10 {
		t.Errorf("expected total_pages 10, got %d", page.Meta.TotalPages)
	}
	if page.Meta.From != 21 {
		t.Errorf("expected from 21, got %d", page.Meta.From)
	}
	if page.Meta.To != 30 {
		t.Errorf("expected to 30, got %d", page.Meta.To)
	}
	if !page.Meta.HasPrev {
		t.Error("expected has_prev true")
	}
	if !page.Meta.HasNext {
		t.Error("expected has_next true")
	}
	if page.Links.Prev == nil {
		t.Error("expected prev link to be set")
	}
	if page.Links.Next == nil {
		t.Error("expected next link to be set")
	}
}

func TestNewPage_FirstPage(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users?limit=10")
	items := []string{"a", "b"}

	page := paginate.NewPage(items, 25, 1, 10, base)

	if page.Meta.HasPrev {
		t.Error("expected has_prev false on first page")
	}
	if page.Links.Prev != nil {
		t.Error("expected prev link to be nil on first page")
	}
	if !page.Meta.HasNext {
		t.Error("expected has_next true when not on last page")
	}
	if page.Links.Next == nil {
		t.Error("expected next link to be set")
	}
}

func TestNewPage_LastPage(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users")
	items := []string{"z"}

	page := paginate.NewPage(items, 21, 3, 10, base)

	if page.Meta.HasNext {
		t.Error("expected has_next false on last page")
	}
	if page.Links.Next != nil {
		t.Error("expected next link to be nil on last page")
	}
	if page.Meta.To != 21 {
		t.Errorf("expected to 21 (capped at total), got %d", page.Meta.To)
	}
}

func TestNewPage_EmptyData(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users")
	items := []string{}

	page := paginate.NewPage(items, 0, 1, 10, base)

	if page.Meta.From != 0 {
		t.Errorf("expected from 0 for empty result, got %d", page.Meta.From)
	}
	if page.Meta.To != 0 {
		t.Errorf("expected to 0 for empty result, got %d", page.Meta.To)
	}
	if page.Meta.TotalPages != 1 {
		t.Errorf("expected total_pages 1 minimum, got %d", page.Meta.TotalPages)
	}
}

func TestNewPage_QueryParamsPreserved(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/users?limit=5&sort=name&filter=active")
	items := []int{1, 2, 3, 4, 5}

	page := paginate.NewPage(items, 20, 2, 5, base)

	parsedSelf, _ := url.Parse(page.Links.Self)
	if parsedSelf.Query().Get("limit") != "5" {
		t.Errorf("expected limit=5 preserved in self link, got %s", parsedSelf.Query().Get("limit"))
	}
	if parsedSelf.Query().Get("sort") != "name" {
		t.Errorf("expected sort=name preserved in self link")
	}
	if parsedSelf.Query().Get("page") != "2" {
		t.Errorf("expected page=2 in self link, got %s", parsedSelf.Query().Get("page"))
	}

	parsedFirst, _ := url.Parse(page.Links.First)
	if parsedFirst.Query().Get("page") != "1" {
		t.Errorf("expected page=1 in first link, got %s", parsedFirst.Query().Get("page"))
	}

	parsedLast, _ := url.Parse(page.Links.Last)
	if parsedLast.Query().Get("page") != "4" {
		t.Errorf("expected page=4 in last link, got %s", parsedLast.Query().Get("page"))
	}
}

func TestNewPage_JSONSerialization(t *testing.T) {
	base, _ := url.Parse("https://api.example.com/items")
	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	items := []Item{{ID: 1, Name: "foo"}, {ID: 2, Name: "bar"}}

	page := paginate.NewPage(items, 50, 1, 10, base)

	data, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("failed to marshal Page: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if _, ok := decoded["data"]; !ok {
		t.Error("expected 'data' key in JSON output")
	}
	if _, ok := decoded["meta"]; !ok {
		t.Error("expected 'meta' key in JSON output")
	}
	if _, ok := decoded["links"]; !ok {
		t.Error("expected 'links' key in JSON output")
	}

	// prev should be null on page 1 (not omitted)
	links := decoded["links"].(map[string]any)
	if _, ok := links["prev"]; !ok {
		t.Error("expected 'prev' key present in links (as null), not omitted")
	}
}
