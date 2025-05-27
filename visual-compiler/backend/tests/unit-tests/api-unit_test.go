package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
)

func TestSetupRouter(t *testing.T) {
	router := routers.SetupUserRouter()
	if router == nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestRouterRoutes(t *testing.T) {
	router := routers.SetupUserRouter()
	endpoints := router.Routes()
	if len(endpoints) != 2 {
		t.Errorf("Expected routes to be registered")
	}
}

// create mock requests
func createTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}

func TestRegisterInvalidEmail(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Password: "passsdw32323@@@",
		Username: "invalid Email",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestRegisterInvalidPassword(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Email:    "user@gmail.com",
		Password: "pas",
		Username: "invalid Password",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestRegisterInvalidUsername(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.Request{
		Email:    "user@gmail.com",
		Password: "passsdw32323@@@",
		Username: "in",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("POST", "/api/register", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Register(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}

}

func TestLoginNoName(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.LoginReq{
		Password: "passsdw32323@@@",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("GET", "/api/login", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Login(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}

func TestLoginNoPassword(t *testing.T) {
	contxt, rec := createTestContext(t)

	user_data := handlers.LoginReq{
		Login: "passsdw32323@@@",
	}

	req, err := json.Marshal(user_data)
	if err != nil {
		t.Errorf("Could not data to json")
	}

	res, err := http.NewRequest("GET", "/api/login", bytes.NewBuffer(req))
	if err != nil {
		t.Errorf("Request could not be created")
	}
	res.Header.Set("Content-Type", "application/json")
	contxt.Request = res

	handlers.Login(contxt)

	if rec.Code == http.StatusBadRequest {
		var mess map[string]string
		t.Logf(mess["error"])
	}
}
