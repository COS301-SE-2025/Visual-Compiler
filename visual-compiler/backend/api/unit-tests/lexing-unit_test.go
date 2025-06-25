package unit_tests

import (
	"bytes"

	//"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
)

// create mock requests
func createPhaseTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}

func TestStoreSourceCode_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/code", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.StoreSourceCode(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestLexing_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/lexer", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Lexing(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestCreateRulesFromCode_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/rules", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Lexing(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestReadDFAFromUser_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfa", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ReadDFAFromUser(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestTokensFromDFA_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfaToTokens", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.TokensFromDFA(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestConvertDFAToRG_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/dfaToRegex", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertDFAToRG(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestConvertRGToNFA_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/regexToNFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertRGToNFA(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestConvertRGToDFA_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/regexToDFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertRGToDFA(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}

func TestConvertNFAToDFA_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/lexing/nfaToDFA", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.ConvertNFAToDFA(contxt)

	if rec.Code != http.StatusBadRequest {
		var mess map[string]string
		t.Errorf(mess["error"])
	}
}
