package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Specifies the JSON body response after successful fetching
type UserPublic struct {
	// User's auto-generated ID
	ID bson.ObjectID `json:"id"`
	// User's Username
	Username string `json:"username"`
	// User's Email
	Email string `json:"email"`
}

// Name: GetAllUsers
//
// Parameters: Gin Context
//
// Return: None
//
// Returns a JSON object with all the users currently in the system.
func GetAllUsers(c *gin.Context) {
	client := db.ConnectClient()
	collection := client.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pointer, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch all users"})
		return
	}
	defer pointer.Close(ctx)

	var all_users_public []UserPublic
	for pointer.Next(ctx) {
		var raw bson.M
		if err := pointer.Decode(&raw); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user: " + err.Error()})
			return
		}

		id := raw["_id"].(bson.ObjectID)
		username := raw["username"].(string)
		email := raw["email"].(string)

		all_users_public = append(all_users_public, UserPublic{
			ID:       id,
			Username: username,
			Email:    email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched all users",
		"users":   all_users_public,
	})
}
