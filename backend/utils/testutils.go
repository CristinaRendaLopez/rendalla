package utils

import (
	"bytes"
	"io"
	"mime/multipart"
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

// TestFile represents a simulated file in multipart/form-data tests
type TestFile struct {
	Filename string
	Content  []byte
}

// CreateMultipartRequest builds a multipart/form-data HTTP request body for testing handlers.
// It populates form fields and file inputs with the provided content.
//
// Parameters:
//   - fields: key-value pairs representing regular form fields (e.g., "type", "instrument[]")
//   - files: a map where each key is the form field name for a file (e.g., "pdf"), and the value is a TestFile struct
//
// Returns:
//   - a *bytes.Buffer containing the full multipart body
//   - a string with the correct Content-Type header
//   - an error if any part of the request construction fails
func CreateMultipartRequest(fields map[string]string, files map[string]TestFile) (*bytes.Buffer, string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, val := range fields {
		if err := w.WriteField(key, val); err != nil {
			return nil, "", err
		}
	}

	for fieldname, tf := range files {
		fw, err := w.CreateFormFile(fieldname, tf.Filename)
		if err != nil {
			return nil, "", err
		}
		if _, err := fw.Write(tf.Content); err != nil {
			return nil, "", err
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return &b, w.FormDataContentType(), nil
}
