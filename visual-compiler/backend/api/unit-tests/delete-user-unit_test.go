package unit_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestDeleteUser_InvalidInput(t *testing.T) {
	contxt, rec := createTestContext(t)
	res, err := http.NewRequest("DELETE", "/api/users/delete", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.DeleteUser(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestDeleteUser_InvalidID(t *testing.T) {
	contxt, rec := createTestContext(t)
	user_data := map[string]string{
		"id": "123",
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}
	res, err := http.NewRequest("DELETE", "/api/users/delete", bytes.NewBuffer(req))

	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.DeleteUser(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}
