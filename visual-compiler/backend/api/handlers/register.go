package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
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
	Username string `json:"username" binding:"required,min=1"`
}

// @Summary Register User
// @Description Register user with credentials
// @Tags Users
// @Accept json
// @Produce json
// @Param request body Request true "Register User"
// @Success 201 {object} map[string]string "User registration successful"
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 409 {object} map[string]string "Username taken"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/register [post]
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

	inserted_result, err := users_collection.InsertOne(ctx, bson.M{
		"email":    req.Email,
		"password": string(hashed_password),
		"username": normalised_username,
		"projects": []bson.M{},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in registering user"})
		return
	}
	user_id := inserted_result.InsertedID.(bson.ObjectID).String()

	auth_domain := os.Getenv("AUTH0_DOMAIN")
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	audience := auth_domain + "/api/v2/"

	token_request_body := map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
		"audience":      audience,
		"grant_type":    "client_credentials",
	}
	json_token, _ := json.Marshal(token_request_body)

	res, err := http.Post(auth_domain+"/oauth/token", "application/json", bytes.NewBuffer(json_token))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token Auth0 request failed: " + err.Error()})
		users_collection.DeleteOne(ctx, bson.M{"_id": inserted_result.InsertedID})
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var res_token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &res_token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token parse failed (Auth0)"})
		users_collection.DeleteOne(ctx, bson.M{"_id": inserted_result.InsertedID})
		return
	}

	new_auth0_user := map[string]any{
		"connection": "Username-Password-Authentication",
		"email":      req.Email,
		"password":   req.Password,
		"app_metadata": map[string]string{
			"mongo_id": user_id,
		},
	}
	json_user, _ := json.Marshal(new_auth0_user)

	auth_request, _ := http.NewRequest("POST", auth_domain+"/api/v2/users", bytes.NewBuffer(json_user))
	auth_request.Header.Set("Authorization", "Bearer "+res_token.AccessToken)
	auth_request.Header.Set("Content-Type", "application/json")

	cli := &http.Client{}
	auth_res, err := cli.Do(auth_request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed (Auth0): " + err.Error()})
		users_collection.DeleteOne(ctx, bson.M{"_id": inserted_result.InsertedID})
		return
	}

	if auth_res.StatusCode >= 300 {
		bytes_body_error, _ := io.ReadAll(auth_res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(bytes_body_error, &body_array)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":       "User creation failed (Auth0)",
			"status_code": auth_res.StatusCode,
			"body":        body_array["message"],
		})
		users_collection.DeleteOne(ctx, bson.M{"_id": inserted_result.InsertedID})
		return
	}
	defer auth_res.Body.Close()

	var resp_auth0 struct {
		UserID string `json:"user_id"`
	}
	json.NewDecoder(auth_res.Body).Decode(&resp_auth0)

	_, err = users_collection.UpdateByID(ctx, inserted_result.InsertedID, bson.M{
		"$set": bson.M{"auth0_id": resp_auth0.UserID},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		users_collection.DeleteOne(ctx, bson.M{"_id": inserted_result.InsertedID})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Successfully registered user",
		"mongoDB_id": user_id,
		"auth0_id":   resp_auth0.UserID,
	})
}
