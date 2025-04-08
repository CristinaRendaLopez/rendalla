package integration_tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Document struct {
	ID         string   `json:"id"`
	SongID     string   `json:"song_id"`
	Type       string   `json:"type"`
	Instrument []string `json:"instrument"`
	PDFURL     string   `json:"pdf_url"`
	AudioURL   string   `json:"audio_url,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type DocumentResponse struct {
	Data Document `json:"data"`
}

type DocumentsResponse struct {
	Data []Document `json:"data"`
}

func TestGetDocumentsBySongID_ShouldReturnSeededDocuments(t *testing.T) {
	w := MakeRequest("GET", "/songs/queen-001/documents", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response DocumentsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Data, 2)
}

func TestGetDocumentsBySongID_ShouldReturnEmptyList(t *testing.T) {
	w := MakeRequest("GET", "/songs/non-existent-id/documents", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response DocumentsResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Data)
}

func TestGetDocumentByID_ShouldReturnSeededDocument(t *testing.T) {
	w := MakeRequest("GET", "/songs/queen-001/documents/doc-br-piano", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response DocumentResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "doc-br-piano", response.Data.ID)
}

func TestGetDocumentByID_ShouldReturn404(t *testing.T) {
	w := MakeRequest("GET", "/songs/queen-001/documents/non-existent-id", nil, "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateDocument_ShouldSucceedWithJWT(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"type": "partitura",
		"instrument": ["violín"],
		"pdf_url": "https://s3.test/test_violin.pdf"
	}`

	w := MakeRequest("POST", "/songs/queen-001/documents", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusCreated, w.Code)

	var res struct {
		Message    string `json:"message"`
		DocumentID string `json:"document_id"`
	}
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)
	assert.Equal(t, "Document created successfully", res.Message)
	assert.NotEmpty(t, res.DocumentID)
}

func TestCreateDocument_ShouldReturn401WithoutToken(t *testing.T) {
	payload := `{
		"type": "partitura",
		"instrument": ["guitarra"],
		"pdf_url": "https://s3.test/test_guitar.pdf"
	}`

	w := MakeRequest("POST", "/songs/queen-001/documents", strings.NewReader(payload), "")
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateDocument_ShouldReturn400ForInvalidJSON(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{`

	w := MakeRequest("POST", "/songs/queen-001/documents", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateDocument_ShouldReturn400ForInvalidFields(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"type": "",
		"instrument": [],
		"pdf_url": ""
	}`

	w := MakeRequest("POST", "/songs/queen-001/documents", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDocument_ShouldSucceed(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{"type": "tablatura"}`

	w := MakeRequest("PUT", "/songs/queen-001/documents/doc-br-piano", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusOK, w.Code)
}

// func TestUpdateDocument_ShouldReturn404IfNotExists(t *testing.T) {
// 	token, err := GenerateTestJWT("admin")
// 	assert.NoError(t, err)

// 	payload := `{"type": "partitura"}`

// 	w := MakeRequest("PUT", "/documents/non-existent-id", strings.NewReader(payload), token)
// 	w.Header().Set("Content-Type", "application/json")

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }

func TestUpdateDocument_ShouldReturn400ForInvalidJSON(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{"type":`

	w := MakeRequest("PUT", "/songs/queen-001/documents/doc-br-voice", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteDocument_ShouldSucceed(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"type": "partitura",
		"instrument": ["flauta"],
		"pdf_url": "https://s3.test/test_flute.pdf"
	}`
	createRes := MakeRequest("POST", "/songs/queen-001/documents", strings.NewReader(payload), token)
	createRes.Header().Set("Content-Type", "application/json")

	var createBody struct {
		DocumentID string `json:"document_id"`
	}
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, createBody.DocumentID)

	deleteRes := MakeRequest("DELETE", "/songs/queen-001/documents/"+createBody.DocumentID, nil, token)
	assert.Equal(t, http.StatusOK, deleteRes.Code)
}

// func TestDeleteDocument_ShouldReturn404IfNotExists(t *testing.T) {
// 	token, err := GenerateTestJWT("admin")
// 	assert.NoError(t, err)

// 	w := MakeRequest("DELETE", "/documents/non-existent-id", nil, token)
// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }

func TestDeleteDocument_ShouldReturn401WithoutToken(t *testing.T) {
	w := MakeRequest("DELETE", "/songs/queen-001/documents/doc-br-voice", nil, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
