package unit_tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupUserRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupUserRouter()
	endpoints := router.Routes()
	if len(endpoints) != 10 {
		t.Errorf("Expected number of routes not registered")
	}
}

// create mock contexts
func createTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}
