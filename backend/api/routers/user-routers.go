package routers

import (
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/user", handlers.userHandler)
	}

	return r
}
