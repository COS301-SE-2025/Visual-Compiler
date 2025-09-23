package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var test_user_id string
var project_name string
var new_project_name string
var no_input_project_name string
var auth0_id string

var server *http.Server

func startServerCore(t *testing.T) *http.Server {
	os.Setenv("test_mode", "true")

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "https://visual-compiler.co.za"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api", routers.Auth0MiddleWare())

	lexing_group := api.Group("/lexing")
	routers.SetupLexingRouter(lexing_group)

	user_router := routers.SetupUserRouter()
	router.Any("/api/users/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		user_router.HandleContext(c)
	})

	parsing_group := api.Group("/parsing")
	routers.SetupParsingRouter(parsing_group)

	analysing_group := api.Group("/analysing")
	routers.SetupAnalysingRouter(analysing_group)

	translator_group := api.Group("/translating")
	routers.SetupTranslatorRouter(translator_group)

	optimising_group := api.Group("/optimising")
	routers.SetupOptimisingRouter(optimising_group)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	//goroutine needed so test does not enter an infinite loop
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Server failed to start: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	return server
}

func closeServerCore(t *testing.T, server *http.Server) {
	cont, cancel_cont := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel_cont()
	server.Shutdown(cont)
}

func loginTestUser(t *testing.T) {
	user_data := map[string]string{
		"login":    "halfstack.testuser@gmail.com",
		"password": "testUser13",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Fatalf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/login", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Fatalf("Login failed: %s", string(body_bytes))
	}

	body_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	var body_array map[string]string
	_ = json.Unmarshal(body_bytes, &body_array)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("Login working: %s", body_array["message"])
	test_user_id = body_array["id"]
	auth0_id = body_array["auth_token"]
	project_name = "project1"

	new_project_name = "project2"
	addProject(t)

}

func addNewProject(t *testing.T) {
	user_data := map[string]string{
		"project_name": new_project_name,
		"users_id":     test_user_id,
		"pipeline":     "",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Fatalf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/save", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Fatalf("New project failed: %v", err)
	}

	defer res.Body.Close()

}
func addProject(t *testing.T) {
	no_input_project_name = fmt.Sprintf("project-%s", time.Now().Format("20060102-150405"))
	user_data := map[string]string{
		"project_name": new_project_name,
		"users_id":     test_user_id,
		"pipeline":     "",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Fatalf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/save", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Fatalf("New project failed: %v", err)
	}

	defer res.Body.Close()

}

func deleteNewProject(t *testing.T) {
	user_data := map[string]string{
		"project_name": new_project_name,
		"users_id":     test_user_id,
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Fatalf("converting data to json failed")
	}

	res, err := http.NewRequest("DELETE",
		"http://localhost:8080/api/users/deleteProject",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Project deletion failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Fatalf("Delete Project failed: %s", string(body_bytes))
	}

}

func deleteProject(t *testing.T) {
	user_data := map[string]string{
		"project_name": no_input_project_name,
		"users_id":     test_user_id,
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Fatalf("converting data to json failed")
	}

	res, err := http.NewRequest("DELETE",
		"http://localhost:8080/api/users/deleteProject",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Project deletion failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

}
