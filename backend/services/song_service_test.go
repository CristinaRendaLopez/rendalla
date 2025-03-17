package services_test

import (
	"testing"
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
)

func TestGetSongByID(t *testing.T) {
	mockDB := new(services.MockDB)
	expectedSong := &models.Song{
		ID:        "123",
		Title:     "Bohemian Rhapsody",
		Author:    "Queen",
		Genres:    []string{"Rock", "Opera"},
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}

	mockDB.On("GetSongByID", "123").Return(expectedSong, nil)

	song, err := mockDB.GetSongByID("123")

	assert.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, expectedSong, song)
}
