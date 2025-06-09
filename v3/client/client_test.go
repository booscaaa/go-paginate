package client

import (
	"net/url"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("https://api.example.com/users")
	if client.baseURL != "https://api.example.com/users" {
		t.Errorf("Expected baseURL to be 'https://api.example.com/users', got '%s'", client.baseURL)
	}
	if client.params == nil {
		t.Error("Expected params to be initialized")
	}
}

func TestNewFromURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedURL string
		expectedErr bool
	}{
		{
			name:        "valid URL with query params",
			input:       "https://api.example.com/users?page=2&limit=10",
			expectedURL: "https://api.example.com/users",
			expectedErr: false,
		},
		{
			name:        "valid URL without query params",
			input:       "https://api.example.com/users",
			expectedURL: "https://api.example.com/users",
			expectedErr: false,
		},
		{
			name:        "invalid URL",
			input:       "://invalid-url",
			expectedURL: "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewFromURL(tt.input)
			if tt.expectedErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if client.baseURL != tt.expectedURL {
				t.Errorf("Expected baseURL to be '%s', got '%s'", tt.expectedURL, client.baseURL)
			}
		})
	}
}

func TestBasicPagination(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Page(2).Limit(25).BuildURL()
	
	expected := "https://api.example.com/users?limit=25&page=2"
	if url != expected {
		t.Errorf("Expected URL '%s', got '%s'", expected, url)
	}
}

func TestSearch(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Search("john").SearchFields("name", "email").BuildURL()
	
	if !strings.Contains(url, "search=john") {
		t.Error("Expected URL to contain 'search=john'")
	}
	if !strings.Contains(url, "search_fields=name%2Cemail") {
		t.Error("Expected URL to contain encoded search fields")
	}
}

func TestSort(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Sort("name", "-created_at").BuildURL()
	
	if !strings.Contains(url, "sort=name") {
		t.Error("Expected URL to contain 'sort=name'")
	}
	if !strings.Contains(url, "sort=-created_at") {
		t.Error("Expected URL to contain 'sort=-created_at'")
	}
}

func TestLikeFilters(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Like("name", "john", "jane").LikeOr("status", "active", "pending").BuildURL()
	
	if !strings.Contains(url, "like%5Bname%5D=john") {
		t.Error("Expected URL to contain encoded like filter")
	}
	if !strings.Contains(url, "likeor%5Bstatus%5D=active") {
		t.Error("Expected URL to contain encoded likeor filter")
	}
}

func TestEqualityFilters(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Eq("age", 25).EqOr("department", "IT", "HR").BuildURL()
	
	if !strings.Contains(url, "eq%5Bage%5D=25") {
		t.Error("Expected URL to contain encoded eq filter")
	}
	if !strings.Contains(url, "eqor%5Bdepartment%5D=IT") {
		t.Error("Expected URL to contain encoded eqor filter")
	}
}

func TestComparisonFilters(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Gte("age", 18).Lt("score", 100).BuildURL()
	
	if !strings.Contains(url, "gte%5Bage%5D=18") {
		t.Error("Expected URL to contain encoded gte filter")
	}
	if !strings.Contains(url, "lt%5Bscore%5D=100") {
		t.Error("Expected URL to contain encoded lt filter")
	}
}

func TestInFilters(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.In("status", "active", "pending").NotIn("role", "admin").BuildURL()
	
	if !strings.Contains(url, "in%5Bstatus%5D=active") {
		t.Error("Expected URL to contain encoded in filter")
	}
	if !strings.Contains(url, "notin%5Brole%5D=admin") {
		t.Error("Expected URL to contain encoded notin filter")
	}
}

func TestBetweenFilter(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Between("age", 18, 65).BuildURL()
	
	if !strings.Contains(url, "between%5Bage%5D%5B0%5D=18") {
		t.Error("Expected URL to contain encoded between min filter")
	}
	if !strings.Contains(url, "between%5Bage%5D%5B1%5D=65") {
		t.Error("Expected URL to contain encoded between max filter")
	}
}

func TestNullFilters(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.IsNull("deleted_at").IsNotNull("email").BuildURL()
	
	if !strings.Contains(url, "isnull=deleted_at") {
		t.Error("Expected URL to contain isnull filter")
	}
	if !strings.Contains(url, "isnotnull=email") {
		t.Error("Expected URL to contain isnotnull filter")
	}
}

func TestVacuumAndNoOffset(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Vacuum(true).NoOffset(true).BuildURL()
	
	if !strings.Contains(url, "vacuum=true") {
		t.Error("Expected URL to contain vacuum=true")
	}
	if !strings.Contains(url, "no_offset=true") {
		t.Error("Expected URL to contain no_offset=true")
	}
}

func TestColumns(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.Columns("id", "name", "email").BuildURL()
	
	if !strings.Contains(url, "columns=id%2Cname%2Cemail") {
		t.Error("Expected URL to contain encoded columns")
	}
}

func TestClone(t *testing.T) {
	client := New("https://api.example.com/users")
	client.Page(1).Limit(10)
	
	cloned := client.Clone()
	cloned.Page(2).Limit(20)
	
	originalURL := client.BuildURL()
	clonedURL := cloned.BuildURL()
	
	if originalURL == clonedURL {
		t.Error("Expected cloned client to have different parameters")
	}
	
	if !strings.Contains(originalURL, "page=1") {
		t.Error("Expected original client to maintain page=1")
	}
	if !strings.Contains(clonedURL, "page=2") {
		t.Error("Expected cloned client to have page=2")
	}
}

func TestReset(t *testing.T) {
	client := New("https://api.example.com/users")
	client.Page(1).Limit(10).Search("test")
	
	urlBefore := client.BuildURL()
	client.Reset()
	urlAfter := client.BuildURL()
	
	if strings.Contains(urlAfter, "page=") || strings.Contains(urlAfter, "limit=") || strings.Contains(urlAfter, "search=") {
		t.Error("Expected all parameters to be cleared after reset")
	}
	
	if urlBefore == urlAfter {
		t.Error("Expected URL to change after reset")
	}
}

func TestCustomParams(t *testing.T) {
	client := New("https://api.example.com/users")
	url := client.SetCustomParam("custom", "value").AddCustomParam("multi", "value1").AddCustomParam("multi", "value2").BuildURL()
	
	if !strings.Contains(url, "custom=value") {
		t.Error("Expected URL to contain custom parameter")
	}
	if !strings.Contains(url, "multi=value1") || !strings.Contains(url, "multi=value2") {
		t.Error("Expected URL to contain multiple values for multi parameter")
	}
}

func TestRemoveParam(t *testing.T) {
	client := New("https://api.example.com/users")
	client.Page(1).Limit(10).Search("test")
	client.RemoveParam("search")
	
	url := client.BuildURL()
	if strings.Contains(url, "search=") {
		t.Error("Expected search parameter to be removed")
	}
	if !strings.Contains(url, "page=1") {
		t.Error("Expected page parameter to remain")
	}
}

func TestGetParams(t *testing.T) {
	client := New("https://api.example.com/users")
	client.Page(1).Limit(10)
	
	params := client.GetParams()
	if params.Get("page") != "1" {
		t.Error("Expected page parameter to be '1'")
	}
	if params.Get("limit") != "10" {
		t.Error("Expected limit parameter to be '10'")
	}
	
	// Test that returned params are a copy
	params.Set("page", "999")
	if client.GetParams().Get("page") == "999" {
		t.Error("Expected returned params to be a copy, not reference")
	}
}

func TestBuildQueryString(t *testing.T) {
	client := New("https://api.example.com/users")
	queryString := client.Page(1).Limit(10).BuildQueryString()
	
	if !strings.Contains(queryString, "page=1") {
		t.Error("Expected query string to contain page=1")
	}
	if !strings.Contains(queryString, "limit=10") {
		t.Error("Expected query string to contain limit=10")
	}
	if strings.HasPrefix(queryString, "?") {
		t.Error("Expected query string to not start with '?'")
	}
}

func TestComplexExample(t *testing.T) {
	client := New("https://api.example.com/users")
	generatedURL := client.
		Page(2).
		Limit(25).
		Search("john").
		SearchFields("name", "email").
		Sort("name", "-created_at").
		LikeOr("status", "active", "pending").
		EqOr("age", 25, 30).
		Gte("created_at", "2023-01-01").
		Lt("score", 100).
		In("department", "IT", "HR").
		IsNotNull("email").
		Vacuum(true).
		BuildURL()
	
	// Verify that the URL contains all expected parameters
	expectedParams := []string{
		"page=2",
		"limit=25",
		"search=john",
		"vacuum=true",
	}
	
	for _, param := range expectedParams {
		if !strings.Contains(generatedURL, param) {
			t.Errorf("Expected URL to contain '%s', got: %s", param, generatedURL)
		}
	}
	
	// Parse the URL to verify it's valid
	parsedURL, err := url.Parse(generatedURL)
	if err != nil {
		t.Errorf("Generated URL is not valid: %v", err)
	}
	
	if parsedURL.Scheme != "https" {
		t.Error("Expected HTTPS scheme")
	}
}