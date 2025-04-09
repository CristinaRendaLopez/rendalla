package services_test

import (
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

func TestSearchService_ListSongs_InvalidSortField_ShouldFallbackToCreatedAt(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Song{
		{ID: "6", Title: "One Vision", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "queen", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("queen", "banana", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "One Vision", result[0].Title)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListSongs_InvalidOrder_ShouldFallbackToDesc(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Song{
		{ID: "7", Title: "Innuendo", Author: "Queen"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListSongs", "innuendo", "title", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("innuendo", "title", "sideways", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Innuendo", result[0].Title)
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
		Return([]models.Song{}, emptyKey, utils.ErrInternalServer)

	var next repository.PagingKey = nil
	result, _, err := service.ListSongs("queen", "created_at", "desc", 10, next)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.ErrorIs(t, err, utils.ErrInternalServer)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_FilterByTitle(t *testing.T) {
	mockRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockRepo)

	expected := []models.Document{
		{ID: "1", TitleNormalized: "love of my life", Type: "sheet_music"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockRepo.
		On("ListDocuments", "love", "", "", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListDocuments("love", "", "", "created_at", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "love of my life", result[0].TitleNormalized)
	mockRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_FilterByInstrument(t *testing.T) {
	mockRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockRepo)

	expected := []models.Document{
		{ID: "2", Instrument: []string{"Guitar"}, Type: "tablature"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockRepo.
		On("ListDocuments", "", "Guitar", "", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	result, _, err := service.ListDocuments("", "Guitar", "", "created_at", "desc", 10, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "tablature", result[0].Type)
	mockRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_FilterByType(t *testing.T) {
	mockRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockRepo)

	expected := []models.Document{
		{ID: "3", Type: "sheet_music"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockRepo.
		On("ListDocuments", "", "", "sheet_music", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	result, _, err := service.ListDocuments("", "", "sheet_music", "created_at", "desc", 10, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "sheet_music", result[0].Type)
	mockRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_CombinedFilters(t *testing.T) {
	mockRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockRepo)

	expected := []models.Document{
		{ID: "4", TitleNormalized: "somebody to love", Type: "sheet_music", Instrument: []string{"Violin"}},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockRepo.
		On("ListDocuments", "love", "Violin", "sheet_music", "title", "asc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	result, _, err := service.ListDocuments("love", "Violin", "sheet_music", "title", "asc", 10, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "somebody to love", result[0].TitleNormalized)
	mockRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_InvalidSortField_ShouldFallbackToCreatedAt(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Document{
		{ID: "6", TitleNormalized: "one vision", Type: "sheet_music"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListDocuments", "queen", "", "", "created_at", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	var next repository.PagingKey = nil
	result, _, err := service.ListDocuments("queen", "", "", "banana", "desc", 10, next)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "one vision", result[0].TitleNormalized)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_InvalidSortOrder_ShouldFallbackToDesc(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	expected := []models.Document{
		{ID: "7", TitleNormalized: "innuendo", Type: "tablature"},
	}

	emptyKey := repository.PagingKey(map[string]interface{}{})
	mockSearchRepo.
		On("ListDocuments", "innuendo", "", "", "title", "desc", 10, mock.Anything).
		Return(expected, emptyKey, nil)

	result, _, err := service.ListDocuments("innuendo", "", "", "title", "sideways", 10, nil)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "innuendo", result[0].TitleNormalized)
	mockSearchRepo.AssertExpectations(t)
}

func TestSearchService_ListDocuments_RepositoryError(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	emptyKey := repository.PagingKey(map[string]interface{}{})

	mockSearchRepo.
		On("ListDocuments", "queen", "", "", "created_at", "desc", 10, mock.Anything).
		Return([]models.Document{}, emptyKey, utils.ErrInternalServer)

	var next repository.PagingKey = nil
	result, _, err := service.ListDocuments("queen", "", "", "created_at", "desc", 10, next)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.ErrorIs(t, err, utils.ErrInternalServer)
	mockSearchRepo.AssertExpectations(t)
}
