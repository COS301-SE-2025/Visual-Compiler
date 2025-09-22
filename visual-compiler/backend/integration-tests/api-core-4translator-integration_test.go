package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestReadRules_Failed(t *testing.T) {
	server = startServerCore(t)
	loginTestUser(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"translation_rules": []map[string]interface{}{
			{
				"sequence":    "KEYWORD IDENTIFIER OPERATOR NUMBER PUNCTUATION",
				"translation": []string{},
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/readRules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)
	}
}

func TestReadRules_RegexFailed(t *testing.T) {
	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"translation_rules": []map[string]interface{}{
			{
				"sequence": []string{
					"KEYWORD OPERATOR NUMBER PUNCTUATION",
				},
				"translation": []string{
					"MOVE REGISTER {IDENTIFIER}",
					"ADDI REGISTER {NUMBER}",
					"MOVE {IDENTIFIER} REGISTER",
				},
			},
			{
				"sequence": []string{
					"KEYWORD  ASSIGNMENT INTEGER OPERATOR INTEGER SEPARATOR",
				},
				"translation": []string{
					"MOVE REGISTER {IDENTIFIER}",
					"ADDI REGISTER {INTEGER}",
					"ADDI REGISTER {INTEGER}",
					"MOVE {IDENTIFIER} REGISTER",
				},
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/readRules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}
}

func TestReadRules_Success(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"translation_rules": []map[string]interface{}{
			{
				"sequence": []string{
					"KEYWORD IDENTIFIER OPERATOR NUMBER PUNCTUATION",
				},
				"translation": []string{
					"MOVE REGISTER {IDENTIFIER}",
					"ADDI REGISTER {NUMBER}",
					"MOVE {IDENTIFIER} REGISTER",
				},
			},
			{
				"sequence": []string{
					"KEYWORD IDENTIFIER ASSIGNMENT INTEGER OPERATOR INTEGER SEPARATOR",
				},
				"translation": []string{
					"MOVE REGISTER {IDENTIFIER}",
					"ADDI REGISTER {INTEGER}",
					"ADDI REGISTER {INTEGER}",
					"MOVE {IDENTIFIER} REGISTER",
				},
			},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/readRules", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}
}

func TestTranslate_Success(t *testing.T) {
	data := map[string]string{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/translate", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}
}

func TestTranslate_NoRules(t *testing.T) {
	data := map[string]string{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/translate", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)
		if body_array["error"] != "Tree not found. Please go back to parsing" {
			t.Errorf("Analyser not working: %s", string(body_bytes))
		}
	}
}

func TestTranslate_NoTree(t *testing.T) {
	defer closeServerCore(t, server)

	data := map[string]string{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/translating/translate", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Translator not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)
		if body_array["error"] != "Tree not found. Please go back to parsing" {
			t.Errorf("Analyser not working: %s", string(body_bytes))
		}
	}

	deleteNewProject(t)
	deleteProject(t)
}
