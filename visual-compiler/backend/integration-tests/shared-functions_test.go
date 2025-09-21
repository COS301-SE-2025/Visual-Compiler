package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var test_user_id string
var no_input_user string
var project_name string

func startServerCore(t *testing.T) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	no_input_user = "689e5c3b2b0a249a86761244"
	project_name = "project1"

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	lexer_router := routers.SetupLexingRouter()
	router.Any("/api/lexing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		lexer_router.HandleContext(c)
	})

	user_router := routers.SetupUserRouter()
	router.Any("/api/users/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		user_router.HandleContext(c)
	})

	parser_router := routers.SetupParsingRouter()
	router.Any("/api/parsing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		parser_router.HandleContext(c)
	})

	analyser_router := routers.SetupAnalysingRouter()
	router.Any("/api/analysing/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		analyser_router.HandleContext(c)
	})

	translator_router := routers.SetupTranslatorRouter()
	router.Any("/api/translating/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		translator_router.HandleContext(c)
	})

	optimiser_router := routers.SetupOptimisingRouter()
	router.Any("/api/optimising/*any", func(c *gin.Context) {
		c.Request.URL.Path = c.Param("any")
		optimiser_router.HandleContext(c)
	})

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

func getID(t *testing.T) {
	user_data := map[string]string{
		"login":    "tiaharripersad",
		"password": "tia1234$$",
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
	project_name = "project1"

}

func generateRandomProjectName() {
	project_name = fmt.Sprintf("project-%s", time.Now().Format("20060102-150405"))
}
