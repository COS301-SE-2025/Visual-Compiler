package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Sets up the router for the lexing specific endpoints
// Sets up the following:
//   - POST request: lexer functionality
func SetupLexingRouter() *gin.Engine {
	r := gin.New()

	r.POST("/lexer", handlers.Lexing)

	return r
}
