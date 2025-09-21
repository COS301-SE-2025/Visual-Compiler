package routers

import (
	// "github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Name: SetupTranslatorRouter
//
// Parameters: None
//
// Return: Endpoints
//
// Creates the endpoints for translation. Links the endpoints to the respective function
func SetupTranslatorRouter() *gin.Engine {
	r := gin.New()

	r.POST("/readRules", handlers.ReadRules)
	r.POST("/translate", handlers.TranslateCode)

	return r
}
