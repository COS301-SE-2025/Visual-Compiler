package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AnalyseUserInputs struct {
	// Scope rules for the variables/functions scopes
	ScopeRules []*services.ScopeRule `json:"scope_rules" binding:"required"`
	// Grammar rules to be analysed
	GrammarRules services.GrammarRules `json:"grammar_rules" binding:"required"`
	// Type rules for Type Checking
	TypeRules []services.TypeRule `json:"type_rules" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// @Summary Analysing phase
// @Description Accepts scope rules, grammar rules and type rules from the user. Searches the database for the syntax tree created from the user. If it exists, the analysing process is performed and the artefacts are stored in the database
// @Tags Analysing
// @Accept json
// @Produce json
// @Param request body AnalyseUserInputs true "Read Analysing Inputs From User"
// @Success 200 {object} map[string]string "Artefacts Successfully stored"
// @Failure 400 {object} map[string]string "Invalid input or artefacts failed to insert"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Syntax Tree not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /analysing/analyse [post]
func Analyse(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AnalyseUserInputs

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	parsing_collection := mongo_cli.Database("visual-compiler").Collection("parsing")
	analyse_collection := mongo_cli.Database("visual-compiler").Collection("analysing")

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

	var parsing_res struct {
		Tree services.SyntaxTree `bson:"tree"`
	}

	err = parsing_collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}).Decode(&parsing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tree not found. Please go back to parsing"})
		return
	}

	artefact, _, err := services.Analyse(req.ScopeRules, parsing_res.Tree, req.GrammarRules, req.TypeRules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Artefacts creation failed", "details": err.Error()})
		return
	}

	filters := bson.M{"users_id": dbUser.UsersID, "project_name": req.Project_Name}
	var userexisting bson.M

	err = analyse_collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = analyse_collection.InsertOne(ctx, bson.M{
			"users_id":              dbUser.UsersID,
			"symbol_table_artefact": artefact,
			"project_name":          req.Project_Name,
			"scope_rules":           req.ScopeRules,
			"grammar_rules":         req.GrammarRules,
			"type_rules":            req.TypeRules,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$set", Value: bson.M{
				"symbol_table_artefact": artefact,
				"scope_rules":           req.ScopeRules,
				"grammar_rules":         req.GrammarRules,
				"type_rules":            req.TypeRules,
			}},
		}
		_, err = analyse_collection.UpdateOne(ctx, filters, update_existing)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Update error"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database lookup error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Symbol Table Artefact successfully inserted.",
		"symbol_table": artefact,
	})
}
