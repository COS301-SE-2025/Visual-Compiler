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

// Name: Register
//
// Parameters: Gin Context
//
// Return: None
//
// Registers a user into the database. If any inputs are missing, an error is displayed
func Register(c *gin.Context) {
	var req Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongo_client := db.ConnectClient()
	users_collection := mongo_client.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	normalised_username := strings.ToLower(req.Username)

	filter_checks := bson.M{
		"$or": []bson.M{
			{"email": req.Email},
			{"username": normalised_username},
		},
	}

	var user_exists bson.M
	err := users_collection.FindOne(ctx, filter_checks).Decode(&user_exists)
	if err == nil {
		if user_exists["email"] == req.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		}
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in hashing password"})
		return
	}

	_, err = users_collection.InsertOne(ctx, bson.M{
		"email":    req.Email,
		"password": string(hashed_password),
		"username": normalised_username,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in registering user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully registered user"})
}
