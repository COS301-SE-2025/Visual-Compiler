// Package Handlers abstracts the functionality from the endpoints.
package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Specifies what is required from the user as a JSON body DELETE request
type DeleteRequest struct {
	ID bson.ObjectID `json:"users_id" binding:"required" example:"685df259c1294de5546b045f"`
}

// @Summary Delete User
// @Description Delete user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param request body DeleteRequest true "Delete user by ID"
// @Success 200 {object} map[string]string "User deleted successfully"
// @Failure 400 {object} map[string]string "Invalid input or ID format"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/delete [delete]
func DeleteUser(c *gin.Context) {
	var req DeleteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	client := db.ConnectClient()
	users_collection := client.Database("visual-compiler").Collection("users")
	lexing_collection := client.Database("visual-compiler").Collection("lexing")
	parsing_collection := client.Database("visual-compiler").Collection("parsing")
	analysing_collection := client.Database("visual-compiler").Collection("analysing")
	translating_collection := client.Database("visual-compiler").Collection("translating")
	optimising_collection := client.Database("visual-compiler").Collection("optimising")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := users_collection.DeleteOne(ctx, bson.M{"_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "id": req.ID})
		return
	}

	lexing_res, err := lexing_collection.DeleteOne(ctx, bson.M{"users_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's related lexing data: " + err.Error()})
		return
	}

	parsing_res, err := parsing_collection.DeleteOne(ctx, bson.M{"users_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's related parsing data: " + err.Error()})
		return
	}

	analysing_res, err := analysing_collection.DeleteOne(ctx, bson.M{"users_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's related analysing data: " + err.Error()})
		return
	}

	translating_res, err := translating_collection.DeleteOne(ctx, bson.M{"users_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's related translating data: " + err.Error()})
		return
	}
	optimising_res, err := optimising_collection.DeleteOne(ctx, bson.M{"users_id": req.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user's related optimising data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":             "User deleted successfully",
		"user":                res.DeletedCount,
		"lexing_deleted":      lexing_res.DeletedCount,
		"parsing_deleted":     parsing_res.DeletedCount,
		"analysing_deleted":   analysing_res.DeletedCount,
		"translating_deleted": translating_res.DeletedCount,
		"optimising_deleted":  optimising_res.DeletedCount,
	})
}
