package main

import (
	"log"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectClient()

	// Set up Gin engine
	router := gin.Default()

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Attach your routes
	api := routers.SetupRouter()
	router.Any("/api/*any", func(c *gin.Context) {
		api.HandleContext(c)
	})

	log.Println("Starting backend server on: http://localhost:8080/api")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
