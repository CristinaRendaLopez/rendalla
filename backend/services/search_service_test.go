package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
)

func TestSearchSongsByTitle(t *testing.T) {
	mockDB := new(services.MockDB)

	songs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}},
		{ID: "2", Title: "We Will Rock You", Author: "Queen", Genres: []string{"Rock"}},
	}

	mockDB.On("SearchSongsByTitle", "rock", 10, nil).Return(songs, nil, nil)

	result, _, err := mockDB.SearchSongsByTitle("rock", 10, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Bohemian Rhapsody", result[0].Title)
}

func TestSearchSongsByTitle_NoResults(t *testing.T) {
	mockDB := new(services.MockDB)

	mockDB.On("SearchSongsByTitle", "nonexistent", 10, nil).Return([]models.Song{}, nil, nil)

	result, _, err := mockDB.SearchSongsByTitle("nonexistent", 10, nil)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(result))
}

func TestSearchDocumentsByTitle(t *testing.T) {
	mockDB := new(services.MockDB)

	docs := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
		{ID: "doc2", SongID: "2", Type: "tablature", Instrument: []string{"Drums"}, PDFURL: "https://s3.amazonaws.com/queen/wewillrockyou.pdf"},
	}

	mockDB.On("SearchDocumentsByTitle", "bohemian", 10, nil).Return(docs, nil, nil)

	result, _, err := mockDB.SearchDocumentsByTitle("bohemian", 10, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "doc1", result[0].ID)
}

func TestFilterDocumentsByInstrument(t *testing.T) {
	mockDB := new(services.MockDB)

	docs := []models.Document{
		{ID: "doc1", SongID: "1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
	}

	mockDB.On("FilterDocumentsByInstrument", "Guitar", 10, nil).Return(docs, nil, nil)

	result, _, err := mockDB.FilterDocumentsByInstrument("Guitar", 10, nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "doc1", result[0].ID)
}
