package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestOptimiseCode_NoSourceCode(t *testing.T) {
	server = startServerCore(t)
	loginTestUser(t)

	optimising_data := map[string]interface{}{
		"users_id":         test_user_id,
		"project_name":     no_input_project_name,
		"constant_folding": true,
		"dead_code":        true,
		"loop_unrolling":   true,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/optimise", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestStoreOptimisingCode_NoInput(t *testing.T) {

	optimising_data := map[string]interface{}{
		"source_code":  "",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/source_code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestStoreOptimisingCode_Optimiser(t *testing.T) {
	optimising_data := map[string]interface{}{
		"source_code":  "package main func main()",
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/source_code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestStoreOptimisingCode_OptimiserNewProject(t *testing.T) {
	optimising_data := map[string]interface{}{
		"source_code":  "package main func main()",
		"users_id":     test_user_id,
		"project_name": new_project_name,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/source_code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestOptimiseCode_InvalidInput(t *testing.T) {
	optimising_data := map[string]interface{}{
		"users_id":       test_user_id,
		"loop_unrolling": true,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/optimise", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestOptimiseCode_NoCode(t *testing.T) {
	optimising_data := map[string]interface{}{
		"users_id":         test_user_id,
		"project_name":     project_name,
		"constant_folding": true,
		"dead_code":        true,
		"loop_unrolling":   true,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/optimise", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusInternalServerError {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestStoreOptimisingCode_UsedOptimiser(t *testing.T) {
	code := "package main\n"
	code += "func main() {\n"
	code += "return\n"
	code += "}"

	optimising_data := map[string]interface{}{
		"source_code":  code,
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/source_code", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}

func TestOptimiseCode_Success(t *testing.T) {
	defer closeServerCore(t, server)

	loginTestUser(t)
	optimising_data := map[string]interface{}{
		"users_id":         test_user_id,
		"project_name":     project_name,
		"constant_folding": true,
		"dead_code":        true,
		"loop_unrolling":   true,
	}

	req, err := json.Marshal(optimising_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/optimising/optimise", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer res.Body.Close()

	if err != nil {
		t.Errorf("Optimiser integration error: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Optimiser not working: %s", string(body_bytes))
	}
}
