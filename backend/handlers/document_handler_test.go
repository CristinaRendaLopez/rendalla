package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupDocumentHandlerTest() (*handlers.DocumentHandler, *mocks.MockDocumentService) {
	mockService := new(mocks.MockDocumentService)
	handler := handlers.NewDocumentHandler(mockService)
	return handler, mockService
}

func TestCreateDocumentHandler(t *testing.T) {
	tests := []struct {
		name          string
		songID        string
		setupParam    bool
		body          string
		expectedCode  int
		mockReturnID  string
		mockReturnErr error
		expectedDocID string
	}{
		{
			name:          "successfully creates document",
			songID:        "1",
			setupParam:    true,
			body:          DocumentValidJSON,
			expectedCode:  http.StatusCreated,
			mockReturnID:  "doc-123",
			mockReturnErr: nil,
			expectedDocID: "doc-123",
		},
		{
			name:         "invalid JSON payload",
			songID:       "1",
			setupParam:   true,
			body:         DocumentInvalidJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing song_id parameter",
			setupParam:   false,
			body:         DocumentValidJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation fails",
			songID:       "1",
			setupParam:   true,
			body:         DocumentInvalidDataJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "internal service error",
			songID:        "1",
			setupParam:    true,
			body:          DocumentValidJSON,
			expectedCode:  http.StatusInternalServerError,
			mockReturnErr: errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupDocumentHandlerTest()

			if tt.setupParam && tt.expectedCode == http.StatusCreated || tt.expectedCode == http.StatusInternalServerError {
				var req dto.CreateDocumentRequest
				_ = json.Unmarshal([]byte(tt.body), &req)
				req.SongID = tt.songID

				mockService.
					On("CreateDocument", req).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}

			path := "/songs"
			if tt.setupParam {
				path += "/" + tt.songID + "/documents"
			}

			c, w := utils.CreateTestContext(http.MethodPost, path, strings.NewReader(tt.body))
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}

			handler.CreateDocumentHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusCreated {
				var response struct {
					Message    string `json:"message"`
					DocumentID string `json:"document_id"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDocID, response.DocumentID)
			}

			if tt.setupParam && tt.mockReturnErr != errors.ErrValidationFailed {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestGetAllDocumentsBySongIDHandler(t *testing.T) {
	tests := []struct {
		name           string
		songID         string
		setupParam     bool
		mockDocuments  []dto.DocumentResponseItem
		mockError      error
		expectedCode   int
		expectedResult []dto.DocumentResponseItem
	}{
		{
			name:           "successfully returns documents",
			songID:         "1",
			setupParam:     true,
			mockDocuments:  []dto.DocumentResponseItem{DocumentResponseScore, DocumentResponseTablature},
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: []dto.DocumentResponseItem{DocumentResponseScore, DocumentResponseTablature}},
		{
			name:           "returns empty document list",
			songID:         "1",
			setupParam:     true,
			mockDocuments:  []dto.DocumentResponseItem{},
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: []dto.DocumentResponseItem{},
		},
		{
			name:         "missing song_id param",
			setupParam:   false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "internal service error",
			songID:       "1",
			setupParam:   true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupDocumentHandlerTest()

			if tt.setupParam && tt.mockError != errors.ErrValidationFailed {
				mockService.
					On("GetDocumentsBySongID", tt.songID).
					Return(tt.mockDocuments, tt.mockError)
			}

			path := "/songs"
			if tt.setupParam {
				path += "/" + tt.songID + "/documents"
			}

			c, w := utils.CreateTestContext(http.MethodGet, path, nil)
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}

			handler.GetAllDocumentsBySongIDHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response struct {
					Data []dto.DocumentResponseItem `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, response.Data)
			}

			if tt.setupParam && tt.mockError != errors.ErrValidationFailed {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestGetDocumentByIDHandler(t *testing.T) {
	tests := []struct {
		name           string
		songID         string
		docID          string
		setupParams    bool
		mockResult     dto.DocumentResponseItem
		mockError      error
		expectedCode   int
		expectedResult dto.DocumentResponseItem
	}{
		{
			name:           "successfully retrieves document",
			songID:         "1",
			docID:          "doc-123",
			setupParams:    true,
			mockResult:     DocumentResponseScore,
			mockError:      nil,
			expectedCode:   http.StatusOK,
			expectedResult: DocumentResponseScore,
		},
		{
			name:         "missing song_id param",
			docID:        "doc-123",
			setupParams:  true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing doc_id param",
			songID:       "1",
			setupParams:  true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "document not found",
			songID:       "1",
			docID:        "doc-999",
			setupParams:  true,
			mockError:    errors.ErrResourceNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "internal service error",
			songID:       "1",
			docID:        "doc-123",
			setupParams:  true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupDocumentHandlerTest()

			if tt.setupParams && tt.expectedCode != http.StatusBadRequest {
				mockService.
					On("GetDocumentByID", tt.songID, tt.docID).
					Return(tt.mockResult, tt.mockError)
			}

			path := "/songs"
			if tt.setupParams {
				path += "/" + tt.songID + "/documents/" + tt.docID
			}

			c, w := utils.CreateTestContext(http.MethodGet, path, nil)
			if tt.setupParams {
				if tt.songID != "" {
					c.Params = append(c.Params, gin.Param{Key: "song_id", Value: tt.songID})
				}
				if tt.docID != "" {
					c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: tt.docID})
				}
			}

			handler.GetDocumentByIDHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response struct {
					Data dto.DocumentResponseItem `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, response.Data)
			}

			if tt.setupParams && tt.expectedCode != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUpdateDocumentHandler(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		docID        string
		setupParams  bool
		body         string
		expectedCode int
		mockError    error
	}{
		{
			name:         "successfully updates document",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			body:         DocumentUpdateValidJSON,
			expectedCode: http.StatusOK,
			mockError:    nil,
		},
		{
			name:         "missing song_id param",
			docID:        "doc-1",
			setupParams:  true,
			body:         DocumentUpdateOnlyTypeJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing doc_id param",
			songID:       "1",
			setupParams:  true,
			body:         DocumentUpdateOnlyTypeJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid JSON payload",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			body:         DocumentUpdateInvalidJSON,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation fails (empty instrument)",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			body:         DocumentUpdateEmptyInstrumentJSON,
			expectedCode: http.StatusBadRequest,
			mockError:    errors.ErrValidationFailed,
		},
		{
			name:         "document not found",
			songID:       "1",
			docID:        "doc-404",
			setupParams:  true,
			body:         DocumentUpdateOnlyTypeJSON,
			expectedCode: http.StatusNotFound,
			mockError:    errors.ErrResourceNotFound,
		},
		{
			name:         "internal service error",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			body:         DocumentUpdateOnlyTypeJSON,
			expectedCode: http.StatusInternalServerError,
			mockError:    errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupDocumentHandlerTest()

			if tt.setupParams && tt.mockError != errors.ErrValidationFailed {
				var update dto.UpdateDocumentRequest
				_ = json.Unmarshal([]byte(tt.body), &update)

				mockService.
					On("UpdateDocument", tt.songID, tt.docID, update).
					Return(tt.mockError)
			}

			path := "/songs"
			if tt.setupParams {
				path += "/" + tt.songID + "/documents/" + tt.docID
			}

			c, w := utils.CreateTestContext(http.MethodPut, path, strings.NewReader(tt.body))
			if tt.setupParams {
				if tt.songID != "" {
					c.Params = append(c.Params, gin.Param{Key: "song_id", Value: tt.songID})
				}
				if tt.docID != "" {
					c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: tt.docID})
				}
			}

			handler.UpdateDocumentHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)
		})
	}
}

func TestDeleteDocumentHandler(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		docID        string
		setupParams  bool
		mockError    error
		expectedCode int
	}{
		{
			name:         "successfully deletes document",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing song_id param",
			docID:        "doc-1",
			setupParams:  true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing doc_id param",
			songID:       "1",
			setupParams:  true,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "document not found",
			songID:       "1",
			docID:        "doc-999",
			setupParams:  true,
			mockError:    errors.ErrResourceNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "internal service error",
			songID:       "1",
			docID:        "doc-1",
			setupParams:  true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupDocumentHandlerTest()

			if tt.setupParams && tt.expectedCode != http.StatusBadRequest {
				mockService.
					On("DeleteDocument", tt.songID, tt.docID).
					Return(tt.mockError)
			}

			path := "/songs"
			if tt.setupParams {
				path += "/" + tt.songID + "/documents/" + tt.docID
			}

			c, w := utils.CreateTestContext(http.MethodDelete, path, nil)
			if tt.setupParams {
				if tt.songID != "" {
					c.Params = append(c.Params, gin.Param{Key: "song_id", Value: tt.songID})
				}
				if tt.docID != "" {
					c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: tt.docID})
				}
			}

			handler.DeleteDocumentHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.setupParams && tt.expectedCode != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}
