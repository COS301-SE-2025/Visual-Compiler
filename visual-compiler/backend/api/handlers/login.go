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

var UsersID bson.ObjectID

// Specifies the JSON body request for logging in
type LoginReq struct {
	// Either the user's email or username
	Login string `json:"login" binding:"required"`
	// User's password
	Password string `json:"password" binding:"required"`
}

// Logs a user in.
// Gets the inputs(Email/Username & Password) from a JSON request from the user.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for fetching or parsing
func Login(c *gin.Context) {
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	normaliseInput := strings.ToLower(req.Login)

	mongoClient := db.ConnectClient()
	usersCollection := mongoClient.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filterLogin := bson.M{
		"$or": []bson.M{
			{"email": normaliseInput},
			{"username": normaliseInput},
		},
	}

	var dbUser struct {
		Username string        `bson:"username"`
		Email    string        `bson:"email"`
		Password string        `bson:"password"`
		ID       bson.ObjectID `bson:"_id"`
	}

	err := usersCollection.FindOne(ctx, filterLogin).Decode(&dbUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successful. Welcome " + dbUser.Username,
	})
	UsersID = dbUser.ID
}
