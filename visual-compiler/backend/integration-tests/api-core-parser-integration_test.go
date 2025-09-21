package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func registerNoInputUser(t *testing.T) {

	user_data := map[string]string{
		"username": "no_input_user",
		"email":    "no_user@gmail.com",
		"password": "noInput1234",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/register", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("registration failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("registration failed: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)

		if err != nil {
			t.Errorf("registration failed: %v", err)
		}
	}
}

func getNoInputUserId(t *testing.T) {

	user_data := map[string]string{
		"login":    "no_input_user",
		"password": "noInput1234",
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
	} else {

		body_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)

		t.Logf("Login working: %s", body_array["message"])
		no_input_user = body_array["id"]
	}
}

func deleteNoInputUser(t *testing.T) {

	user_data := map[string]string{
		"users_id": no_input_user,
	}
	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.NewRequest("DELETE",
		"http://localhost:8080/api/users/delete",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("User deletion failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(response.Body)
		t.Errorf("User deletion failed: %s", string(body_bytes))
	} else {

		body_bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)

		t.Logf("Deletion working: %s", body_array["message"])
	}
}

func TestReadGrammar_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)

	data := map[string]interface{}{
		"users_id":  test_user_id,
		"variables": []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		"terminals": []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		"start":     "RANDOM",
		"rules": []map[string]interface{}{
			{"input": "STATEMENT", "output": []string{"DECLARATION", "SEPARATOR"}},
			{"input": "DECLARATION", "output": []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{"input": "EXPRESSION", "output": []string{"TERM", "OPERATOR", "TERM"}},
			{"input": "TERM", "output": []string{"INTEGER"}},
			{"input": "TYPE", "output": []string{"KEYWORD"}},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/grammar", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusBadRequest {
		t.Logf("ReadGrammar Error: success")
	}
}

func TestReadGrammar_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	getID(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"variables":    []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		"terminals":    []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		"start":        "STATEMENT",
		"rules": []map[string]interface{}{
			{"input": "STATEMENT", "output": []string{"DECLARATION", "SEPARATOR"}},
			{"input": "DECLARATION", "output": []string{"TYPE", "IDENTIFIER", "ASSIGNMENT", "EXPRESSION"}},
			{"input": "EXPRESSION", "output": []string{"TERM", "OPERATOR", "TERM"}},
			{"input": "TERM", "output": []string{"INTEGER"}},
			{"input": "TYPE", "output": []string{"KEYWORD"}},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/grammar", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["message"] != "Grammar successfully inserted. Ready to create Syntax Tree" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestReadGrammar_Valid2(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"variables":    []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		"terminals":    []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "PUNCTUATION", "NUMBER"},
		"start":        "STATEMENT",
		"rules": []map[string]interface{}{
			{"input": "STATEMENT", "output": []string{"DECLARATION", "PUNCTUATION"}},
			{"input": "DECLARATION", "output": []string{"TYPE", "IDENTIFIER", "OPERATOR", "EXPRESSION"}},
			{"input": "EXPRESSION", "output": []string{"TERM"}},
			{"input": "TERM", "output": []string{"NUMBER"}},
			{"input": "TYPE", "output": []string{"KEYWORD"}},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/grammar", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["message"] != "Grammar successfully inserted. Ready to create Syntax Tree" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestCreateSyntaxTree_NoTokens(t *testing.T) {

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
		"http://localhost:8080/api/parsing/tree", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Tokens code not found. Please go back to lexing" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func addProjectToNoInputUser(t *testing.T) {

	getNoInputUserId(t)
	project_name = "project1"
	user_data := map[string]string{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/save", "application/json",
		bytes.NewBuffer(req),
	)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("saving project name failed: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)

		if err != nil {
			t.Errorf("saving project name failed: %v", err)
		}
	}
}

func addSourceCodeToNoInputUser(t *testing.T) {

	getNoInputUserId(t)
	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     no_input_user,
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

func readRulesToNoInputUser(t *testing.T) {

	getNoInputUserId(t)
	data := map[string]interface{}{
		"users_id":     no_input_user,
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

func tokeniseNoInputUser(t *testing.T) {

	getNoInputUserId(t)
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

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestCreateSyntaxTree_NoGrammarCode(t *testing.T) {

	server := startServerCore(t)
	defer closeServerCore(t, server)

	addProjectToNoInputUser(t)
	addSourceCodeToNoInputUser(t)
	readRulesToNoInputUser(t)
	tokeniseNoInputUser(t)

	data := map[string]interface{}{
		"users_id":     no_input_user,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/tree", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Grammar code not found. Please create one" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestCreateSyntaxTree_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	getID(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/tree", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Syntax Tree creation failed" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestCreateSyntaxTree_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/tree", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Parser not working")
	}
}

func TestTreeToString_NoSyntaxTree(t *testing.T) {
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
		"http://localhost:8080/api/parsing/treeString", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Syntax tree not found. Please create one" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestTreeToString_Valid(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	loginUser(t)
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/parsing/treeString", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Parser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["message"] != "Successfully generated Syntax Tree into a string" {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}
