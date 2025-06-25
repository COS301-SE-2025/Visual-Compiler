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

func startServerCore(t *testing.T) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/login", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Login failed: %s", string(body_bytes))
	}

	body_bytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	var body_array map[string]string
	err = json.Unmarshal(body_bytes, &body_array)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("Login working: %s", body_array["message"])
	test_user_id = body_array["id"]

}

func TestStoreSourceCode_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	re_data := map[string]interface{}{
		"source_code": "int x = 2 ;",
		"users_id":    test_user_id,
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
		"users_id": test_user_id,
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
			t.Logf("ReadDFAFromUser: success")
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
		"users_id": test_user_id,
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

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

	data := map[string]string{
		"users_id": test_user_id,
	}
	lexer_req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	lexer_res, err := http.Post(
		"http://localhost:8080/api/lexing/lexer", "application/json",
		bytes.NewBuffer(lexer_req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer lexer_res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(lexer_res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}
	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(lexer_res.Body)
		if string(body_bytes) == `{"message":"Successfully tokenised your code","tokens":[{"type":"KEYWORD","value":"int"},{"type":"IDENTIFIER","value":"x"},{"type":"OPERATOR","value":"="},{"type":"NUMBER","value":"2"},{"type":"PUNCTUATION","value":";"}],"tokens_unidentified":[],"users_id":"682ef437ee52c95438ba65c1"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestReadDFAFromUser_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
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

func TestTokensFromDFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id": test_user_id,
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

func TestConvertDFAToRG_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id": test_user_id,
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
		if string(body_bytes) == `{"message":"Successfully converted DFA to Regex","rules":[{"type":"IDENTIFIER","regex":"[a-zA-Z_]\\w*"},{"type":"KEYWORD","regex":"\\bint\\b"},{"type":"NUMBER","regex":"\\d+"}]}` {
			t.Logf("ConvertDFAToRG: success")
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
		"users_id": test_user_id,
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
			t.Logf("ConvertRGToNFA: success")
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
		"users_id": test_user_id,
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

func TestConvertNFAToDFA_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	get_id(t)
	data := map[string]interface{}{
		"users_id": test_user_id,
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
