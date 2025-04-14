package handlers_test

import (
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
	"github.com/stretchr/testify/mock"
)

func setupDocumentHandlerTest() (*handlers.DocumentHandler, *mocks.MockDocumentService) {
	mockService := new(mocks.MockDocumentService)
	handler := handlers.NewDocumentHandler(mockService)
	return handler, mockService
}

func TestGetAllDocumentsBySongIDHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	documents := []dto.DocumentResponseItem{
		{
			ID:     "doc1",
			SongID: "1",
			Type:   "sheet_music",
			PDFURL: "https://example.com/doc1.pdf",
		},
	}

	mockService.On("GetDocumentsBySongID", "1").Return(documents, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1/documents", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "doc1")
	mockService.AssertExpectations(t)
}

func TestGetAllDocumentsBySongIDHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("GetDocumentsBySongID", "1").Return([]dto.DocumentResponseItem{}, errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1/documents", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllDocumentsBySongIDHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/", nil)
	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDocumentByIDHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	document := dto.DocumentResponseItem{
		ID:     "doc1",
		SongID: "queen-001",
		Type:   "sheet_music",
		PDFURL: "https://example.com/doc1.pdf",
	}

	mockService.On("GetDocumentByID", "queen-001", "doc1").Return(document, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/queen-001/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc1"})

	handler.GetDocumentByIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "doc1")
	mockService.AssertExpectations(t)
}

func TestGetDocumentByIDHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/queen-001/documents/", nil)
	handler.GetDocumentByIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDocumentByIDHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("GetDocumentByID", "queen-001", "doc1").Return(dto.DocumentResponseItem{}, errors.ErrResourceNotFound)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/queen-001/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc1"})

	handler.GetDocumentByIDHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	validDoc := `{
		"type": "score",
		"instrument": ["guitar"],
		"pdf_url": "https://s3.amazonaws.com/bucket/file.pdf"
	}`

	mockService.On("CreateDocument", mock.Anything).Return("doc123", nil)

	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(validDoc))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "doc123")
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_MissingSongID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	validDoc := `{"type":"score", "instrument":["guitar"], "pdf_url":"file.pdf"}`
	c, w := utils.CreateTestContext(http.MethodPost, "/songs//documents", strings.NewReader(validDoc))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateDocumentHandler_InvalidInput(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	invalidDoc := `{"type": "", "instrument": [], "pdf_url": ""}`
	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(invalidDoc))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_ServiceError(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	validDoc := `{
		"type": "score",
		"instrument": ["violin"],
		"pdf_url": "https://example.com/file.pdf"
	}`
	mockService.On("CreateDocument", mock.Anything).Return("", errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(validDoc))
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	badJSON := `{"type":`

	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(badJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDocumentHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	updates := dto.UpdateDocumentRequest{
		Type: "tablature",
	}
	mockService.On("UpdateDocument", "queen-001", "doc123", updates).Return(nil)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/queen-001/documents/doc123", strings.NewReader(`{"type": "tablature"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc123"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateDocumentHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	update := `{"type":"score"}`
	c, w := utils.CreateTestContext(http.MethodPut, "/songs/queen-001/documents/", strings.NewReader(update))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDocumentHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	badJSON := `{"type":`

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/queen-001/documents/1", strings.NewReader(badJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "1"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateDocumentHandler_ServiceError(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	updates := dto.UpdateDocumentRequest{
		Type: "score",
	}
	mockService.On("UpdateDocument", "queen-001", "doc123", updates).Return(errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/queen-001/documents/doc123", strings.NewReader(`{"type": "score"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc123"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteDocumentHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("DeleteDocument", "queen-001", "doc1").Return(nil)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/queen-001/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc1"})

	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteDocumentHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/queen-001/documents/", nil)
	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteDocumentHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("DeleteDocument", "queen-001", "doc1").Return(errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/queen-001/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "queen-001"})
	c.Params = append(c.Params, gin.Param{Key: "doc_id", Value: "doc1"})

	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}
