// Package Routers defines the endpoints that the frontend calls.
package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

// Sets up the router for the user specific endpoints
// Sets up the following:
//   - POST requests: register and login
//   - GET request: get all users
//   - DELETE request: delete user
func SetupUserRouter() *gin.Engine {
	r := gin.New()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/getUsers", handlers.GetAllUsers)
	r.DELETE("/delete", handlers.DeleteUser)

	return r
}
