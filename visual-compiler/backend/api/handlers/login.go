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

// Specifies the JSON body request for logging in
type LoginReq struct {
	// Either the user's email or username
	Login string `json:"login" binding:"required"`
	// User's password
	Password string `json:"password" binding:"required"`
}

// @Summary Login User
// @Description Login user with correct credentials
// @Tags Users
// @Accept json
// @Produce json
// @Param request body LoginReq true "Login User"
// @Success 200 {object} map[string]string "User login successful"
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 404 {object} map[string]string "User not found/Invalid credentials"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/login [post]
func Login(c *gin.Context) {
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	normalise_input := strings.ToLower(req.Login)

	mongo_client := db.ConnectClient()
	users_collection := mongo_client.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter_login := bson.M{
		"$or": []bson.M{
			{"email": normalise_input},
			{"username": normalise_input},
		},
	}

	var db_user struct {
		Username string        `bson:"username"`
		Email    string        `bson:"email"`
		Password string        `bson:"password"`
		ID       bson.ObjectID `bson:"_id"`
		Is_Admin bool          `bson:"is_admin"`
	}

	err := users_collection.FindOne(ctx, filter_login).Decode(&db_user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login Successful. Welcome " + db_user.Username,
		"id":       db_user.ID,
		"is_admin": db_user.Is_Admin,
	})
}
