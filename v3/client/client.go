package client

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Client provides methods to build query string parameters compatible with go-paginate
type Client struct {
	baseURL string
	params  url.Values
}

// New creates a new Client instance
func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		params:  make(url.Values),
	}
}

// NewFromURL creates a new Client instance from an existing URL with query parameters
func NewFromURL(fullURL string) (*Client, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	
	baseURL := fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
	client := &Client{
		baseURL: baseURL,
		params:  u.Query(),
	}
	
	return client, nil
}

// Reset clears all query parameters
func (c *Client) Reset() *Client {
	c.params = make(url.Values)
	return c
}

// Clone creates a copy of the current client
func (c *Client) Clone() *Client {
	newParams := make(url.Values)
	for k, v := range c.params {
		newParams[k] = append([]string(nil), v...)
	}
	
	return &Client{
		baseURL: c.baseURL,
		params:  newParams,
	}
}

// Page sets the page number
func (c *Client) Page(page int) *Client {
	if page > 0 {
		c.params.Set("page", strconv.Itoa(page))
	}
	return c
}

// Limit sets the number of items per page
func (c *Client) Limit(limit int) *Client {
	if limit > 0 {
		c.params.Set("limit", strconv.Itoa(limit))
	}
	return c
}

// ItemsPerPage is an alias for Limit
func (c *Client) ItemsPerPage(itemsPerPage int) *Client {
	return c.Limit(itemsPerPage)
}

// Search sets the search term
func (c *Client) Search(search string) *Client {
	if search != "" {
		c.params.Set("search", search)
	}
	return c
}

// SearchFields sets the fields to search in
func (c *Client) SearchFields(fields ...string) *Client {
	if len(fields) > 0 {
		c.params.Set("search_fields", strings.Join(fields, ","))
	}
	return c
}

// Sort sets sorting parameters. Use "-" prefix for descending order
// Example: Sort("name", "-created_at")
func (c *Client) Sort(fields ...string) *Client {
	if len(fields) > 0 {
		c.params.Del("sort")
		for _, field := range fields {
			c.params.Add("sort", field)
		}
	}
	return c
}

// SortColumns sets the columns to sort by
func (c *Client) SortColumns(columns ...string) *Client {
	if len(columns) > 0 {
		c.params.Set("sort_columns", strings.Join(columns, ","))
	}
	return c
}

// SortDirections sets the sort directions (asc/desc)
func (c *Client) SortDirections(directions ...string) *Client {
	if len(directions) > 0 {
		c.params.Set("sort_directions", strings.Join(directions, ","))
	}
	return c
}

// Columns sets the columns to select
func (c *Client) Columns(columns ...string) *Client {
	if len(columns) > 0 {
		c.params.Set("columns", strings.Join(columns, ","))
	}
	return c
}

// Vacuum enables vacuum mode
func (c *Client) Vacuum(enable bool) *Client {
	c.params.Set("vacuum", strconv.FormatBool(enable))
	return c
}

// NoOffset enables no offset mode
func (c *Client) NoOffset(enable bool) *Client {
	c.params.Set("no_offset", strconv.FormatBool(enable))
	return c
}

// Like adds LIKE filter for a field
func (c *Client) Like(field string, values ...string) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("like[%s]", field), value)
	}
	return c
}

// LikeOr adds LIKE OR filter for a field
func (c *Client) LikeOr(field string, values ...string) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("likeor[%s]", field), value)
	}
	return c
}

// LikeAnd adds LIKE AND filter for a field
func (c *Client) LikeAnd(field string, values ...string) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("likeand[%s]", field), value)
	}
	return c
}

// Eq adds equality filter for a field
func (c *Client) Eq(field string, values ...any) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("eq[%s]", field), fmt.Sprintf("%v", value))
	}
	return c
}

// EqOr adds equality OR filter for a field
func (c *Client) EqOr(field string, values ...any) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("eqor[%s]", field), fmt.Sprintf("%v", value))
	}
	return c
}

// EqAnd adds equality AND filter for a field
func (c *Client) EqAnd(field string, values ...any) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("eqand[%s]", field), fmt.Sprintf("%v", value))
	}
	return c
}

// Gte adds greater than or equal filter
func (c *Client) Gte(field string, value any) *Client {
	c.params.Set(fmt.Sprintf("gte[%s]", field), fmt.Sprintf("%v", value))
	return c
}

// Gt adds greater than filter
func (c *Client) Gt(field string, value any) *Client {
	c.params.Set(fmt.Sprintf("gt[%s]", field), fmt.Sprintf("%v", value))
	return c
}

// Lte adds less than or equal filter
func (c *Client) Lte(field string, value any) *Client {
	c.params.Set(fmt.Sprintf("lte[%s]", field), fmt.Sprintf("%v", value))
	return c
}

// Lt adds less than filter
func (c *Client) Lt(field string, value any) *Client {
	c.params.Set(fmt.Sprintf("lt[%s]", field), fmt.Sprintf("%v", value))
	return c
}

// In adds IN filter for a field
func (c *Client) In(field string, values ...any) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("in[%s]", field), fmt.Sprintf("%v", value))
	}
	return c
}

// NotIn adds NOT IN filter for a field
func (c *Client) NotIn(field string, values ...any) *Client {
	for _, value := range values {
		c.params.Add(fmt.Sprintf("notin[%s]", field), fmt.Sprintf("%v", value))
	}
	return c
}

// Between adds BETWEEN filter for a field
func (c *Client) Between(field string, min, max any) *Client {
	c.params.Set(fmt.Sprintf("between[%s][0]", field), fmt.Sprintf("%v", min))
	c.params.Set(fmt.Sprintf("between[%s][1]", field), fmt.Sprintf("%v", max))
	return c
}

// IsNull adds IS NULL filter for a field
func (c *Client) IsNull(fields ...string) *Client {
	for _, field := range fields {
		c.params.Add("isnull", field)
	}
	return c
}

// IsNotNull adds IS NOT NULL filter for a field
func (c *Client) IsNotNull(fields ...string) *Client {
	for _, field := range fields {
		c.params.Add("isnotnull", field)
	}
	return c
}

// BuildURL builds the complete URL with query parameters
func (c *Client) BuildURL() string {
	if len(c.params) == 0 {
		return c.baseURL
	}
	return fmt.Sprintf("%s?%s", c.baseURL, c.params.Encode())
}

// BuildQueryString builds only the query string part
func (c *Client) BuildQueryString() string {
	return c.params.Encode()
}

// GetParams returns a copy of the current query parameters
func (c *Client) GetParams() url.Values {
	params := make(url.Values)
	for k, v := range c.params {
		params[k] = append([]string(nil), v...)
	}
	return params
}

// SetCustomParam sets a custom query parameter
func (c *Client) SetCustomParam(key, value string) *Client {
	c.params.Set(key, value)
	return c
}

// AddCustomParam adds a custom query parameter (allows multiple values)
func (c *Client) AddCustomParam(key, value string) *Client {
	c.params.Add(key, value)
	return c
}

// RemoveParam removes a query parameter
func (c *Client) RemoveParam(key string) *Client {
	c.params.Del(key)
	return c
}