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

// @Summary Reads new code for optimising
// @Description Reads in new code from the user. Starts the optimising process
// @Tags Optimising
// @Accept json
// @Produce json
// @Param request body SourceCodeOnlyRequest true "Read optimising code from User"
// @Success 200 {object} map[string]string "Optimising code successfully and stored"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /optimising/source_code [post]
func StoreOptimisingCode(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req SourceCodeOnlyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	collection := mongo_cli.Database("visual-compiler").Collection("optimising")

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

	filters := bson.M{
		"users_id":     dbUser.UsersID,
		"project_name": req.Project_Name,
	}
	var userexisting bson.M

	err = collection.FindOne(ctx, filters).Decode(&userexisting)

	if err == mongo.ErrNoDocuments {
		_, err = collection.InsertOne(ctx, bson.M{
			"optimising_source_code": req.Code,
			"users_id":               dbUser.UsersID,
			"project_name":           req.Project_Name,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database Insertion error"})
			return
		}
	} else if err == nil {
		update_existing := bson.D{
			bson.E{Key: "$unset", Value: bson.M{
				"optimised_code": "",
			}},
			bson.E{Key: "$set", Value: bson.M{
				"optimising_source_code": req.Code,
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

	c.JSON(http.StatusOK, gin.H{"message": "Code is ready for further processing"})
}

type OptimiseCodeRequest struct {
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
	// Simplify any code before execution (3 + 4 becomes 7 before execution)
	ConstantFolding bool `json:"constant_folding"`
	// Remove any code that will never be executed
	DeadCode bool `json:"dead_code"`
	// Expand any loops
	LoopUnrolling bool `json:"loop_unrolling"`
}

// @Summary Optimises stored code
// @Description Optimises the user's stored code with parameters for optimisation
// @Tags Optimising
// @Accept json
// @Produce json
// @Param request body OptimiseCodeRequest true "Optimise code"
// @Success 200 {object} map[string]string "Code successfully optimised and stored"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /optimising/optimise [post]
func OptimiseCode(c *gin.Context) {
	authID, is_existing := c.Get("auth0_id")
	if !is_existing {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req OptimiseCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid", "details": err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	collection := mongo_cli.Database("visual-compiler").Collection("optimising")

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
		OptimisingSourceCode string `bson:"optimising_source_code"`
	}

	err = collection.FindOne(ctx, bson.M{"users_id": dbUser.UsersID}).Decode(&res)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Source code not found. Please enter a source code"})
		return
	}

	filters := bson.M{
		"users_id":     dbUser.UsersID,
		"project_name": req.Project_Name,
	}

	optimised_code, err := services.OptimiseGoCode(res.OptimisingSourceCode, req.ConstantFolding, req.DeadCode, req.LoopUnrolling)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Optimising failed", "details": err.Error()})
		return
	}

	update_users_optimising := bson.M{
		"$set": bson.M{
			"optimised_code": optimised_code,
		},
	}

	_, err = collection.UpdateOne(ctx, filters, update_users_optimising)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save optimised code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Successfully optimised code",
		"optimised_code": optimised_code,
	})
}
