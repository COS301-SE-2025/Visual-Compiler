package unit_tests

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupParsingRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupParsingRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestParsingRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupParsingRouter()
	endpoints := router.Routes()
	if len(endpoints) != 2 {
		t.Errorf("Amount of routes does not match")
	}
}
