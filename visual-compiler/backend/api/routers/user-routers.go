// Package Routers defines the endpoints that the frontend calls.
package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Name: SetupUserRouter
//
// Parameters: None
//
// Return: Endpoints
//
// Creates the endpoints for user management. Links the endpoints to the respective function
func SetupUserRouter() *gin.Engine {
	r := gin.New()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/getUsers", handlers.GetAllUsers)
	r.DELETE("/delete", handlers.DeleteUser)
	r.POST("/save", handlers.SaveProjectName)
	r.POST("/savePipeline", handlers.SaveProjectPipeline)
	r.GET("/getProjects", handlers.GetAllProjects)
	r.GET("/getProject", handlers.GetProject)
	r.DELETE("/deleteProject", handlers.DeleteProject)
	r.PATCH("/update", handlers.EditUser)

	return r
}
