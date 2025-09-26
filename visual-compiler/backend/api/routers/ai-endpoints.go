//go:build !coverage
// +build !coverage

package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Name: SetupAIRouter
//
// Parameters: None
//
// Return: Endpoints
//
// Creates the endpoints for the AI assistant. Links the endpoints to the respective function
func SetupAIRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/answer", handlers.AnswerUsersQuestion)
	r.POST("/generate", handlers.InputGeneration)

	return r
}
