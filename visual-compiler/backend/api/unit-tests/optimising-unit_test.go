package unit_tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// create mock requests
func createOptmiseTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}

func TestStoreOptimisingCode_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createOptmiseTestContext(t)

	res, err := http.NewRequest("POST", "/api/optimising/source_code", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.StoreOptimisingCode(contxt)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("StatusUnauthorized status code expected")
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
		if body_array["error"] != "Unauthorized" {
			t.Errorf("Incorrect error")
		}
	}

}

func TestOptimiseCode_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createOptmiseTestContext(t)

	res, err := http.NewRequest("POST", "/api/optimising/optimise", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.OptimiseCode(contxt)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("StatusUnauthorized status code expected")
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
		if body_array["error"] != "Unauthorized" {
			t.Errorf("Incorrect error")
		}
	}

}
