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

	r.POST("/code", handlers.StoreSourceCode)
	r.POST("/lexer", handlers.Lexing)
	r.POST("/dfa", handlers.ReadDFAFromUser)
	r.POST("/dfaToTokens", handlers.TokensFromDFA)
	r.POST("/dfaToRegex", handlers.ConvertDFAToRG)
	r.POST("/regexToNFA", handlers.ConvertRGToNFA)
	r.POST("/regexToDFA", handlers.ConvertRGToDFA)
	r.POST("/nfaToDFA", handlers.ConvertNFAToDFA)

	return r
}
