package main

import (
	"log"
)

func main() {
	router := routers.setupRouter()

	log.Println("Starting backend server on: http://localhost:8080")
	if error := router.Run(":8080"); error != nil {
		log.Fatalf("Error starting server: %v", error)
	}
}
