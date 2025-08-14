package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/core/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProjectRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
	// User's project name
	Project_Name string `json:"project_name" binding:"required"`
	// User's saved pipeline
	Pipeline map[string]any `json:"pipeline" binding:"required"`
}

type IDOnlyRequest struct {
	// Represents the User's ID from frontend
	UsersID bson.ObjectID `json:"users_id" binding:"required"`
}

// @Summary Save project name
// @Description Saves the project name in the Users table
// @Tags Users
// @Accept json
// @Produce json
// @Param request body IDRequest true "Save User's Defined Project Name"
// @Success 200 {object} map[string]interface{} "Project added successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 409 {object} map[string]string "Project name already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/save [post]
func SaveProjectName(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{
		"_id":           req.UsersID,
		"projects.name": req.Project_Name,
	}

	err := collection.FindOne(ctx, filters).Err()
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Project name already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	new_project := bson.M{"name": req.Project_Name}
	update_project := bson.M{"$push": bson.M{"projects": new_project}}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": req.UsersID}, update_project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project added successfully",
	})
}

// @Summary Save project pipeline from frontend
// @Description Saves the project's pipeline to be shown when the user opens that project on the frontend
// @Tags Users
// @Accept json
// @Produce json
// @Param request body ProjectRequest true "Save User's Pipeline"
// @Success 200 {object} map[string]interface{} "Pipeline saved"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/savePipeline [post]
func SaveProjectPipeline(c *gin.Context) {
	var req ProjectRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{
		"_id":           req.UsersID,
		"projects.name": req.Project_Name,
	}

	err := collection.FindOne(ctx, filters).Err()

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	update_project := bson.M{
		"$set": bson.M{"projects.$.pipeline": req.Pipeline},
	}

	_, err = collection.UpdateOne(ctx, filters, update_project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save pipeline: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Pipeline successfully saved",
		"pipeline": req.Pipeline,
	})
}

// @Summary Get all user's projects
// @Description Gets all the user's project names
// @Tags Users
// @Accept json
// @Produce json
// @Param users_id query string true "User's ID"
// @Success 200 {object} map[string]interface{} "Retrieved all project names"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/getProjects [get]
func GetAllProjects(c *gin.Context) {
	users_id := c.Query("users_id")

	if users_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: Missing query parameters."})
		return
	}

	mongo_cli := db.ConnectClient()
	collection := mongo_cli.Database("visual-compiler").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var res struct {
		Projects []struct {
			Name string `bson:"name"`
		} `bson:"projects"`
	}

	id, err := bson.ObjectIDFromHex(users_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
		return
	}

	err = collection.FindOne(
		ctx,
		bson.M{"_id": id},
		options.FindOne().SetProjection(bson.M{"projects.name": 1, "_id": 0}),
	).Decode(&res)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	var project_names []string
	for _, proj := range res.Projects {
		project_names = append(project_names, proj.Name)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Retrieved all projects",
		"all_projects": project_names,
	})
}

// @Summary Delete a project
// @Description Deletes a project and all associated fields in every collection in the database
// @Tags Users
// @Accept json
// @Produce json
// @Param request body IDRequest true "Delete Project"
// @Success 200 {object} map[string]interface{} "Successfully deleted project"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to delete project"
// @Router /users/deleteProject [delete]
func DeleteProject(c *gin.Context) {
	var req IDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: " + err.Error()})
		return
	}

	mongo_cli := db.ConnectClient()
	users_collection := mongo_cli.Database("visual-compiler").Collection("users")
	lexing_collection := mongo_cli.Database("visual-compiler").Collection("lexing")
	parsing_collection := mongo_cli.Database("visual-compiler").Collection("parsing")
	analysing_collection := mongo_cli.Database("visual-compiler").Collection("analysing")
	translating_collection := mongo_cli.Database("visual-compiler").Collection("translating")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := bson.M{
		"_id":           req.UsersID,
		"projects.name": req.Project_Name,
	}

	if err := users_collection.FindOne(ctx, filters).Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	update_project := bson.M{
		"$pull": bson.M{
			"projects": bson.M{"name": req.Project_Name},
		},
	}

	_, err := users_collection.UpdateOne(ctx, bson.M{"_id": req.UsersID}, update_project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove project"})
		return
	}

	filters_for_collections := bson.M{
		"users_id":     req.UsersID,
		"project_name": req.Project_Name,
	}

	if _, err := lexing_collection.DeleteOne(ctx, filters_for_collections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if _, err := parsing_collection.DeleteOne(ctx, filters_for_collections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if _, err := analysing_collection.DeleteOne(ctx, filters_for_collections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if _, err := translating_collection.DeleteOne(ctx, filters_for_collections); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully deleted your projects",
	})
}

// @Summary Get a project
// @Description Gets all the details about a user's project in every collection in the database
// @Tags Users
// @Accept json
// @Produce json
// @Param project_name query string true "Name of project"
// @Param users_id query string true "User's ID"
// @Success 200 {object} map[string]interface{} "Retrieved all project details"
// @Failure 400 {object} map[string]string "Invalid input/Conversion failed"
// @Failure 500 {object} map[string]string "Error in reading from specific collection"
// @Router /users/getProject [get]
func GetProject(c *gin.Context) {
	project_name := c.Query("project_name")
	users_id := c.Query("users_id")

	if project_name == "" || users_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input is invalid: Missing query parameters."})
		return
	}

	mongo_cli := db.ConnectClient()
	database := "visual-compiler"

	db_collections := []string{"lexing", "parsing", "analysing", "translating"}
	res := make(map[string]any)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := bson.ObjectIDFromHex(users_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is invalid"})
		return
	}

	filters := bson.M{
		"users_id":     id,
		"project_name": project_name,
	}

	for _, col := range db_collections {
		var user_doc bson.M

		opts := options.FindOne().SetProjection(bson.M{
			"users_id":     0,
			"project_name": 0,
			"_id":          0,
		})

		err := mongo_cli.Database(database).Collection(col).FindOne(ctx, filters, opts).Decode(&user_doc)

		if err == mongo.ErrNoDocuments {
			continue
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error in reading %s: %v", col, err)})
			return
		}

		res[col] = user_doc
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Retrieved users project details",
		"results": res,
	})
}
