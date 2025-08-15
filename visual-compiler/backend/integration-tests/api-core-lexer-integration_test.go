package tests

import (
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
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

func TestStoreSourceCode_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(re_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Code is ready for further processing"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestStoreSourceCode_Valid2(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)
	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(re_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Code is ready for further processing"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"pairs": []map[string]string{
			{
				"Type":  "KEYWORD",
				"Regex": "\\b(if|else|int)\\b",
			},
			{
				"Type":  "IDENTIFIER",
				"Regex": "[a-zA-Z_]\\w*",
			},
			{
				"Type":  "NUMBER",
				"Regex": "\\d+(\\.\\d+)?",
			},
			{
				"Type":  "OPERATOR",
				"Regex": "=",
			},
			{
				"Type":  "PUNCTUATION",
				"Regex": ";",
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/rules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Rules successfully created."}` {
			t.Logf("CreateRulesFromCode: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_Valid2(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"pairs": []map[string]string{
			{
				"Type":  "KEYWORD",
				"Regex": "\\b(if|else|int)\\b",
			},
			{
				"Type":  "IDENTIFIER",
				"Regex": "[a-zA-Z_]\\w*",
			},
			{
				"Type":  "NUMBER",
				"Regex": "\\d+(\\.\\d+)?",
			},
			{
				"Type":  "OPERATOR",
				"Regex": "=",
			},
			{
				"Type":  "PUNCTUATION",
				"Regex": ";",
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/rules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Rules successfully created."}` {
			t.Logf("CreateRulesFromCode: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	getNoInputUserId(t)
	deleteNoInputUser(t)
	registerNoInputUser(t)
	getNoInputUserId(t)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
		"pairs": []map[string]string{
			{
				"Type":  "KEYWORD",
				"Regex": "[a-z",
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/rules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"error":"Regex rule creation failed"}{"error":"Source code not found. Please enter a source code"}` {
			t.Logf("CreateRulesFromCode Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_Val(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
		"pairs":        []map[string]string{},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/rules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Rules successfully created."}` {
			t.Logf("CreateRulesFromCode: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestLexing_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	reg_expr_data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(reg_expr_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestLexing_Valid2(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)
	reg_expr_data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(reg_expr_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestLexing_NoSourceCode(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	reg_expr_data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(reg_expr_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}
	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Source code not found" {
			t.Logf("Lexing Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestLexing_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}
	code_req, err := json.Marshal(re_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}
	code_res, err := http.Post(
		"http://localhost:8080/api/lexing/code", "application/json",
		bytes.NewBuffer(code_req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer code_res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}
	if code_res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(code_res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	reg_expr_data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(reg_expr_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}
	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Tokenization failed" {
			t.Logf("Lexing Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	} else {
		t.Errorf("Error expected")
	}

}

func TestReadDFAFromUser_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"project_name": project_name,
		"states": []string{
			"START",
			"S1",
			"S2",
			"S3",
			"S4",
			"S5",
		},
		"transitions": []map[string]string{
			{
				"from":  "START",
				"to":    "S1",
				"label": "i",
			},
			{
				"from":  "S1",
				"to":    "S5",
				"label": "n",
			},
			{
				"from":  "S5",
				"to":    "S4",
				"label": "t",
			},
			{
				"from":  "START",
				"to":    "S2",
				"label": "0123456789",
			},
			{
				"from":  "S2",
				"to":    "S2",
				"label": "0123456789",
			},
			{
				"from":  "START",
				"to":    "S3",
				"label": "abcdefghijklmnopqrstuvwxyz",
			},
			{
				"from":  "S3",
				"to":    "S3",
				"label": "abcdefghijklmnopqrstuvwxyz0123456789",
			},
		},
		"start_state": "START",
		"accepting_states": []map[string]string{
			{
				"state":      "S3",
				"token_type": "IDENTIFIER",
			},
			{
				"state":      "S4",
				"token_type": "KEYWORD",
			},
			{
				"state":      "S2",
				"token_type": "NUMBER",
			},
		},
		"users_id": test_user_id,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfa", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"DFA successfuly created. Ready to create tokens"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestReadDFAFromUser_TestUser(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"project_name": project_name,
		"states": []string{
			"START",
			"S1",
			"S2",
			"S3",
			"S4",
			"S5",
		},
		"transitions": []map[string]string{},
		"start_state": "START",
		"accepting_states": []map[string]string{
			{
				"state":      "S3",
				"token_type": "IDENTIFIER",
			},
			{
				"state":      "S4",
				"token_type": "KEYWORD",
			},
			{
				"state":      "S2",
				"token_type": "NUMBER",
			},
		},
		"users_id": "6834279e18addf82669c9acd",
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfa", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"DFA successfuly created. Ready to create tokens"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestTokensFromDFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	//loginUser(t)
	get_id(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToTokens", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Successfully tokenised your code","tokens":[{"type":"KEYWORD","value":"int"},{"type":"IDENTIFIER","value":"x"},{"type":"NUMBER","value":"2"}],"tokens_unidentified":["=",";"]}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestTokensFromDFA_NoSourceCode(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToTokens", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Source code not found" {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestTokensFromDFA_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToTokens", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Tokenization from DFA failed" {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertDFAToRG_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToRegex", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["message"] == "Successfully converted DFA to Regex" {
			t.Logf("ConvertDFAToRG Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertDFAToRG_NoDFA(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToRegex", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "DFA not found. Please create one" {
			t.Logf("ConvertDFAToRG Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}
}

func TestConvertDFAToRG_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/dfaToRegex", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Conversion from DFA to Regex failed" {
			t.Logf("ConverDFAToRG error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}
}

func TestConvertRGToNFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToNFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["message"] == "Successfully converted Regex to NFA" {
			t.Logf("ConvertRGToNFA Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToNFA_NoRegexRules(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToNFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Regex rules not found. Please create one" {
			t.Logf("ConvertRGToNFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToNFA_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToNFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Conversion from Regex to NFA failed" {
			t.Logf("ConvertRGToNFA Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToDFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["message"] == "Successfully converted Regex to DFA" {
			t.Logf("ConvertRGToDFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToDFA_NoRegexRules(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Regex rules not found. Please create one" {
			t.Logf("ConvertRGToDFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToDFA_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/regexToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Conversion from Regex to DFA failed" {
			t.Logf("ConvertRGToDFA Error: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertNFAToDFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/nfaToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["message"] == "Successfully converted NFA to DFA" {
			t.Logf("ConvertNFAToDFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertNFAToDFA_NoNFA(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/nfaToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "NFA not found. Please create one" {
			t.Logf("ConvertNFAToDFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertNFAToDFA_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     "6834279e18addf82669c9acd",
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/lexing/nfaToDFA", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Conversion from NFA to DFA failed" {
			t.Logf("ConvertNFAToDFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}
