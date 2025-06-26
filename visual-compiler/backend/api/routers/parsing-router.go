package routers

import (
	// "github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Sets up the router for the parsing specific endpoints
// Sets up the following:
//   - POST request: parsing functionality
func SetupParsingRouter() *gin.Engine {
	r := gin.New()

	r.POST("/grammar", handlers.ReadGrammar)
	r.POST("/tree", handlers.CreateSyntaxTree)

	return r
}
