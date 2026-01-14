package paginate

import (
	"net/url"
	"strings"
	"testing"
)

func TestIsNullIsNotNullQueryStringBinding(t *testing.T) {
	t.Run("IsNull and IsNotNull from query string", func(t *testing.T) {
		// Simular query string real: ?isnull=deleted_at&isnull=archived_at&isnotnull=email&isnotnull=phone
		values := url.Values{}
		values.Add("isnull", "deleted_at")
		values.Add("isnull", "archived_at")
		values.Add("isnotnull", "email")
		values.Add("isnotnull", "phone")

		// Fazer bind
		params, err := BindQueryParamsToStruct(values)
		if err != nil {
			t.Fatalf("BindQueryParamsToStruct failed: %v", err)
		}

		// Verificar se os campos foram populados corretamente
		if len(params.IsNull) != 2 {
			t.Errorf("Expected IsNull to have 2 elements, got %d: %+v", len(params.IsNull), params.IsNull)
		}
		if len(params.IsNotNull) != 2 {
			t.Errorf("Expected IsNotNull to have 2 elements, got %d: %+v", len(params.IsNotNull), params.IsNotNull)
		}

		// Verificar valores específicos
		if !containsString(params.IsNull, "deleted_at") {
			t.Errorf("Expected IsNull to contain 'deleted_at', got: %+v", params.IsNull)
		}
		if !containsString(params.IsNull, "archived_at") {
			t.Errorf("Expected IsNull to contain 'archived_at', got: %+v", params.IsNull)
		}
		if !containsString(params.IsNotNull, "email") {
			t.Errorf("Expected IsNotNull to contain 'email', got: %+v", params.IsNotNull)
		}
		if !containsString(params.IsNotNull, "phone") {
			t.Errorf("Expected IsNotNull to contain 'phone', got: %+v", params.IsNotNull)
		}

		// Testar com builder
		type TestUser struct {
			ID         int    `json:"id" paginate:"users.id"`
			Email      string `json:"email" paginate:"users.email"`
			Phone      string `json:"phone" paginate:"users.phone"`
			DeletedAt  string `json:"deleted_at" paginate:"users.deleted_at"`
			ArchivedAt string `json:"archived_at" paginate:"users.archived_at"`
		}

		builder := NewBuilder().
			Table("users").
			Model(TestUser{}).
			FromStruct(params)

		queryParams, err := builder.Build()
		if err != nil {
			t.Fatalf("Build failed: %v", err)
		}

		sql, _ := queryParams.GenerateSQL()

		// Verificar se a SQL contém as cláusulas IS NULL e IS NOT NULL
		if !strings.Contains(sql, "users.deleted_at IS NULL") {
			t.Errorf("Expected SQL to contain 'users.deleted_at IS NULL', got: %s", sql)
		}
		if !strings.Contains(sql, "users.archived_at IS NULL") {
			t.Errorf("Expected SQL to contain 'users.archived_at IS NULL', got: %s", sql)
		}
		if !strings.Contains(sql, "users.email IS NOT NULL") {
			t.Errorf("Expected SQL to contain 'users.email IS NOT NULL', got: %s", sql)
		}
		if !strings.Contains(sql, "users.phone IS NOT NULL") {
			t.Errorf("Expected SQL to contain 'users.phone IS NOT NULL', got: %s", sql)
		}
	})
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
