package tests

import (
	//"bytes"
	//"encoding/json"
	//"io"
	"testing"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"

	//"bytes"
	//"encoding/json"
	"net/http"
	//"io"
	"context"
	"errors"
	"time"
)

func startServer(t *testing.T) *http.Server {
	router := routers.SetupUserRouter()
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

	return server
}

func closeServer(t *testing.T, server *http.Server) {
	cont, cancel_cont := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel_cont()
	server.Shutdown(cont)
}

func TestRegisterExistingUser(t *testing.T) {
	/*server := startServer(t)

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
		"http://localhost:8080/api/register", "application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		t.Errorf("could not send request to api")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Logf("register working: %s", string(bodyBytes))
	}

	if res.StatusCode == http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Errorf("register works for existing user: %s", string(bodyBytes))
	}

	closeServer(t, server)*/
}

func TestRegisterNewUser(t *testing.T) {
	/*server := startServer(t)

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
		"http://localhost:8080/api/register", "application/json",
		bytes.NewBuffer(req),
	)
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(res.Body)
		t.Errorf("registration failed: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	var respMap map[string]string
	err = json.Unmarshal(bodyBytes, &respMap)

	t.Logf("Register working: %s", respMap["message"])

	closeServer(t, server)*/
}

func LoginExistingUser(t *testing.T) {
	//login integration testing to be implemented
}

func LoginNewUser(t *testing.T) {
	//login integration testing to be implemented
}
