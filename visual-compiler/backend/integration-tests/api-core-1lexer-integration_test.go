package tests

import (
	"fmt"
	"testing"

	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

func TestStoreSourceCode_InvalidInput(t *testing.T) {
	server = startServerCore(t)

	loginTestUser(t)
	addNewProject(t)

	re_data := map[string]interface{}{
		"source_code":  "",
		"users_id":     test_user_id,
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
	if code_res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(code_res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}
}

func TestLexing_CoreError(t *testing.T) {

	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
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

func TestTokensFromDFA_NoSourceCode(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestStoreSourceCode_Valid(t *testing.T) {
	lexer_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req_data, err := json.Marshal(lexer_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/code",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testToken")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Code is ready for further processing"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_InvalidInput(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
		"pairs":        "",
	}

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] != "Input is invalid" {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_CoreError(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
		"pairs": []map[string]string{
			{
				"Type":  "KEYWORD",
				"Regex": "[a-z",
			},
		},
	}

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestLexing_NoSourceCode(t *testing.T) {
	reg_expr_data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestStoreSourceCode_ValidNewProject(t *testing.T) {

	re_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
		"project_name": new_project_name,
	}

	req_data, err := json.Marshal(re_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/code",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
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

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
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

func TestCreateRulesFromCode_ValidNewProject(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
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

func TestLexing_ValidNewProject(t *testing.T) {

	reg_expr_data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestReadDFAFromUser_InvalidInput(t *testing.T) {

	data := map[string]interface{}{

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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestReadDFAFromUser_Valid(t *testing.T) {

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

func TestConvertDFAToRG_NoDFA(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}
}

func TestReadDFAFromUser_ValidNewProject(t *testing.T) {

	data := map[string]interface{}{
		"project_name": new_project_name,
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

func TestReadDFAFromUser_NewProject(t *testing.T) {

	data := map[string]interface{}{
		"project_name": new_project_name,
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

func TestTokensFromDFA_InvalidInput(t *testing.T) {

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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestTokensFromDFA_Valid(t *testing.T) {

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

func TestTokensFromDFA_CoreError(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestConvertDFAToRG_InvalidInput(t *testing.T) {
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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestConvertDFAToRG_Valid(t *testing.T) {
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

func TestConvertDFAToRG_NoDFANewProject(t *testing.T) {
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

}

func TestConvertDFAToRG_CoreError(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestConvertRGToNFA_InvalidInput(t *testing.T) {
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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestConvertRGToNFA_Valid(t *testing.T) {
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

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] != "Regex rules not found. Please create one" {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToNFA_Failed(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]interface{}
		err := json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
		if body_array["error"] == "Conversion from Regex to NFA failed" {
			t.Logf("ConvertRGToNFA: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestConvertRGToNFA_CoreError(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestConvertRGToDFA_InvalidInput(t *testing.T) {

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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestConvertRGToDFA_Valid(t *testing.T) {

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

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

}

func TestConvertNFAToDFA_InvalidInput(t *testing.T) {
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

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestConvertNFAToDFA_Valid(t *testing.T) {
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
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
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

}

func TestConvertNFAToDFA_CoreError(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

func TestGetCode_InvalidInput(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getCode?project_name=%s",
		url.QueryEscape(""),
	)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestGetCode_SourceCodeNotFound(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getCode?project_name=%s",
		url.QueryEscape(no_input_project_name),
	)

	res, err := http.Get(url)
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

}

func TestGetCode_Valid(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getCode?project_name=%s",
		url.QueryEscape(project_name),
	)

	res, err := http.Get(url)
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

}

func TestGetTokens_InvalidInput(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getTokens?project_name=%s",
		url.QueryEscape(""),
	)

	res, err := http.Get(url)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	}

}

func TestGetTokens_SourceCodeNotFound(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getTokens?project_name=%s",
		url.QueryEscape(no_input_project_name),
	)

	res, err := http.Get(url)
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

}

func TestGetTokens_Valid(t *testing.T) {
	url := fmt.Sprintf(
		"http://localhost:8080/api/lexing/getTokens?project_name=%s",
		url.QueryEscape(project_name),
	)

	res, err := http.Get(url)
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

}

func TestStoreSourceCode_Val(t *testing.T) {
	lexer_data := map[string]interface{}{
		"source_code":  "int x = 2 ;",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req_data, err := json.Marshal(lexer_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/code",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testToken")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Lexing integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Lexer not working: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		if string(body_bytes) == `{"message":"Code is ready for further processing"}` {
			t.Logf("ReadDFAFromUser: success")
		} else {
			t.Errorf("Error: %v", string(body_bytes))
		}
	}

}

func TestCreateRulesFromCode_Val(t *testing.T) {

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

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
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
func TestLexing_Val(t *testing.T) {

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

func TestCreateRulesFromCode_VaNewProjectl(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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

	req_data, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	req, err := http.NewRequest(
		"POST",
		"http://localhost:8080/api/lexing/rules",
		bytes.NewBuffer(req_data),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+auth0_id)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
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

func TestLexing_ValNewProject(t *testing.T) {
	defer closeServerCore(t, server)

	reg_expr_data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": new_project_name,
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
