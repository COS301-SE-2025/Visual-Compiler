package unit_tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// create mock requests
func createPhaseTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	contxt, eng := gin.CreateTestContext(rec)
	if eng == nil {
		t.Errorf("context not created")
	}
	return contxt, rec
}
