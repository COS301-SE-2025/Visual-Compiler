package routers

import (
	// "github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Name: SetupLexingRouter
//
// Parameters: None
//
// Return: Endpoints
//
// Creates the endpoints for lexing. Links the endpoints to the respective function
func SetupAnalysingRouter() *gin.Engine {
	r := gin.New()

	r.POST("/analyse", handlers.Analyse)

	return r
}
