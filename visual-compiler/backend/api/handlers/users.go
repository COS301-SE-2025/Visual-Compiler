package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Specifies the JSON body response after successful fetching
type UserPublic struct {
	// User's auto-generated ID
	ID bson.ObjectID `json:"users_id"`
	// User's Username
	Username string `json:"username"`
	// User's Email
	Email string `json:"email"`
	// User's projects
	Projects bson.A `json:"projects"`
}

// @Summary Get all users
// @Description Gets all users currently in the database
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Users successfully found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/getUsers [get]
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
		var projects bson.A
		if raw["projects"] != nil {
			projects = raw["projects"].(bson.A)
		} else {
			projects = bson.A{}
		}

		all_users_public = append(all_users_public, UserPublic{
			ID:       id,
			Username: username,
			Email:    email,
			Projects: projects,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully fetched all users",
		"users":   all_users_public,
	})
}

type AdminRequest struct {
	// Represents the Admin's ID from frontend
	AdminID bson.ObjectID `json:"admin_id" binding:"required"`
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// Optional: User's new username
	Username string `json:"username,omitempty"`
	// Optional: User's new email address
	Email string `json:"email,omitempty"`
}

// @Summary Update User
// @Description Update user details. Either updating username, email or both
// @Tags Users
// @Accept json
// @Produce json
// @Param request body AdminRequest true "Edit User"
// @Success 200 {object} map[string]string "User update successful"
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 403 {object} map[string]string "Unauthorized access"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 409 {object} map[string]string "Username/Email taken"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/update [patch]
func EditUser(c *gin.Context) {
	var req AdminRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var admin_verification bson.M
	err := collection.FindOne(ctx, bson.M{"_id": req.AdminID, "is_admin": true}).Decode(&admin_verification)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authorization failed. You are not an admin"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify admin: " + err.Error()})
		return
	}

	var user_to_edit bson.M
	err = collection.FindOne(ctx, bson.M{"_id": req.UsersID}).Decode(&user_to_edit)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find error: " + err.Error()})
		return
	}

	fields_to_update := bson.M{}

	if req.Username != "" {
		normalised_username := strings.ToLower(req.Username)

		checks := bson.M{
			"username": normalised_username,
		}

		var user_exists bson.M
		err := collection.FindOne(ctx, checks).Decode(&user_exists)
		if err == nil {
			if user_exists["username"] == req.Username {
				c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			}
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		fields_to_update["username"] = normalised_username
	}

	if req.Email != "" {
		checks := bson.M{
			"email": req.Email,
		}

		var user_exists bson.M
		err := collection.FindOne(ctx, checks).Decode(&user_exists)
		if err == nil {
			if user_exists["email"] == req.Email {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			}
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		fields_to_update["email"] = req.Email
	}

	if len(fields_to_update) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields were provided to update"})
		return
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"_id": req.UsersID},
		bson.M{"$set": fields_to_update},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user: " + err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated user successfully"})
}
