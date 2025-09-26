package unit_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestDeleteUser_Unauthorised(t *testing.T) {
	contxt, rec := createTestContext(t)
	res, err := http.NewRequest("DELETE", "/api/users/delete", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.DeleteUser(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("StatusBadRequest status code expected")
	} else {
		body_bytes, err := io.ReadAll(rec.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Input is invalid" {
			t.Errorf("Incorrect error")
		}
	}
}
