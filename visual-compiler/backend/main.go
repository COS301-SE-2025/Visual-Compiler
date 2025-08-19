//go:build !coverage
// +build !coverage

// Package Main for the main entry into the API.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Visual Compiler API
// @version 1.0
// @description Documentation for Visual Compiler API endpoints
// @host localhost:8080
// @BasePath /api
func main() {
	db.ConnectClient()

	// Set up Gin engine
	router := gin.Default()

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// health check for AWS
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Attach your routes
	api_user_routes := routers.SetupUserRouter()
	api_lexing_routes := routers.SetupLexingRouter()
	api_parsing_routes := routers.SetupParsingRouter()
	api_analysing_routes := routers.SetupAnalysingRouter()
	api_translating_routes := routers.SetupTranslatorRouter()

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

	router.Any("/api/analysing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		api_analysing_routes.HandleContext(c)
	})

	router.Any("/api/translating/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		api_translating_routes.HandleContext(c)
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Swagger available at: http://localhost:8080/swagger/index.html")

	log.Println("Starting backend server on: http://localhost:8080/api")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
