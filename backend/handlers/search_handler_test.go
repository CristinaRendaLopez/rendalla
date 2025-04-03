package handlers_test

import (
	"net/http"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSearchHandlerTest() (*handlers.SearchHandler, *mocks.MockSearchService) {
	mockService := new(mocks.MockSearchService)
	handler := handlers.NewSearchHandler(mockService)
	return handler, mockService
}

func TestSearchSongsByTitleHandler_Success(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("SearchSongsByTitle", "love", 10, mock.Anything).Return([]models.Song{
		{ID: "1", Title: "Love Me Do"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=love", nil)
	c.Request.URL.RawQuery = "title=love"

	handler.SearchSongsByTitleHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Love Me Do")
	mockService.AssertExpectations(t)
}

func TestSearchSongsByTitleHandler_MissingTitle(t *testing.T) {
	handler, _ := setupSearchHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search", nil)
	c.Request.URL.RawQuery = "title="

	handler.SearchSongsByTitleHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing title parameter")
}

func TestSearchSongsByTitleHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("SearchSongsByTitle", "love", 10, mock.Anything).Return([]models.Song{}, nil, utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=love", nil)
	c.Request.URL.RawQuery = "title=love"

	handler.SearchSongsByTitleHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error searching for songs")
	mockService.AssertExpectations(t)
}

func TestSearchDocumentsByTitleHandler_Success(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("SearchDocumentsByTitle", "score", 10, mock.Anything).Return([]models.Document{
		{ID: "doc1", Type: "score", Instrument: []string{"violin"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=score", nil)
	c.Request.URL.RawQuery = "title=score"

	handler.SearchDocumentsByTitleHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "violin")
	mockService.AssertExpectations(t)
}

func TestSearchDocumentsByTitleHandler_MissingTitle(t *testing.T) {
	handler, _ := setupSearchHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search", nil)
	c.Request.URL.RawQuery = "title="

	handler.SearchDocumentsByTitleHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing title parameter")
}

func TestFilterDocumentsByInstrumentHandler_Success(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("FilterDocumentsByInstrument", "guitar", 10, mock.Anything).Return([]models.Document{
		{ID: "doc1", Type: "score", Instrument: []string{"guitar"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/filter?instrument=guitar", nil)
	c.Request.URL.RawQuery = "instrument=guitar"

	handler.FilterDocumentsByInstrumentHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "guitar")
	mockService.AssertExpectations(t)
}

func TestSearchDocumentsByTitleHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("SearchDocumentsByTitle", "score", 10, mock.Anything).Return([]models.Document{}, nil, utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=score", nil)
	c.Request.URL.RawQuery = "title=score"

	handler.SearchDocumentsByTitleHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error searching for documents")
	mockService.AssertExpectations(t)
}

func TestFilterDocumentsByInstrumentHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("FilterDocumentsByInstrument", "guitar", 10, mock.Anything).Return([]models.Document{}, nil, utils.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/filter?instrument=guitar", nil)
	c.Request.URL.RawQuery = "instrument=guitar"

	handler.FilterDocumentsByInstrumentHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error filtering documents by instrument")
	mockService.AssertExpectations(t)
}

func TestFilterDocumentsByInstrumentHandler_MissingInstrument(t *testing.T) {
	handler, _ := setupSearchHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/filter", nil)
	c.Request.URL.RawQuery = "instrument="

	handler.FilterDocumentsByInstrumentHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing instrument parameter")
}

func TestExtractPaginationParams_InvalidLimit(t *testing.T) {
	c, _ := utils.CreateTestContext(http.MethodGet, "/test", nil)
	c.Request.URL.RawQuery = "limit=abc"

	limit, next := utils.ExtractPaginationParams(c)

	assert.Equal(t, 10, limit)
	assert.Equal(t, 0, len(next))
}

func TestExtractPaginationParams_WithNextToken(t *testing.T) {
	c, _ := utils.CreateTestContext(http.MethodGet, "/test", nil)
	c.Request.URL.RawQuery = "next_token=abc123"

	limit, next := utils.ExtractPaginationParams(c)

	assert.Equal(t, 10, limit)
	_, ok := next["abc123"]
	assert.True(t, ok)
}
