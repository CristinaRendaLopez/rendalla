package handlers_test

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDocumentHandlerTest() (*handlers.DocumentHandler, *mocks.MockDocumentService, *mocks.MockFileService, *mocks.MockFileExtractor) {
	mockService := new(mocks.MockDocumentService)
	mockFileService := new(mocks.MockFileService)
	mockFileExtractor := new(mocks.MockFileExtractor)
	handler := handlers.NewDocumentHandler(mockService, mockFileService, mockFileExtractor)
	return handler, mockService, mockFileService, mockFileExtractor
}

func TestCreateDocumentHandler(t *testing.T) {
	tests := []struct {
		name              string
		setupRequest      func() (*http.Request, *multipart.FileHeader)
		setupMock         bool
		expectSongIDParam bool
		expectOpenFile    bool
		expectUpload      bool
		expectCreate      bool
		mockUploadURL     string
		mockUploadErr     error
		mockServiceErr    error
		expectedCode      int
		expectedDocID     string
	}{
		{
			name: "success",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				return buildMultipartRequest("123", "score", []string{"guitarra"}, MockPDFFile())
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectOpenFile:    true,
			expectUpload:      true,
			expectCreate:      true,
			mockUploadURL:     "https://mock-s3/document.pdf",
			expectedCode:      http.StatusCreated,
			expectedDocID:     "abc123",
		},
		{
			name: "missing song_id",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				req, fh := buildMultipartRequest("", "score", []string{"guitarra"}, MockPDFFile())
				return req, fh
			},
			setupMock:         false,
			expectSongIDParam: false,
			expectedCode:      http.StatusBadRequest,
		},
		{
			name: "invalid pdf file",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				fh := MockInvalidFile()
				return buildMultipartRequest("123", "score", []string{"guitarra"}, fh)
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectedCode:      http.StatusBadRequest,
		},
		{
			name: "invalid document data",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				return buildMultipartRequest("123", "", []string{""}, MockPDFFile())
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectedCode:      http.StatusBadRequest,
		},
		{
			name: "validation fails after PDF upload",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				return buildMultipartRequest("123", "score", []string{}, MockPDFFile())
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectOpenFile:    true,
			expectUpload:      true,
			mockUploadURL:     "https://mock-s3/document.pdf",
			expectedCode:      http.StatusBadRequest,
		},
		{
			name: "error getting file header",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				return buildMultipartRequest("123", "score", []string{"guitarra"}, MockPDFFile())
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectedCode:      http.StatusBadRequest,
		},
		{
			name: "error opening file",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				fh := MockPDFFile()
				return buildMultipartRequest("123", "score", []string{"guitarra"}, fh)
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectOpenFile:    true,
			expectUpload:      false,
			expectedCode:      http.StatusInternalServerError,
		},
		{
			name: "upload to S3 fails",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				fh := MockPDFFile()
				return buildMultipartRequest("123", "score", []string{"guitarra"}, fh)
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectOpenFile:    true,
			expectUpload:      true,
			mockUploadErr:     fmt.Errorf("upload error"),
			expectedCode:      http.StatusInternalServerError,
		},
		{
			name: "failed to create document",
			setupRequest: func() (*http.Request, *multipart.FileHeader) {
				return buildMultipartRequest("123", "score", []string{"guitarra"}, MockPDFFile())
			},
			setupMock:         true,
			expectSongIDParam: true,
			expectOpenFile:    true,
			expectUpload:      true,
			expectCreate:      true,
			mockUploadURL:     "https://mock-s3/document.pdf",
			mockServiceErr:    fmt.Errorf("db error"),
			expectedCode:      http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockDocService, mockFileService, mockFileExtractor := setupDocumentHandlerTest()

			req, fileHeader := tt.setupRequest()
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tt.expectSongIDParam {
				c.Params = []gin.Param{{Key: "song_id", Value: "123"}}
			}

			if tt.setupMock {
				if tt.name == "error getting file header" {
					mockFileExtractor.On("GetHeader", c, "pdf").Return(nil, fmt.Errorf("missing file"))
				} else {
					mockFileExtractor.On("GetHeader", c, "pdf").Return(fileHeader, nil)

					if tt.name == "invalid document data" {
						return
					}

					if tt.expectOpenFile {
						if tt.expectUpload {
							mockFileExtractor.On("OpenFile", fileHeader).Return(mocks.NewFakeMultipartFile([]byte("PDF")), nil)
						} else {
							mockFileExtractor.On("OpenFile", fileHeader).Return(nil, fmt.Errorf("failed to open"))
						}
					}

					if tt.expectUpload {
						mockFileService.On("UploadPDFToS3", mock.Anything, fileHeader, "123").Return(tt.mockUploadURL, tt.mockUploadErr)
					}

					if tt.expectCreate {
						mockDocService.On("CreateDocument", mock.Anything).Return(tt.expectedDocID, tt.mockServiceErr)
					}
				}
			}

			handler.CreateDocumentHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusCreated {
				var resp dto.CreateDocumentResponse
				err := json.Unmarshal(w.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDocID, resp.DocumentID)
			}

			mockDocService.AssertExpectations(t)
			mockFileService.AssertExpectations(t)
			mockFileExtractor.AssertExpectations(t)
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
			handler, mockService, _, _ := setupDocumentHandlerTest()

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
			handler, mockService, _, _ := setupDocumentHandlerTest()

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
			handler, mockService, _, _ := setupDocumentHandlerTest()

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
			handler, mockService, _, _ := setupDocumentHandlerTest()

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
