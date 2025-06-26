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
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReadGrammerFromUser struct {
	UsersID   bson.ObjectID          `json:"users_id" binding:"required"`
	Vars      []string               `json:"variables" binding:"required"`
	Terminals []string               `json:"terminals" binding:"required"`
	StartVar  string                 `json:"start" binding:"required"`
	Rules     []services.ParsingRule `json:"rules" binding:"required"`
}

func ReadGrammar(c *gin.Context) {
	var req ReadGrammerFromUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("parsing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	users_grammer_rules := services.Grammar{
		Variables: req.Vars,
		Terminals: req.Terminals,
		Start:     req.StartVar,
		Rules:     req.Rules,
	}

	json_as_bytes, err := json.Marshal(users_grammer_rules)
	if err != nil {
		panic(err)
	}

	grammar, err := services.ReadGrammar(json_as_bytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Grammar creation failed", "details": err.Error()})
	}

	filters := bson.M{"users_id": req.UsersID}
	var userexisting bson.M

	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"users_id": req.UsersID,
			"grammar":  grammar,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"tree":        "",
				"tree_string": "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"grammar": grammar,
			}},
		}
		_, err = collection.UpdateOne(ctx, filters, update_existing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database lookup error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Grammar successfully inserted. Ready to create Syntax Tree"})
}

func CreateSyntaxTree(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	lexing_collection := mongo_cli.Database("visual-compiler").Collection("lexing")
	parsing_collection := mongo_cli.Database("visual-compiler").Collection("parsing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var lexing_res struct {
		Tokens []services.TypeValue `bson:"tokens"`
	}

	var parsing_res struct {
		Grammar services.Grammar `bson:"grammar"`
	}

	err := lexing_collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&lexing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tokens code not found. Please go back to lexing"})
		return
	}

	err = parsing_collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&parsing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grammar code not found. Please create one"})
		return
	}

	tree, err := services.CreateSyntaxTree(lexing_res.Tokens, parsing_res.Grammar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Syntax Tree creation failed", "details": err.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID}
	update_users_tree := bson.M{"$set": bson.M{
		"tree": tree,
	}}

	_, err = parsing_collection.UpdateOne(ctx, filters, update_users_tree)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully created Syntax Tree",
		"tree":    tree,
	})
}

func TreeToString(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("parsing")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Tree services.SyntaxTree `bson:"tree"`
	}

	err := collection.FindOne(ctx, bson.M{"users_id": req.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Syntax tree not found. Please create one"})
		return
	}

	tree_as_string := services.ConvertTreeToString(res.Tree.Root, "", true)

	filters := bson.M{"users_id": req.UsersID}
	update_users_tree_string := bson.M{"$set": bson.M{
		"tree_string": tree_as_string,
	}}

	_, err = collection.UpdateOne(ctx, filters, update_users_tree_string)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Successfully generated Syntax Tree into a string",
		"tree_string": tree_as_string,
	})
}
