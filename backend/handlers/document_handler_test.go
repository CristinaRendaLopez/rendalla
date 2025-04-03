package handlers_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
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

	documents := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", PDFURL: "https://example.com/doc1.pdf"},
	}

	mockService.On("GetDocumentsBySongID", "1").Return(documents, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1/documents", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "doc1")
	mockService.AssertExpectations(t)
}

func TestGetAllDocumentsBySongIDHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("GetDocumentsBySongID", "1").Return([]models.Document{}, utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1/documents", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetAllDocumentsBySongIDHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/", nil)
	handler.GetAllDocumentsBySongIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing id")
}

func TestGetDocumentByIDHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	document := &models.Document{ID: "doc1", SongID: "1", Type: "sheet_music", PDFURL: "https://example.com/doc1.pdf"}
	mockService.On("GetDocumentByID", "doc1").Return(document, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc1"})

	handler.GetDocumentByIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "doc1")
	mockService.AssertExpectations(t)
}

func TestGetDocumentByIDHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/", nil)
	handler.GetDocumentByIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing id")
}

func TestGetDocumentByIDHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("GetDocumentByID", "doc1").Return((*models.Document)(nil), utils.ErrResourceNotFound)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc1"})

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
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Document created successfully")
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
	assert.Contains(t, w.Body.String(), "Missing id")
}

func TestCreateDocumentHandler_InvalidInput(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	invalidDoc := `{"type": "", "instrument": [], "pdf_url": ""}`
	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(invalidDoc))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid document data")
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_ServiceError(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	validDoc := `{"type": "score", "instrument": ["violin"], "pdf_url": "file.pdf"}`
	mockService.On("CreateDocument", mock.Anything).Return("", utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(validDoc))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to create document")
	mockService.AssertExpectations(t)
}

func TestCreateDocumentHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	badJSON := `{"type":`

	c, w := utils.CreateTestContext(http.MethodPost, "/songs/1/documents", strings.NewReader(badJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.CreateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request data")
}

func TestUpdateDocumentHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	updates := map[string]interface{}{"type": "tablature"}
	mockService.On("UpdateDocument", "doc123", updates).Return(nil)

	c, w := utils.CreateTestContext(http.MethodPut, "/documents/doc123", strings.NewReader(`{"type": "tablature"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc123"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Document updated successfully")
	mockService.AssertExpectations(t)
}

func TestUpdateDocumentHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	update := `{"type":"score"}`
	c, w := utils.CreateTestContext(http.MethodPut, "/documents/", strings.NewReader(update))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing id")
}

func TestUpdateDocumentHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	badJSON := `{"type":`

	c, w := utils.CreateTestContext(http.MethodPut, "/documents/1", strings.NewReader(badJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request data")
}

func TestUpdateDocumentHandler_InvalidType(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/documents/doc123", strings.NewReader(`{"type": ""}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc123"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid document update data")
}

func TestUpdateDocumentHandler_ServiceError(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	updates := map[string]interface{}{"type": "score"}
	mockService.On("UpdateDocument", "doc123", updates).Return(utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPut, "/documents/doc123", strings.NewReader(`{"type": "score"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc123"})

	handler.UpdateDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to update document")
	mockService.AssertExpectations(t)
}

func TestDeleteDocumentHandler_Success(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("DeleteDocument", "doc1").Return(nil)

	c, w := utils.CreateTestContext(http.MethodDelete, "/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc1"})

	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Document deleted successfully")
	mockService.AssertExpectations(t)
}

func TestDeleteDocumentHandler_MissingID(t *testing.T) {
	handler, _ := setupDocumentHandlerTest()

	c, w := utils.CreateTestContext(http.MethodDelete, "/documents/", nil)
	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing id")
}

func TestDeleteDocumentHandler_Service_Error(t *testing.T) {
	handler, mockService := setupDocumentHandlerTest()

	mockService.On("DeleteDocument", "doc1").Return(utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodDelete, "/documents/doc1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "doc1"})

	handler.DeleteDocumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to delete document")
	mockService.AssertExpectations(t)
}
