package handlers_test

import (
	"net/http"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
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

func TestListSongsHandler_FilterByTitle_DefaultSort(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "love", "", "", 10, mock.Anything).Return([]models.Song{
		{ID: "1", Title: "Love of My Life", Author: "Queen"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=love", nil)
	c.Request.URL.RawQuery = "title=love"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":\"1\"")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_TitleWithSortDesc(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "love", "title", "desc", 10, mock.Anything).Return([]models.Song{
		{ID: "2", Title: "Somebody to Love", Author: "Queen"},
		{ID: "1", Title: "Love of My Life", Author: "Queen"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=love&sort=title&order=desc", nil)
	c.Request.URL.RawQuery = "title=love&sort=title&order=desc"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Love of My Life")
	assert.Contains(t, w.Body.String(), "Somebody to Love")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_SortByCreatedAtAsc(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "", "created_at", "asc", 10, mock.Anything).Return([]models.Song{
		{ID: "3", Title: "Seven Seas of Rhye", Author: "Queen", CreatedAt: "1974-01-01T00:00:00Z"},
		{ID: "4", Title: "Radio Ga Ga", Author: "Queen", CreatedAt: "1984-01-01T00:00:00Z"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?sort=created_at&order=asc", nil)
	c.Request.URL.RawQuery = "sort=created_at&order=asc"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Seven Seas of Rhye")
	assert.Contains(t, w.Body.String(), "Radio Ga Ga")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_DefaultParams(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "", "", "", 10, mock.Anything).Return([]models.Song{
		{ID: "5", Title: "The Show Must Go On", Author: "Queen"},
		{ID: "6", Title: "I Want It All", Author: "Queen"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search", nil)
	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "The Show Must Go On")
	assert.Contains(t, w.Body.String(), "I Want It All")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_EmptyResult(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "radio", "", "", 10, mock.Anything).Return([]models.Song{}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=radio", nil)
	c.Request.URL.RawQuery = "title=radio"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"data":[]`)
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_SortByTitleAsc_NoFilter(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "", "title", "asc", 10, mock.Anything).Return([]models.Song{
		{ID: "2", Title: "Bohemian Rhapsody", Author: "Queen"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?sort=title&order=asc", nil)
	c.Request.URL.RawQuery = "sort=title&order=asc"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Bohemian Rhapsody")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_NextTokenPresent(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "", "", "", 10, mock.Anything).Return([]models.Song{
		{ID: "4", Title: "One Vision", Author: "Queen"},
	}, "some-token", nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?next_token=some-token", nil)
	c.Request.URL.RawQuery = "next_token=some-token"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "One Vision")
	assert.Contains(t, w.Body.String(), "next_token")
	mockService.AssertExpectations(t)
}

func TestListSongsHandler_ServiceError(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListSongs", "love", "", "", 10, mock.Anything).Return([]models.Song{}, nil, errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs/search?title=love", nil)
	c.Request.URL.RawQuery = "title=love"

	handler.ListSongsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_FilterByTitle_DefaultSort(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "queen", "", "", "", "", 10, mock.Anything).Return([]models.Document{
		{ID: "1", SongID: "s1", TitleNormalized: "queen", Type: "sheet_music", Instrument: []string{"Guitar"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=queen", nil)
	c.Request.URL.RawQuery = "title=queen"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "sheet_music")
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_FilterByInstrument(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "", "Piano", "", "", "", 10, mock.Anything).Return([]models.Document{
		{ID: "2", SongID: "s2", TitleNormalized: "bohemian rhapsody", Type: "sheet_music", Instrument: []string{"Piano"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?instrument=Piano", nil)
	c.Request.URL.RawQuery = "instrument=Piano"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":\"2\"")
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_FilterByType(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "", "", "tablature", "", "", 10, mock.Anything).Return([]models.Document{
		{ID: "3", SongID: "s3", TitleNormalized: "we will rock you", Type: "tablature", Instrument: []string{"Guitar"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?type=tablature", nil)
	c.Request.URL.RawQuery = "type=tablature"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "tablature")
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_CombinedFiltersWithSorting(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "love", "Violin", "sheet_music", "title", "asc", 10, mock.Anything).Return([]models.Document{
		{ID: "3", SongID: "s3", TitleNormalized: "love of my life", Type: "sheet_music", Instrument: []string{"Violin"}},
		{ID: "4", SongID: "s4", TitleNormalized: "somebody to love", Type: "sheet_music", Instrument: []string{"Violin"}},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=love&instrument=Violin&type=sheet_music&sort=title&order=asc", nil)
	c.Request.URL.RawQuery = "title=love&instrument=Violin&type=sheet_music&sort=title&order=asc"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":\"3\"")
	assert.Contains(t, w.Body.String(), "\"id\":\"4\"")

	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_SortByTitleAsc(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "", "", "", "title", "asc", 10, mock.Anything).Return([]models.Document{
		{ID: "7", SongID: "s7", TitleNormalized: "a kind of magic", Type: "sheet_music"},
		{ID: "8", SongID: "s8", TitleNormalized: "bohemian rhapsody", Type: "sheet_music"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?sort=title&order=asc", nil)
	c.Request.URL.RawQuery = "sort=title&order=asc"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":\"7\"")
	assert.Contains(t, w.Body.String(), "\"id\":\"8\"")
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_SortByCreatedAtDesc(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "", "", "", "created_at", "desc", 10, mock.Anything).Return([]models.Document{
		{ID: "5", SongID: "s5", TitleNormalized: "under pressure", CreatedAt: "1985-01-01T00:00:00Z"},
		{ID: "6", SongID: "s6", TitleNormalized: "innuendo", CreatedAt: "1991-01-01T00:00:00Z"},
	}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?sort=created_at&order=desc", nil)
	c.Request.URL.RawQuery = "sort=created_at&order=desc"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "\"id\":\"5\"")
	assert.Contains(t, w.Body.String(), "\"id\":\"6\"")

	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_EmptyResult(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "nonexistent", "", "", "", "", 10, mock.Anything).Return([]models.Document{}, nil, nil)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=nonexistent", nil)
	c.Request.URL.RawQuery = "title=nonexistent"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"data":[]`)
	mockService.AssertExpectations(t)
}

func TestListDocumentsHandler_ServiceError(t *testing.T) {
	handler, mockService := setupSearchHandlerTest()

	mockService.On("ListDocuments", "queen", "", "", "", "", 10, mock.Anything).
		Return([]models.Document{}, nil, errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodGet, "/documents/search?title=queen", nil)
	c.Request.URL.RawQuery = "title=queen"

	handler.ListDocumentsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.ErrorIs(t, errors.ErrInternalServer, errors.ErrInternalServer)
	mockService.AssertExpectations(t)
}
