package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDocumentHandlerTest() (*handlers.DocumentHandler, *mocks.MockDocumentService, utils.FileUploader) {
	mockService := new(mocks.MockDocumentService)
	mockFileService := new(mocks.MockFileService)
	handler := handlers.NewDocumentHandler(mockService, mockFileService)
	return handler, mockService, mockFileService
}

func TestCreateDocumentHandler(t *testing.T) {
	fmt.Println("MAX_PDF_SIZE from env:", os.Getenv("MAX_PDF_SIZE"))
	fmt.Println("bootstrap.MaxPDFSize:", bootstrap.MaxPDFSize)
	tests := []struct {
		name          string
		songID        string
		setupParam    bool
		fields        map[string]string
		file          utils.TestFile
		expectedCode  int
		mockReturnID  string
		mockReturnErr error
		expectedDocID string
		expectUpload  bool
		expectCreate  bool
	}{
		{
			name:          "successfully creates document",
			songID:        "1",
			setupParam:    true,
			fields:        MultipartFieldsValid,
			file:          MultipartPDFMock,
			expectedCode:  http.StatusCreated,
			mockReturnID:  "doc-123",
			expectedDocID: "doc-123",
			expectUpload:  true,
			expectCreate:  true,
		},
		{
			name:         "missing song_id parameter",
			setupParam:   false,
			expectedCode: http.StatusBadRequest,
			expectUpload: false,
			expectCreate: false,
		},
		{
			name:         "validation fails (empty fields)",
			songID:       "1",
			setupParam:   true,
			fields:       MultipartFieldsInvalid,
			file:         MultipartPDFMock,
			expectedCode: http.StatusBadRequest,
			expectUpload: true,
			expectCreate: false,
		},
		{
			name:         "validation fails (invalid PDF)",
			songID:       "1",
			setupParam:   true,
			fields:       MultipartFieldsValid,
			file:         MultipartPDFInvalid,
			expectedCode: http.StatusBadRequest,
			expectUpload: false,
			expectCreate: false,
		},
		{
			name:         "file service upload fails",
			songID:       "1",
			setupParam:   true,
			fields:       MultipartFieldsValid,
			file:         MultipartPDFMock,
			expectedCode: http.StatusInternalServerError,
			expectUpload: true,
			expectCreate: false,
		},
		{
			name:          "internal service error",
			songID:        "1",
			setupParam:    true,
			fields:        MultipartFieldsValid,
			file:          MultipartPDFMock,
			expectedCode:  http.StatusInternalServerError,
			mockReturnErr: errors.ErrInternalServer,
			expectUpload:  true,
			expectCreate:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockDocumentService)
			mockFileService := new(mocks.MockFileService)

			if tt.expectUpload {
				if tt.name == "file service upload fails" {
					mockFileService.
						On("UploadPDFToS3", mock.Anything, mock.Anything, tt.songID).
						Return("", errors.ErrInternalServer)
				} else {
					mockFileService.
						On("UploadPDFToS3", mock.Anything, mock.Anything, tt.songID).
						Return("https://mock.url/file.pdf", nil)
				}
			}

			if tt.expectCreate {
				mockService.
					On("CreateDocument", mock.AnythingOfType("dto.CreateDocumentRequest")).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}

			handler := handlers.NewDocumentHandler(mockService, mockFileService)

			var body *bytes.Buffer
			var contentType string
			var err error

			if tt.fields != nil && tt.file.Filename != "" {
				body, contentType, err = utils.CreateMultipartRequest(tt.fields, map[string]utils.TestFile{
					"pdf": tt.file,
				})
				assert.NoError(t, err)
			} else {
				body = &bytes.Buffer{}
			}

			path := "/songs"
			if tt.setupParam {
				path += "/" + tt.songID + "/documents"
			}

			c, w := utils.CreateTestContext(http.MethodPost, path, body)
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}
			if contentType != "" {
				c.Request.Header.Set("Content-Type", contentType)
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

			if tt.expectUpload {
				mockFileService.AssertCalled(t, "UploadPDFToS3", mock.Anything, mock.Anything, tt.songID)
			} else {
				mockFileService.AssertNotCalled(t, "UploadPDFToS3", mock.Anything, mock.Anything, tt.songID)
			}

			if tt.expectCreate {
				mockService.AssertCalled(t, "CreateDocument", mock.AnythingOfType("dto.CreateDocumentRequest"))
			} else {
				mockService.AssertNotCalled(t, "CreateDocument", mock.Anything)
			}

			mockFileService.AssertExpectations(t)
			mockService.AssertExpectations(t)
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
			handler, mockService, _ := setupDocumentHandlerTest()

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
			handler, mockService, _ := setupDocumentHandlerTest()

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
			handler, mockService, _ := setupDocumentHandlerTest()

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
			handler, mockService, _ := setupDocumentHandlerTest()

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
