package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestReadGrammar_CoreError(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

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
		"http://51.21.245.160:8080/api/parsing/grammar", "application/json",
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

	data := map[string]interface{}{
		"users_id":  test_user_id,
		"variables": []string{"STATEMENT", "DECLARATION", "EXPRESSION", "TYPE", "TERM"},
		"terminals": []string{"KEYWORD", "IDENTIFIER", "ASSIGNMENT", "INTEGER", "OPERATOR", "SEPARATOR"},
		"start":     "STATEMENT",
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
		"http://51.21.245.160:8080/api/parsing/grammar", "application/json",
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
		"users_id": no_input_user,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/tree", "application/json",
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

func TestCreateSyntaxTree_NoGrammarCode(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	test_user6 := "665c5f116aae29d323dc6a7c"
	data := map[string]interface{}{
		"users_id": test_user6,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/tree", "application/json",
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

	data := map[string]interface{}{
		"users_id": test_user_id,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/tree", "application/json",
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

	test_user7 := "675c5f116aae29d323dc6a7c"
	data := map[string]interface{}{
		"users_id": test_user7,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/tree", "application/json",
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
		if string(body_bytes) != `{"message":"Successfully created Syntax Tree","tree":{"root":{"symbol":"STATEMENT","value":"","children":[{"symbol":"DECLARATION","value":"","children":[{"symbol":"TYPE","value":"","children":[{"symbol":"KEYWORD","value":"int","children":null}]},{"symbol":"IDENTIFIER","value":"blue","children":null},{"symbol":"ASSIGNMENT","value":"=","children":null},{"symbol":"EXPRESSION","value":"","children":[{"symbol":"TERM","value":"","children":[{"symbol":"INTEGER","value":"13","children":null}]},{"symbol":"OPERATOR","value":"+","children":null},{"symbol":"TERM","value":"","children":[{"symbol":"INTEGER","value":"89","children":null}]}]}]},{"symbol":"SEPARATOR","value":";","children":null}]}}}` {
			t.Errorf("Parser not working: %s", string(body_bytes))
		}
	}
}

func TestTreeToString_NoSyntaxTree(t *testing.T) {
	server := startServerCore(t)
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id": no_input_user,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/treeString", "application/json",
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

	test_user7 := "675c5f116aae29d323dc6a7c"
	data := map[string]interface{}{
		"users_id": test_user7,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://51.21.245.160:8080/api/parsing/treeString", "application/json",
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
