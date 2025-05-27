package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupUserRouter() *gin.Engine {
	r := gin.New()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/getUsers", handlers.GetAllUsers)
	r.DELETE("/delete", handlers.DeleteUser)

	return r
}
