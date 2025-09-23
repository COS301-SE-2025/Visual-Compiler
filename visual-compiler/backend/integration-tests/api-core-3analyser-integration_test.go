package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestAnalyse_ArtefactCreationFailed(t *testing.T) {
	server = startServerCore(t)
	loginTestUser(t)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"scope_rules": []map[string]string{
			{"Start": "{"},
			{"End": "}"},
		},
		"grammar_rules": map[string]string{
			"TypeRule":       "TYPE",
			"VariableRule":   "IDENTIFIER",
			"FunctionRule":   "FUNC",
			"ParameterRule":  "PARAM",
			"AssignmentRule": "ASSIGNMENT",
			"OperatorRule":   "OPERATOR",
			"TermRule":       "TERM",
		},
		"type_rules": []map[string]interface{}{
			{"ResultData": "int"},
			{"Assignment": "="},
			{"LHSData": "INTEGER"},
			{"Operator": []string{"+"}},
			{"RHSData": "INTEGER"},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/analysing/analyse", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Analyser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)
		if body_array["error"] != "Artefacts creation failed" {
			t.Errorf("Analyser not working: %s", string(body_bytes))
		}
	}
}

func TestAnalyse_Success(t *testing.T) {

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": project_name,
		"scope_rules": []map[string]string{
			{
				"Start": "{",
				"End":   "}",
			},
		},
		"grammar_rules": map[string]string{
			"TypeRule":       "TYPE",
			"VariableRule":   "IDENTIFIER",
			"FunctionRule":   "FUNC",
			"ParameterRule":  "PARAM",
			"AssignmentRule": "ASSIGNMENT",
			"OperatorRule":   "OPERATOR",
			"TermRule":       "TERM",
		},
		"type_rules": []map[string]interface{}{
			{"ResultData": "int"},
			{"Assignment": "="},
			{"LHSData": "INTEGER"},
			{"Operator": []string{"+"}},
			{"RHSData": "INTEGER"},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/analysing/analyse", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Analyser not working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)
		if body_array["message"] != "Symbol Table Artefact successfully inserted." {
			t.Errorf("Analyser not working: %s", string(body_bytes))
		}
	}
}

func TestAnalyse_NoTree(t *testing.T) {
	defer closeServerCore(t, server)

	data := map[string]interface{}{
		"users_id":     test_user_id,
		"project_name": no_input_project_name,
		"scope_rules": []map[string]string{
			{
				"Start": "{",
				"End":   "}",
			},
		},
		"grammar_rules": map[string]string{
			"TypeRule":       "TYPE",
			"VariableRule":   "IDENTIFIER",
			"FunctionRule":   "FUNC",
			"ParameterRule":  "PARAM",
			"AssignmentRule": "ASSIGNMENT",
			"OperatorRule":   "OPERATOR",
			"TermRule":       "TERM",
		},
		"type_rules": []map[string]interface{}{
			{"ResultData": "int"},
			{"Assignment": "="},
			{"LHSData": "INTEGER"},
			{"Operator": []string{"+"}},
			{"RHSData": "INTEGER"},
		},
	}

	req, err := json.Marshal(data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/analysing/analyse", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error not expected")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Analyser not working: %s", string(body_bytes))
	}
}
