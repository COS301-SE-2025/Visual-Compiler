package core

import (
	"fmt" 
	routers "github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-contrib/cors"
)

func StartServer() {
	router:= routers.SetupRouter()
	router.Use(cors.Default())
	fmt.Printf("server is alive on port 8080")
	router.Run(":8080")
}