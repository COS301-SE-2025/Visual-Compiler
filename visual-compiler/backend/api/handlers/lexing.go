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

type LexingRequest struct {
	Code  string               `json:"source_code" bson:"required"`
	Pairs []services.TypeRegex `json:"pairs" bson:"required"`
}

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
	unexpected_tokens := services.GetUnexpectedTokens()

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
		"unexpected_tokens": unexpected_tokens,
	})
}
