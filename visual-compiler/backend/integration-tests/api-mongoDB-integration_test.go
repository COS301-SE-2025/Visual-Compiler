package tests

import (
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

var user_id string

func startServer(t *testing.T) *http.Server {
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
	user_router := routers.SetupUserRouter()
	//go run main.go
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

func closeServer(t *testing.T, server *http.Server) {
	cont, cancel_cont := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel_cont()
	server.Shutdown(cont)
}

func TestRegister_ExistingEmail(t *testing.T) {
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"username": "tiaharripersad",
		"email":    "t@gmail.com",
		"password": "tia1234$$",
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
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"username": "jasmine1",
		"email":    "j@gmail.com",
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
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"username": "jasmine1",
		"email":    "jah@gmail.com",
		"password": "jazzyo234$$",
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

func TestLogin_ExistingUser(t *testing.T) {
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"login":    "jasmine1",
		"password": "jazzy1234$$",
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
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		t.Logf("Login working: %s", body_array["message"])
		user_id = body_array["id"]
	}
}

func TestLogin_IncorrectPassword(t *testing.T) {
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"login":    "jasmine1",
		"password": "jazzy16899",
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
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"login":    "rando",
		"password": "jazzy1234$$",
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

func TestDeleteUser_Existing(t *testing.T) {
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"id": user_id,
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
		err = json.Unmarshal(body_bytes, &body_array)
		if err != nil {
			t.Errorf("Error: %v", err)
		}

		t.Logf("Deletion working: %s", body_array["message"])
	}
}

func TestDeleteUser_Invalid(t *testing.T) {
	server := startServer(t)
	defer closeServer(t, server)

	user_data := map[string]string{
		"id": user_id,
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
	server := startServer(t)
	defer closeServer(t, server)

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
