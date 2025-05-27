package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupLexingRouter() *gin.Engine {
	r := gin.New()

	r.GET("/lexer", handlers.Lexing)

	return r
}
