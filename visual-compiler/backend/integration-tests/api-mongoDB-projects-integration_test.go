package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func loginUser(t *testing.T) {
	//server := startServer(t)
	//defer closeServer(server)

	user_data := map[string]string{
		"login":    "integration_tester",
		"password": "Password1234$",
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

		test_user_id = body_array["id"]
		project_name = "project1"
	}
}

func TestSaveProjectName_Success(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	user_data := map[string]string{
		"users_id":     test_user_id,
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

func TestSaveProjectName_ExistingProject(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	user_data := map[string]string{
		"users_id":     test_user_id,
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

	if res.StatusCode == http.StatusConflict {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)

		if err != nil {
			t.Errorf("Save project name failed: %v", err)
		}
		if body_array["error"] != "Project name already exists" {
			t.Errorf("Incorrect error")
		}
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Save project works for existing project name: %s", string(body_bytes))
	}

}

func TestSaveProjectPipeline_Success(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	user_data := map[string]any{
		"users_id":     test_user_id,
		"project_name": project_name,
		"pipeline": map[string]any{
			"node-1": "source_code",
			"node-2": "lexer",
			"node-3": "parser",
		},
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/savePipeline", "application/json",
		bytes.NewBuffer(req),
	)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("saving project pipeline failed: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)

	}

}

func TestDeleteProject_Success(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	user_data := map[string]string{
		"users_id":     test_user_id,
		"project_name": project_name,
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.NewRequest("DELETE",
		"http://localhost:8080/api/users/deleteProject",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Project deletion failed")
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
		t.Errorf("Project deletion failed: %s", string(body_bytes))
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

func TestSaveProjectPipeline_UserNotFound(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	user_data := map[string]any{
		"users_id":     test_user_id,
		"project_name": project_name,
		"pipeline": map[string]any{
			"node-1": "source_code",
			"node-2": "lexer",
			"node-3": "parser",
		},
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.Post(
		"http://localhost:8080/api/users/savePipeline", "application/json",
		bytes.NewBuffer(req),
	)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Project deletion failed: %s", string(body_bytes))
	}

}

func TestDeleteProject_UserNotFound(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	user_data := map[string]string{
		"users_id":     "123e45df6f7b89f50014505a",
		"project_name": project_name,
	}

	req, err := json.Marshal(user_data)

	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.NewRequest("DELETE",
		"http://localhost:8080/api/users/deleteProject",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("Project deletion failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(response.Body)
		t.Errorf("Project deletion failed: %s", string(body_bytes))
	}

}

func TestGetAllProjects_Success(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)
	user_url := fmt.Sprintf("http://localhost:8080/api/users/getProjects?users_id=%s", test_user_id)

	res, err := http.Get(user_url)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Get all projects failed: %s", string(body_bytes))
	}

}

func TestGetAllProjects_UserNotFound(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	res, err := http.Get("http://localhost:8080/api/users/getProjects?users_id=123e45df6f7b89f50014505a")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Get all projects failed: %s", string(body_bytes))
	}

}

func TestGetProject_Success(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)
	project_url := "http://localhost:8080/api/users/getProject"
	url_param := url.Values{}
	url_param.Add("users_id", test_user_id)
	fmt.Printf(test_user_id)
	url_param.Add("project_name", "test_project")
	user_url := fmt.Sprintf("%s?%s", project_url, url_param.Encode())

	res, err := http.Get(user_url)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Get all projects failed: %s", string(body_bytes))
	}

}

func TestGetProject_UserIdInvalid(t *testing.T) {
	server := startServer(t)
	defer closeServer(server)

	loginUser(t)

	res, err := http.Get("http://localhost:8080/api/users/getProject?users_id=123e45df6f7b4505a&project_name=12")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Get all projects failed: %s", string(body_bytes))
	}

}
