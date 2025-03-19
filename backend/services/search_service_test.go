package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/guregu/dynamo"
	"github.com/stretchr/testify/assert"
)

func TestSearchSongsByTitle(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	songs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}},
		{ID: "2", Title: "We Will Rock You", Author: "Queen", Genres: []string{"Rock"}},
	}

	mockSearchRepo.On("SearchSongsByTitle", "rock", 10, dynamo.PagingKey(nil)).Return(songs, dynamo.PagingKey(nil), nil)

	result, _, err := service.SearchSongsByTitle("rock", 10, dynamo.PagingKey(nil))

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Bohemian Rhapsody", result[0].Title)

	mockSearchRepo.AssertExpectations(t)
}

func TestFilterDocumentsByInstrument(t *testing.T) {
	mockSearchRepo := new(mocks.MockSearchRepository)
	service := services.NewSearchService(mockSearchRepo)

	docs := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
	}

	mockSearchRepo.On("FilterDocumentsByInstrument", "Guitar", 10, nil).Return(docs, dynamo.PagingKey(nil), nil)

	result, _, err := service.FilterDocumentsByInstrument("Guitar", 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "doc1", result[0].ID)

	mockSearchRepo.AssertExpectations(t)
}
