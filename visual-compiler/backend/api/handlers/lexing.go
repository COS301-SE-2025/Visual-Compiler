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
type SourceCodeRequest struct {
	// Represents the source code the user enters
	Code string `json:"source_code" bson:"required"`
	// Represents the pairs of Type and Regex
	Pairs []services.TypeRegex `json:"pairs" bson:"required"`
}

// Locally store a user's source code and regex expressions.
// Gets the source code from a JSON request.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for parsing errors
func StoreSourceCode(c *gin.Context) {
	var req SourceCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	jsonAsBytes, err := json.Marshal(req.Pairs)
	if err != nil {
		panic(err)
	}

	pairs := jsonAsBytes

	services.SourceCode(req.Code)
	services.ReadRegexRules(pairs)

	c.JSON(http.StatusOK, gin.H{"message": "Code is ready for lexing"})
}

// Lexes the user's source code that is locally stored.
// Stores the tokens, unidentified input, user's source code and the user's id in the database.
// Gets the source code from the function GetSourceCode.
// Gets the users id from a global variable `UsersID`.
// Formats the response as a JSON Body
//
// Returns:
//   - A JSON response body.
//   - A 200 OK response if successful
//   - A 500 Internal Server Error if any errors are caught for parsing and lexing errors
func Lexing(c *gin.Context) {
	mongoCli := db.ConnectClient()
	collection := mongoCli.Database("visual-compiler").Collection("lexing")

	services.CreateTokens()

	tokens := services.GetTokens()
	tokens_unidentified := services.GetInvalidInput()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, bson.M{
		"code":     services.GetSourceCode(),
		"tokens":   tokens,
		"users_id": UsersID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users_id":            UsersID,
		"message":             "Successfully tokenised your code",
		"tokens":              tokens,
		"tokens_unidentified": tokens_unidentified,
	})
}
