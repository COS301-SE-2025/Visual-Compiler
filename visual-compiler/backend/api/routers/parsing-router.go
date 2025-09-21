package routers

import (
	// "github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Name: SetupParsingRouter
//
// Parameters: None
//
// Return: Endpoints
//
// Creates the endpoints for parsing. Links the endpoints to the respective function
func SetupParsingRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/grammar", handlers.ReadGrammar)
	r.POST("/tree", handlers.CreateSyntaxTree)
	r.POST("/treeString", handlers.TreeToString)
	r.GET("/getTree", handlers.GetTree)

	return r
}
