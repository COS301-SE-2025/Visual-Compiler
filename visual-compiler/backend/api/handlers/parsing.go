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
	// User's defined variables
	Vars []string `json:"variables" binding:"required" example:"S, Decl"`
	// User's defined terminal variables
	Terminals []string `json:"terminals" binding:"required" example:"KEYWORD, IDENTIFIER, OPERATOR, NUMBER, PUNCTUATION"`
	// User's defined start variable
	StartVar string `json:"start" binding:"required" example:"S"`
	// User's defined rules
	Rules []services.ParsingRule `json:"rules" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// @Summary Processs and store user-defined grammer
// @Description Accepts grammar variables, terminals, start variable, and rules from the user and stores them in the database. If it already exists, it updates the current grammar
// @Tags Parsing
// @Accept json
// @Produce json
// @Param request body ReadGrammerFromUser true "Read Grammer From User"
// @Success 200 {object} map[string]string "Grammar successfully read and stored/updated"
// @Failure 400 {object} map[string]string "Invalid input or Grammar failed to insert"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /parsing/grammar [post]
func ReadGrammar(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ReadGrammerFromUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	collection := mongo_cli.Database("visual-compiler").Collection("parsing")

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

	filters := bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}
	var userexisting bson.M

	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"users_id":     dbUser.UsersID,
			"grammar":      grammar,
			"project_name": req.Project_Name,
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

// @Summary Create and store syntax tree from stored grammar and tokens
// @Description Searches database for Grammar and Tokens. If found, and creates and stores the tree.
// @Tags Parsing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Syntax tree successfully created and stored/updated"
// @Failure 400 {object} map[string]string "Invalid input or Syntax Tree failed to insert"
// @Failure 404 {object} map[string]string "Tokens or Grammer not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /parsing/tree [post]
func CreateSyntaxTree(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ProjectNameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	lexing_collection := mongo_cli.Database("visual-compiler").Collection("lexing")
	parsing_collection := mongo_cli.Database("visual-compiler").Collection("parsing")

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

	var lexing_res struct {
		Tokens []services.TypeValue `bson:"tokens"`
	}

	var parsing_res struct {
		Grammar services.Grammar `bson:"grammar"`
	}

	err = lexing_collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}).Decode(&lexing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tokens code not found. Please go back to lexing"})
		return
	}

	err = parsing_collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}).Decode(&parsing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Grammar code not found. Please create one"})
		return
	}

	tree, err := services.CreateSyntaxTree(lexing_res.Tokens, parsing_res.Grammar)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Syntax Tree creation failed", "details": err.Error()})
		return
	}

	filters := bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}
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

// @Summary Create and store syntax tree as a string from stored tree
// @Description Searches database for an existing syntax tree. If found, and creates and stores the tree as a string.
// @Tags Parsing
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Syntax tree String successfully created and stored/updated"
// @Failure 400 {object} map[string]string "Invalid input or Syntax Tree String failed to insert"
// @Failure 404 {object} map[string]string "Syntax Tree not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /parsing/treeString [post]
func TreeToString(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req ProjectNameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	collection := mongo_cli.Database("visual-compiler").Collection("parsing")

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

	var res struct {
		Tree services.SyntaxTree `bson:"tree"`
	}

	err = collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Syntax tree not found. Please create one"})
		return
	}

	tree_as_string := services.ConvertTreeToString(res.Tree.Root, "", true)

	filters := bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}
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

// @Summary Get user's syntax tree
// @Description Searches the database for the user's syntax tree
// @Tags Lexing
// @Accept json
// @Produce json
// @Param project_name query string true "Project Name"
// @Success 200 {object} map[string]string "Syntax Tree retrieved"
// @Failure 400 {object} map[string]string "Invalid input/Response failed"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /lexing/getTree [get]
func GetTree(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	project_name := c.Query("project_name")

	if project_name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: Missing query parameters."})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	collection := mongo_cli.Database("visual-compiler").Collection("parsing")

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

	var res struct {
		Tree services.SyntaxTree `bson:"tree"`
	}

	err = collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID, "project_name": project_name}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Syntax Tree retrieved",
		"tree":    res.Tree,
	})
}
