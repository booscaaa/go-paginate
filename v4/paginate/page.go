package paginate

import (
	"math"
	"net/url"
	"strconv"
)

// Page is the generic paginated response envelope with HATEOAS links.
type Page[T any] struct {
	Data  []T       `json:"data"`
	Meta  PageMeta  `json:"meta"`
	Links PageLinks `json:"links"`
}

// PageMeta contains complete pagination metadata.
type PageMeta struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	From        int  `json:"from"`
	To          int  `json:"to"`
	HasPrev     bool `json:"has_prev"`
	HasNext     bool `json:"has_next"`
}

// PageLinks contains HATEOAS navigation links.
// Prev and Next are pointers so they serialize as null (not omitted) on boundaries.
type PageLinks struct {
	Self  string  `json:"self"`
	First string  `json:"first"`
	Last  string  `json:"last"`
	Prev  *string `json:"prev"`
	Next  *string `json:"next"`
}

// NewPage constructs a Page[T] with full metadata and HATEOAS links.
// baseURL should be the current request URL (*http.Request).URL — all existing
// query params (filters, sorts, limit) are preserved; only ?page=N is rewritten.
//
// page and perPage are derived from params automatically:
//   - page    ← params.Page  (defaults to 1)
//   - perPage ← params.Limit (defaults to GetDefaultLimit())
func NewPage[T any](data []T, totalItems int, params *PaginationParams, baseURL *url.URL) Page[T] {
	page := params.Page
	if page < 1 {
		page = 1
	}
	perPage := params.Limit
	if perPage <= 0 {
		perPage = GetDefaultLimit()
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	from := (page-1)*perPage + 1
	to := page * perPage
	if to > totalItems {
		to = totalItems
	}
	if totalItems == 0 {
		from, to = 0, 0
	}

	meta := PageMeta{
		CurrentPage: page,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		From:        from,
		To:          to,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
	}

	pageURL := func(p int) string {
		u := *baseURL
		q := u.Query()
		q.Set("page", strconv.Itoa(p))
		u.RawQuery = q.Encode()
		return u.String()
	}

	links := PageLinks{
		Self:  pageURL(page),
		First: pageURL(1),
		Last:  pageURL(totalPages),
	}
	if page > 1 {
		prev := pageURL(page - 1)
		links.Prev = &prev
	}
	if page < totalPages {
		next := pageURL(page + 1)
		links.Next = &next
	}

	return Page[T]{Data: data, Meta: meta, Links: links}
}
