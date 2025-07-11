package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Name: ConnectToMongo
//
// Parameters: Gin Context
//
// Return: None
//
// Sets up the initial connection to the Mongo Database. Any subsequent calls to connect to Mongo will use this one and not create a new one each time
func ConnectToMongo(c *gin.Context) {
	client := db.ConnectClient()

	cont, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(cont, readpref.Primary()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "MongoDB unreachable",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Connected to MongoDB",
	})
}
