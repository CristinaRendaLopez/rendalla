package services_test

import (
	"fmt"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchService_ListSongs_FilterByTitle(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	songs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}},
		{ID: "2", Title: "We Will Rock You", Author: "Queen", Genres: []string{"Rock"}},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "rock", "created_at", "desc", 10, mock.Anything).
		Return(songs, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("rock", "created_at", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Bohemian Rhapsody", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_NoFilter_SortByTitleAsc(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Song{
		{ID: "1", Title: "A Kind of Magic", Author: "Queen"},
		{ID: "2", Title: "Bohemian Rhapsody", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "", "title", "asc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("", "title", "asc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "A Kind of Magic", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_FilterByTitle_SortByTitleDesc(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Song{
		{ID: "3", Title: "Somebody to Love", Author: "Queen"},
		{ID: "4", Title: "Love of My Life", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "love", "title", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("love", "title", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Somebody to Love", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_FilterByTitle_DefaultSort(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Song{
		{ID: "5", Title: "Love of My Life", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "love", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("love", "created_at", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Love of My Life", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_NoResults(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "unknown", "created_at", "desc", 10, mock.Anything).
		Return([]models.Song{}, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("unknown", "created_at", "desc", 10, next)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_WithNextToken(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	next := map[string]interface{}{"id": "last-song-id"}

	expected := []models.Song{
		{ID: "9", Title: "Keep Yourself Alive", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "", "created_at", "desc", 10, next).
		Return(expected, emptyKey, nil)

	result, _, err := service.ListSongs("", "created_at", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Keep Yourself Alive", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_RepositoryError(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	emptyKey := repository.PagingKey(map[string]interface{}{})

	mockSearchRepo.
		On("ListSongs", "queen", "created_at", "desc", 10, mock.Anything).
		Return([]models.Song{}, emptyKey, fmt.Errorf("repository failure"))

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("queen", "created_at", "desc", 10, next)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Contains(t, err.Error(), "repository failure")
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchDocumentsByTitle(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	docs := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("SearchDocumentsByTitle", "bohemian", 10, mock.Anything).
		Return(docs, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.SearchDocumentsByTitle("bohemian", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "doc1", result[0].ID)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchDocumentsByTitle_Error(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("SearchDocumentsByTitle", "error", 10, mock.Anything).
		Return([]models.Document{}, emptyKey, utils.ErrInternalServer)

	var next repository.PagingKey = nil
	result, _, err := service.SearchDocumentsByTitle("error", 10, next)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.ErrorIs(t, err, utils.ErrInternalServer)
	mockSearchRepo.AssertExpectations(t)
}

func TestFilterDocumentsByInstrument(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	docs := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("FilterDocumentsByInstrument", "Guitar", 10, nil).
		Return(docs, emptyKey, nil)

	result, _, err := service.FilterDocumentsByInstrument("Guitar", 10, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "doc1", result[0].ID)
	mockSearchRepo.AssertExpectations(t)
}

func TestFilterDocumentsByInstrument_NoResults(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("FilterDocumentsByInstrument", "Drums", 10, nil).
		Return([]models.Document{}, emptyKey, nil)

	result, _, err := service.FilterDocumentsByInstrument("Drums", 10, nil)

	assert.NoError(t, err)
	assert.Empty(t, result)
	mockSearchRepo.AssertExpectations(t)
}
