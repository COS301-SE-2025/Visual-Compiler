package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/ai"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AnswerQuestion struct {
	Question string `json:"question" binding:"required"`
}

// @Summary AI assistant answers user's question
// @Description The user sends a question to the AI assistance about the compilation process
// @Tags AI
// @Accept json
// @Produce json
// @Param request body AnswerQuestion true "Answer user's question"
// @Success 200 {object} map[string]string "AI response received"
// @Failure 400 {object} map[string]string "Invalid input/Response failed"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /ai/answer [post]
func AnswerUsersQuestion(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AnswerQuestion

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbUser struct {
		UsersID bson.ObjectID `bson:"_id"`
		Auth0ID string        `bson:"auth0_id"`
	}

	err := users_collection.FindOne(ctx, bson.M{"auth0_id": authID}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	answer := ai.QuestionAnswer(req.Question, ctx)

	if answer == "The AI Assistant is currently unavailable. Please try again later!" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error.", "details": answer})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "AI response received", "response": answer})
}

type InputGenerate struct {
	Phase    string `json:"phase" binding:"required"`
	Artefact string `json:"artefact" binding:"required"`
}

// @Summary AI assistant generates input
// @Description The AI generates input (upon request) based on the current process the user is on in the compilation process
// @Tags AI
// @Accept json
// @Produce json
// @Param request body InputGenerate true "Generate input"
// @Success 200 {object} map[string]string "AI response received"
// @Failure 400 {object} map[string]string "Invalid input/Response failed"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /ai/generate [post]
func InputGeneration(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req InputGenerate

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var dbUser struct {
		UsersID bson.ObjectID `bson:"_id"`
		Auth0ID string        `bson:"auth0_id"`
	}

	err := users_collection.FindOne(ctx, bson.M{"auth0_id": authID}).Decode(&dbUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	response, err := ai.GenerateInput(req.Phase, req.Artefact, ctx)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI response error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "AI response received",
		"response": response,
	})
}
