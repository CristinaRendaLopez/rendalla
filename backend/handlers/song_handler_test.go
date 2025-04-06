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

func setupSongHandlerTest() (*handlers.SongHandler, *mocks.MockSongService) {
	mockService := new(mocks.MockSongService)
	handler := handlers.NewSongHandler(mockService)
	return handler, mockService
}

func TestGetAllSongsHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	expectedSongs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen"},
		{ID: "2", Title: "Imagine", Author: "John Lennon"},
	}
	mockService.On("GetAllSongs").Return(expectedSongs, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs", nil)
	handler.GetAllSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bohemian Rhapsody")
	mockService.AssertExpectations(t)
}

func TestGetAllSongsHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("GetAllSongs").Return([]models.Song{}, assert.AnError)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs", nil)
	handler.GetAllSongsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetSongByIDHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("GetSongByID", "1").Return(&models.Song{ID: "1", Title: "Song"}, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.GetSongByIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Song")
	mockService.AssertExpectations(t)
}

func TestGetSongByIDHandler_NotFound(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("GetSongByID", "999").Return((*models.Song)(nil), utils.ErrResourceNotFound)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/999", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.GetSongByIDHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetSongByIDHandler_MissingID(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/", nil)
	handler.GetSongByIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSongHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	data := `{"title": "Yesterday", "author": "The Beatles", "genres": ["Rock"]}`
	mockService.On("CreateSongWithDocuments", mock.Anything, mock.Anything).Return("123", nil)

	c, w := utils.CreateTestContext(http.MethodPost, "/songs", strings.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSongHandler(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateSongHandler_InvalidInput(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	data := `{"title": "", "author": "", "genres": []}`
	c, w := utils.CreateTestContext(http.MethodPost, "/songs", strings.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSongHandler_InvalidDocument(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	data := `{"title": "Title", "author": "Auth", "genres": ["Pop"], "documents":[{"type": "", "instrument": [], "pdf_url": ""}]}`
	c, w := utils.CreateTestContext(http.MethodPost, "/songs", strings.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSongHandler_ServiceError(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	data := `{"title": "Imagine", "author": "John Lennon", "genres": ["Rock"]}`
	mockService.On("CreateSongWithDocuments", mock.Anything, mock.Anything).Return("", utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPost, "/songs", strings.NewReader(data))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateSongHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateSongHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	update := map[string]interface{}{"title": "New Title"}
	mockService.On("UpdateSong", "1", update).Return(nil)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":"New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateSongHandler_MissingID(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/", strings.NewReader(`{"title": "New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidJSON(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title": 123}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_EmptyUpdate(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidFields(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":""}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	invalidJSON := `{"title":`

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_ServiceError(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	update := map[string]interface{}{"title": "New Title"}
	mockService.On("UpdateSong", "1", update).Return(utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":"New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("DeleteSongWithDocuments", "1").Return(nil)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "1"})

	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("DeleteSongWithDocuments", "999").Return(utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/999", nil)
	c.Params = append(c.Params, gin.Param{Key: "id", Value: "999"})

	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_MissingID(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/", nil)
	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
