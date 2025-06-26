package unit_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestReadGrammar_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/parsing/grammar", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.StoreSourceCode(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
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
			t.Errorf("BadRequest status code expected")
		}
	}
}

func TestCreateSyntaxTree_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/parsing/tree", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.StoreSourceCode(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
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
