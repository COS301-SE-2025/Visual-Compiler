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

type TranslatorRules struct {
	// User's ID for searching and storing purposes
	UsersID bson.ObjectID `json:"users_id" binding:"required" example:"685df259c1294de5546b045f"`
	// Translation Rules
	Rules []services.TranslationRule `json:"translation_rules" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
}

// @Summary Create translation rules
// @Description Takes the user's rules, reads it, verifies it and stores it with the user's ID. If translation rules already exist, it overwrites those rules and removes any other created fields (translation)
// @Tags Translating
// @Accept json
// @Produce json
// @Param request body TranslatorRules true "Read Translation Rules From User"
// @Success 200 {object} map[string]string "Translation Rules successfully stored"
// @Failure 400 {object} map[string]string "Invalid input/Translation rules failed to insert"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /translating/rules [post]
func ReadRules(c *gin.Context) {
	var req TranslatorRules

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	json_as_bytes, err := json.Marshal(req.Rules)
	if err != nil {
		panic(err)
	}

	rules, err := services.ReadTranslationRules(json_as_bytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Regex rule creation failed"})
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("translating")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}
	var userexisting bson.M

	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"translating_rules": rules,
			"users_id":          req.UsersID,
			"project_name":      req.Project_Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"translation": "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"translating_rules": rules,
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

	c.JSON(http.StatusOK, gin.H{"message": " Ready for further processing"})
}

// @Summary Translates syntax tree using the translation rules into code
// @Description Searches the database for the user's Syntax tree and Translation Rules. If found, both are used to translate the tree into code. The code is either created or ,if already existing, updated. If the syntax tree and translation rules are not found, returns an error
// @Tags Translating
// @Accept json
// @Produce json
// @Param request body IDRequest true "Create Code from Syntax Tree and Translation Rules"
// @Success 200 {object} map[string]string "Code successfully created and stored"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 404 {object} map[string]string "Tree not found/Rules not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /translating/translate [post]
func TranslateCode(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	parsing_collection := mongo_cli.Database("visual-compiler").Collection("parsing")
	translating_collection := mongo_cli.Database("visual-compiler").Collection("translating")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var parsing_res struct {
		Tree services.SyntaxTree `bson:"tree"`
	}

	var translating_res struct {
		Rules []services.TranslationRule `bson:"translating_rules"`
	}

	err := parsing_collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&parsing_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tree not found. Please go back to parsing"})
		return
	}

	err = translating_collection.FindOne(ctx, bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}).Decode(&translating_res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Translating rules not found. Please go back to Translation"})
		return
	}

	translated_code, err := services.Translate(parsing_res.Tree, translating_res.Rules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code creation failed", "details": err.Error()})
		return
	}

	filters := bson.M{"users_id": req.UsersID, "project_name": req.Project_Name}
	update_users_translating := bson.M{"$set": bson.M{
		"code": translated_code,
	}}

	_, err = translating_collection.UpdateOne(ctx, filters, update_users_translating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert NFA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully translated your code",
		"code":    translated_code,
	})
}
