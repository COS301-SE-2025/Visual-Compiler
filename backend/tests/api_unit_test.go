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


/*// MockCollection implements the necessary methods for testing
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) *mongo.SingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*mongo.SingleResult)
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// MockDatabase for testing
type MockDatabase struct {
	mock.Mock
	collections map[string]*MockCollection
}

func (m *MockDatabase) Collection(name string) *MockCollection {
	if m.collections == nil {
		m.collections = make(map[string]*MockCollection)
	}
	if m.collections[name] == nil {
		m.collections[name] = &MockCollection{}
	}
	return m.collections[name]
}

// MockClient for testing
type MockClient struct {
	mock.Mock
	databases map[string]*MockDatabase
}

func (m *MockClient) Database(name string) *MockDatabase {
	if m.databases == nil {
		m.databases = make(map[string]*MockDatabase)
	}
	if m.databases[name] == nil {
		m.databases[name] = &MockDatabase{}
	}
	return m.databases[name]
}

// Helper function to create a test gin context
func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

// Helper function to create JSON request body
func createJSONBody(data interface{}) *bytes.Buffer {
	body, _ := json.Marshal(data)
	return bytes.NewBuffer(body)
}

func TestRegister_ValidInput(t *testing.T) {
	// Setup
	c, w := createTestContext()
	
	// Create valid request body
	reqData := Request{
		Email:    "test@example.com",
		Password: "password123",
		Username: "testuser",
	}
	
	body := createJSONBody(reqData)
	req, _ := http.NewRequest("POST", "/api/register", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// Mock database calls
	// Note: This would require modifying your code to use dependency injection
	// For now, this shows the structure of how you'd test with mocks

	// Execute
	Register(c)

	// This test would need your code to be modified to accept a database interface
	// to properly mock the MongoDB calls
	
	// For now, we can test the input validation part
	assert.Equal(t, http.StatusCreated, w.Code) // This would fail without proper mocking
}*/