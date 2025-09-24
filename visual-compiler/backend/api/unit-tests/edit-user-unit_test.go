package unit_tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestEditUser_NoInput(t *testing.T) {
	contxt, rec := createTestContext(t)
	res, err := http.NewRequest("PATCH", "/api/users/update", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.EditUser(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}
