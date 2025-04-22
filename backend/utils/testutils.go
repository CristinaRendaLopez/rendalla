package utils

import (
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// CreateTestContext creates a new test context and response recorder for unit testing handlers.
// It sets up a JSON request with the given HTTP method, URL, and optional body payload.
//
// Returns:
//   - a *gin.Context initialized with the request
//   - a *httptest.ResponseRecorder to inspect the response
func CreateTestContext(method, url string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}
