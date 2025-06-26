package unit_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestRegisterInvalidEmail(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Password: "passsdw32323@@@",
		Username: "invalid Email",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Email:    "user@gmail.com",
		Password: "pas",
		Username: "invalid Password",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestRegisterInvalidUsername(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Email:    "user@gmail.com",
		Password: "passsdw32323@@@",
		Username: "in",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/users/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}

}
