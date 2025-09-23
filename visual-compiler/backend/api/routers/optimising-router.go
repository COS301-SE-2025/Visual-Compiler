package routers

import (
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupOptimisingRouter(r *gin.RouterGroup) *gin.RouterGroup {
	r.POST("/source_code", handlers.StoreOptimisingCode)
	r.POST("/optimise", handlers.OptimiseCode)

	return r
}
