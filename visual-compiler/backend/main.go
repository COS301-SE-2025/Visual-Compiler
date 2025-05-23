package main

import (
	"log"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
)

func main() {
	db.ConnectClient()

	router := routers.SetupRouter()

	log.Println("Starting backend server on: http://localhost:8080/api")
	if error := router.Run(":8080"); error != nil {
		log.Fatalf("Error starting server: %v", error)
	}
}
