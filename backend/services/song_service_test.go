package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSongWithDocuments(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	song := models.Song{
		Title:  "Bohemian Rhapsody",
		Author: "Queen",
		Genres: []string{"Rock", "Opera"},
	}

	documents := []models.Document{
		{
			Type:       "sheet_music",
			Instrument: []string{"Guitar"},
			PDFURL:     "https://s3.amazonaws.com/queen/bohemian.pdf",
		},
	}

	mockSongRepo.On("CreateSongWithDocuments", mock.MatchedBy(func(s models.Song) bool {
		return s.Title == "Bohemian Rhapsody" &&
			s.Author == "Queen" &&
			len(s.Genres) == 2 &&
			s.Genres[0] == "Rock" &&
			s.Genres[1] == "Opera"
	}), documents).Return("123", nil)

	mockDocRepo.On("CreateDocument", mock.MatchedBy(func(d models.Document) bool {
		return d.Type == "sheet_music" &&
			d.PDFURL == "https://s3.amazonaws.com/queen/bohemian.pdf" &&
			d.SongID == "123"
	})).Return("doc1", nil)

	songID, err := service.CreateSongWithDocuments(song, documents)

	assert.NoError(t, err)
	assert.Equal(t, "123", songID)

	mockSongRepo.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
}
