package unit_tests

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupTranslatingRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupTranslatorRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestTranslatingRouterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routers.SetupTranslatorRouter()
	endpoints := router.Routes()
	if len(endpoints) != 2 {
		t.Errorf("Amount of routes does not match")
	}
}
