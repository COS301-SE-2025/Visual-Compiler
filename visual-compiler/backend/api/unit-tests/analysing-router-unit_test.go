package unit_tests

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupAnalysingRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupAnalysingRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestAnalysingRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupAnalysingRouter()
	endpoints := router.Routes()
	if len(endpoints) != 1 {
		t.Errorf("Amount of routes does not match")
	}
}
