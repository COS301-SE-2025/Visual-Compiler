// Package Handlers abstracts the deleting functionality from the DELETE route themselves
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// A DeleteRequest specifies what is required from the user as a JSON body DELETE request
type DeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

// DeleteUser returns JSON objects notifying if the user deletion was successful or not
// The returned JSON uses JSON formatting
// All errors are handled with a JSON object starting with the word 'error'
// Successful requests are a JSON object with the word 'message'
func DeleteUser(c *gin.Context) {
	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	objectID, err := bson.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	client := db.ConnectClient()
	collection := client.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "id": req.ID})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
