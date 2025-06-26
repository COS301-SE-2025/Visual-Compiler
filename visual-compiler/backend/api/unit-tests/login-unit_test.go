package unit_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestLoginNoName(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.LoginReq{
		Password: "passsdw32323@@@",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("GET", "/api/users/login", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Login(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestLoginNoPassword(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.LoginReq{
		Login: "passsdw32323@@@",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("GET", "/api/users/login", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Login(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}
