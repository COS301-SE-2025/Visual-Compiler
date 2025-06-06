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

// Returns all users currently in the Database as an array.
// It fetches the users from the "users" collection.
// Formats the response as JSON
//
// Returns:
//   - An array of UserPublic structs containing each user's details.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for fetching or parsing
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

	var allUsersPublic []UserPublic
	for pointer.Next(ctx) {
		var raw bson.M
		if err := pointer.Decode(&raw); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user: " + err.Error()})
			return
		}

		id := raw["_id"].(bson.ObjectID)
		username := raw["username"].(string)
		email := raw["email"].(string)

		allUsersPublic = append(allUsersPublic, UserPublic{
			ID:       id,
			Username: username,
			Email:    email,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched all users",
		"users":   allUsersPublic,
	})
}
