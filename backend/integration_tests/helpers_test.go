package integration_tests

import (
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func MakeRequest(router *gin.Engine, method, path string, body io.Reader, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
