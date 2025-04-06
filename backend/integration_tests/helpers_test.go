package integration_tests

import (
	"io"
	"net/http/httptest"
)

func MakeRequest(method, path string, body io.Reader, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)
	return w
}
