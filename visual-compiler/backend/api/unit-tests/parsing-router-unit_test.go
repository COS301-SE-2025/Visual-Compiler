package unit_tests

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupParsingRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router := routers.SetupParsingRouter(r.Group("/"))
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
	endpoints := r.Routes()
	if len(endpoints) != 4 {
		t.Errorf("Amount of routes does not match")
	}
}
