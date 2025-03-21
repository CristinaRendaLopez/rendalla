package utils

import (
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func CreateTestContext(method, url string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, url, body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	return c, w
}
