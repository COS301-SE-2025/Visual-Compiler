package unit_tests

import (
	"bytes"
	//"encoding/json"
	"net/http/httptest"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

// create mock requests
func createTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}

func TestSetupLexingRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupLexingRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestLexingRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupLexingRouter()
	endpoints := router.Routes()
	if len(endpoints) != 1 {
		t.Errorf("Amount of routes does not match")
	}
}

func TestLexing_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createTestContext(t)

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
