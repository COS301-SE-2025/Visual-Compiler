package tests

import(
	"testing"
	"github.com/COS301-SE-2025/Visual-Compiler/backend/api/routers"
	//"github.com/COS301-SE-2025/Visual-Compiler/backend/api/handlers"

	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func TestSetupRouter(t *testing.T) {
	router := routers.SetupRouter()
	if router==nil {
		t.Errorf("SetupRouter function does not initialise router")
	}
}

func TestRouterRoutes(t *testing.T) {
	router := routers.SetupRouter()
	endpoints := router.Routes()
	if len(endpoints) != 2 {
		t.Errorf("Expected routes to be registered")
	}
}

func TestRegisterEndpoint(t *testing.T) {
	//res,err := router.Register("tiaharripersad","t@gmail.com","tia1234$$")
}


