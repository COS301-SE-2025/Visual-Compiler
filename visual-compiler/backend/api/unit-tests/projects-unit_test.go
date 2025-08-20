package unit_tests

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestSaveProjectName_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/users/save", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.SaveProjectName(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
	}

}

func TestSaveProjectPipeline_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/users/savePipeline", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.SaveProjectPipeline(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
	}

}

func TestGetAllProjects_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/users/getProjects", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.GetAllProjects(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
	}

}

func TestGetProject_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/users/getProject", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.GetProject(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
	}

}

func TestDeleteProject_InvalidInput(t *testing.T) {
	gin.SetMode(gin.TestMode)
	contxt, rec := createPhaseTestContext(t)

	res, err := http.NewRequest("POST", "/api/users/deleteProject", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.DeleteProject(contxt)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("BadRequest status code expected")
	}

}
