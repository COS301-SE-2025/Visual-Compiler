package tests

import (
	"testing"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"bytes"
	"encoding/json"
	"net/http"
	"io"
	"context"
	"errors"
	"time"
)

var test_user_id string

func startServerCore(t *testing.T) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Attach CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:	  []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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

func get_id(t *testing.T) {
	user_data :=map[string]string{
		"login": "tiaharripersad",
		"password": "tia1234$$",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/login", "application/json",
		bytes.NewBuffer(req),
	)
	if err!=nil {
		t.Errorf("Login failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Errorf("Login failed: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	var respMap map[string]string
	err = json.Unmarshal(bodyBytes, &respMap)

	t.Logf("Login working: %s", respMap["message"])
	test_user_id = respMap["id"]

}

func TestLexerCode(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	re_data :=map[string]interface{}{
		"id": test_user_id,
		"source_code": "int x = 2 ;",
		"pairs": []map[string]string{
			{
				"Type": "KEYWORD",
				"Regex": "\\b(if|else|int)\\b",
			},
			{
				"Type": "IDENTIFIER",
				"Regex": "[a-zA-Z_]\\w*",
			},
			{
				"Type": "NUMBER",
				"Regex": "\\d+(\\.\\d+)?",
			},
			{
				"Type": "OPERATOR",
				"Regex": "=",
			},
			{
				"Type": "PUNCTUATION",
				"Regex": ";",
			},
		},
	}

	req, err := json.Marshal(re_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/code", "application/json",
		bytes.NewBuffer(req),
	)
	defer res.Body.Close()

	if err!= nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(bodyBytes))
	}

	if res.StatusCode == http.StatusOK{
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Logf("Lexer working: %s", string(bodyBytes))
	}

}

func TestLexer(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	reg_expr_data :=map[string]interface{}{
		"id": test_user_id,
		"source_code": "int x = 2 ;",
		"pairs": []map[string]string{
			{
				"Type": "KEYWORD",
				"Regex": "\\b(if|else|int)\\b",
			},
			{
				"Type": "IDENTIFIER",
				"Regex": "[a-zA-Z_]\\w*",
			},
			{
				"Type": "NUMBER",
				"Regex": "\\d+(\\.\\d+)?",
			},
			{
				"Type": "OPERATOR",
				"Regex": "=",
			},
			{
				"Type": "PUNCTUATION",
				"Regex": ";",
			},
		},
	}

	req, err := json.Marshal(reg_expr_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/code", "application/json",
		bytes.NewBuffer(req),
	)
	defer res.Body.Close()

	if err!= nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(bodyBytes))
	}

	
	data := map[string]string{
		"id":test_user_id,
	}
	lexer_req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	lexer_res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(lexer_req),
	)
	defer lexer_res.Body.Close()

	if err!= nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(lexer_res.Body)
		t.Errorf("Lexer not working: %s", string(bodyBytes))
	}
	if res.StatusCode == http.StatusOK{
		bodyBytes, _ := io.ReadAll(lexer_res.Body)
		t.Logf("Lexer working: %s", string(bodyBytes))
	}

}