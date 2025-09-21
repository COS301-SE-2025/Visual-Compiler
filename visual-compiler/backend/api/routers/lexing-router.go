package routers

import (
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
func SetupLexingRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/code", handlers.StoreSourceCode)
	r.POST("/rules", handlers.CreateRulesFromCode)
	r.POST("/lexer", handlers.Lexing)
	r.POST("/dfa", handlers.ReadDFAFromUser)
	r.POST("/dfaToTokens", handlers.TokensFromDFA)
	r.POST("/dfaToRegex", handlers.ConvertDFAToRG)
	r.POST("/regexToNFA", handlers.ConvertRGToNFA)
	r.POST("/regexToDFA", handlers.ConvertRGToDFA)
	r.POST("/nfaToDFA", handlers.ConvertNFAToDFA)

	return r
}
