package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestReadGrammar_CoreError(t *testing.T) {
	server = startServerCore(t)
	loginTestUser(t)

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
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"variables":    []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		"terminals":    []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "NUMBER", "OPERATOR", "PUNCTUATION"},
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

func TestCreateSyntaxTree_Valid(t *testing.T) {

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

func TestCreateSyntaxTree_NoGrammarCode(t *testing.T) {
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestReadGrammar_ValidNewProject(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

func TestCreateSyntaxTree_CoreError(t *testing.T) {
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestTreeToString_NoSyntaxTree(t *testing.T) {
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

	defer closeServerCore(t, server)
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
