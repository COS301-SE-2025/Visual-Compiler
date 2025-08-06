// Package Main for the main entry into the API.
package main

import (
	"log"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Name: GetAllUsers
//
// Parameters: Gin Context
//
// Return: None
//
// Entry point of the API. Initializes the necessary services, sets up the routes and the GIN framework and starts the http server. Responsibile for bootstrapping the backend infrastructure required. If any setup fails, the application will log the error and exit the application
func main() {
	db.ConnectClient()

	// Set up Gin engine
	router := gin.Default()

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://51.21.245.160:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Attach your routes
	api_user_routes := routers.SetupUserRouter()
	api_lexing_routes := routers.SetupLexingRouter()
	api_parsing_routes := routers.SetupParsingRouter()

	router.Any("/api/users/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		api_user_routes.HandleContext(c)
	})

	router.Any("/api/lexing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		api_lexing_routes.HandleContext(c)
	})

	router.Any("/api/parsing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		api_parsing_routes.HandleContext(c)
	})

	log.Println("Starting backend server on: http://localhost:8080/api")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
