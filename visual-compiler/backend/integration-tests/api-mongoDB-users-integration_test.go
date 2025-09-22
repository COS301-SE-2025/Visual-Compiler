package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestRegister_ExistingEmail(t *testing.T) {
	server = startServerCore(t)

	user_data := map[string]string{
		"username": "test u",
		"email":    "halfstack.testuser@gmail.com",
		"password": "password",
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
		t.Errorf("Error: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusConflict {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)

		if err != nil {
			t.Errorf("registration failed: %v", err)
		}
		if body_array["error"] != "Email already exists" {
			t.Errorf("Incorrect error")
		}
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("register works for existing user: %s", string(body_bytes))
	}

}

func TestRegister_NewUser(t *testing.T) {

	username := fmt.Sprintf("user-%s", time.Now().Format("20060102-150405"))
	email := fmt.Sprintf("%s@gmail.com", time.Now().Format("20060102-150405"))
	user_data := map[string]string{
		"username": username,
		"email":    email,
		"password": "jazzy1234$$",
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

func TestRegister_ExistingUsername(t *testing.T) {

	user_data := map[string]string{
		"username": "test user",
		"email":    "halfstack.testuse@gmail.com",
		"password": "password",
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
	if res.StatusCode != http.StatusConflict {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("registration failed: %s", string(body_bytes))
	} else {
		body_bytes, _ := io.ReadAll(res.Body)
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)

		if err != nil {
			t.Errorf("registration failed: %v", err)
		}

		if body_array["error"] != "Username is already taken" {
			t.Errorf("Incorrect error")
		}
	}

}

func TestLogin_Success(t *testing.T) {

	user_data := map[string]string{
		"login":    "test user",
		"email":    "halfstack.testuser@gmail.com",
		"password": "testUser13",
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
	}
}

func TestEditUser_NoAdmin(t *testing.T) {

	user_data := map[string]string{
		"login":    "test user",
		"email":    "halfstack.testuser@gmail.com",
		"password": "password",
	}
	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.NewRequest("PATCH",
		"http://localhost:8080/api/users/update",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("User editing failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusBadRequest {
		body_bytes, _ := io.ReadAll(response.Body)
		t.Errorf("User edit failed: %s", string(body_bytes))
	}
}

func TestEditUser_NotAdmin(t *testing.T) {

	user_data := map[string]string{
		"users_id": test_user_id,
		"admin_id": test_user_id,
		"username": "jasmine1",
	}
	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("converting data to json failed")
	}

	res, err := http.NewRequest("PATCH",
		"http://localhost:8080/api/users/update",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("User editing failed")
	}
	res.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	response, err := client.Do(res)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusForbidden {
		body_bytes, _ := io.ReadAll(response.Body)
		t.Errorf("User edit failed: %s", string(body_bytes))
	}
}

func TestLogin_IncorrectPassword(t *testing.T) {

	user_data := map[string]string{
		"login":    "test user",
		"email":    "halfstack.testuser@gmail.com",
		"password": "halfStack",
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
	if res.StatusCode != http.StatusUnauthorized {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Login failed: %s", string(body_bytes))
	} else {

		body_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if body_array["error"] != "Password is incorrect" {
			t.Errorf("Incorrect error")
		}
	}
}

func TestLogin_InvalidUser(t *testing.T) {

	user_data := map[string]string{
		"login":    "",
		"email":    "halfstack.testuser@gmail.com",
		"password": "password",
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
		t.Logf("Login working: %s", string(body_bytes))
	}

	if res.StatusCode == http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Logf("Login error: %s", string(body_bytes))
	}
}

func TestDeleteUser_Invalid(t *testing.T) {

	user_data := map[string]string{
		"users_id": "1",
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
		t.Logf("User deletion working: %s", string(body_bytes))
	}

	if response.StatusCode == http.StatusOK {
		body_bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Deletion error")
		}

		t.Errorf("Deletion invalid: %s", body_array["message"])
	}
}

func TestGetAllUsers(t *testing.T) {

	res, err := http.Get(
		"http://localhost:8080/api/users/getUsers",
	)
	if err != nil {
		t.Errorf("Get all users failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body_bytes, _ := io.ReadAll(res.Body)
		t.Errorf("Get all users failed: %s", string(body_bytes))
	} else {

		body_bytes, err := io.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		_ = json.Unmarshal(body_bytes, &body_array)

		t.Logf("Get all users working: %s", body_array["message"])
	}
}

func TestConnectToMongo_Valid(t *testing.T) {
	defer closeServerCore(t, server)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/test-mongo", handlers.ConnectToMongo)

	req, err := http.NewRequest("GET", "/test-mongo", nil)
	if err != nil {
		t.Errorf("Connection failed")
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		body_bytes, _ := io.ReadAll(response.Body)
		t.Errorf("Connection not working: %s", string(body_bytes))
	}

	if response.Code == http.StatusOK {
		body_bytes, err := io.ReadAll(response.Body)
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		var body_array map[string]string
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Connection error")
		}
	}
}
