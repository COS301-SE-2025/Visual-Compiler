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
	"golang.org/x/crypto/bcrypt"
)

// Specifies the JSON body request for Registering
type Request struct {
	// User's Email address
	Email string `json:"email" binding:"required,email"`
	// User's password
	Password string `json:"password" binding:"required,min=8"`
	// User's username
	Username string `json:"username" binding:"required,min=6"`
}

// Registers a user on the Database
// Gets the user's details from a JSON request.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for registering/inserting or parsing
func Register(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongoClient := db.ConnectClient()
	usersCollection := mongoClient.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	normalisedUsername := strings.ToLower(req.Username)

	filterChecks := bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": normalisedUsername},
		},
	}

	var userExists bson.M
	err := usersCollection.FindOne(ctx, filterChecks).Decode(&userExists)
	if err == nil {
		if userExists["email"] == req.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		}
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in hashing password"})
		return
	}

	_, err = usersCollection.InsertOne(ctx, bson.M{
		"email":    req.Email,
		"password": string(hashedPassword),
		"username": normalisedUsername,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in registering user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully registered user"})
}
