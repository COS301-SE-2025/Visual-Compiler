package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Specifies the JSON body request.
type LexingRequest struct {
	// Represents the source code the user enters
	Code string `json:"source_code" bson:"required"`
	// Represents the pairs of Type and Regex
	Pairs []services.TypeRegex `json:"pairs" bson:"required"`
}

// Lexes a user's source code.
// Gets the source code from a JSON request.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for lexing or parsing
func Lexing(c *gin.Context) {
	var req LexingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	jsonAsBytes, err := json.Marshal(req.Pairs)
	if err != nil {
		panic(err)
	}

	pairs := jsonAsBytes

	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	services.SourceCode(req.Code)
	services.ReadRegexRules(pairs)

	tokens := services.CreateTokens()

	_, err = collection.InsertOne(ctx, bson.M{
		"code":   req.Code,
		"tokens": tokens,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully tokenised your code",
		"tokens":  tokens,
	})
}
