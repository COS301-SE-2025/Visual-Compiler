package unit_tests

import (
	"bytes"
	"encoding/json"
	"io"

	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

func TestStoreSourceCode_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/code", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.StoreSourceCode(contxt)

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

func TestLexing_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/lexer", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Lexing(contxt)

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

func TestCreateRulesFromCode_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/rules", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.CreateRulesFromCode(contxt)

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

func TestReadDFAFromUser_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfa", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ReadDFAFromUser(contxt)

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

func TestTokensFromDFA_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfaToTokens", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.TokensFromDFA(contxt)

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

func TestConvertDFAToRG_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfaToRegex", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertDFAToRG(contxt)

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

func TestConvertRGToNFA_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/regexToNFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertRGToNFA(contxt)

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

func TestConvertRGToDFA_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/regexToDFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertRGToDFA(contxt)

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

func TestConvertNFAToDFA_Unauthorised(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/nfaToDFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertNFAToDFA(contxt)

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
