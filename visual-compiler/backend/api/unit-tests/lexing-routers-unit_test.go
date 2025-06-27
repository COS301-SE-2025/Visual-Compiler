package unit_tests

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

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
	if len(endpoints) != 9 {
		t.Errorf("Amount of routes does not match")
	}
}
