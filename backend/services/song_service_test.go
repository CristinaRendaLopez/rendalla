package services_test

import (
	"errors"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllSongs_Success(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	expectedSongs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}},
		{ID: "2", Title: "Imagine", Author: "John Lennon", Genres: []string{"Pop"}},
	}

	mockSongRepo.On("GetAllSongs").Return(expectedSongs, nil)

	songs, err := service.GetAllSongs()

	assert.NoError(t, err)
	assert.Len(t, songs, 2)
	assert.Equal(t, "Bohemian Rhapsody", songs[0].Title)
	mockSongRepo.AssertExpectations(t)
}

func TestGetAllSongs_Error(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	mockSongRepo.On("GetAllSongs").Return([]models.Song{}, errors.New("database error"))

	songs, err := service.GetAllSongs()

	assert.Error(t, err)
	assert.Empty(t, songs)
	assert.Equal(t, "database error", err.Error())
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongByID_Success(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	expectedSong := &models.Song{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}}

	mockSongRepo.On("GetSongByID", "1").Return(expectedSong, nil)

	song, err := service.GetSongByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, "Bohemian Rhapsody", song.Title)
	mockSongRepo.AssertExpectations(t)
}

func TestGetSongByID_NotFound(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	mockSongRepo.On("GetSongByID", "999").Return(nil, errors.New("song not found"))

	song, err := service.GetSongByID("999")

	assert.Error(t, err)
	assert.Nil(t, song)
	assert.Equal(t, "song not found", err.Error())
	mockSongRepo.AssertExpectations(t)
}

func TestCreateSongWithDocuments_Success(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	song := models.Song{
		Title:  "Hey Jude",
		Author: "The Beatles",
		Genres: []string{"Rock"},
	}

	documents := []models.Document{
		{
			Type:       "sheet_music",
			Instrument: []string{"Piano"},
			PDFURL:     "https://s3.amazonaws.com/beatles/heyjude.pdf",
		},
	}

	mockSongRepo.On("CreateSongWithDocuments", mock.Anything, documents).Return("123", nil)
	mockDocRepo.On("CreateDocument", mock.Anything).Return("doc1", nil)

	songID, err := service.CreateSongWithDocuments(song, documents)

	assert.NoError(t, err)
	assert.Equal(t, "123", songID)
	mockSongRepo.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
}

func TestCreateSongWithDocuments_Error(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	song := models.Song{
		Title:  "Yesterday",
		Author: "The Beatles",
		Genres: []string{"Rock"},
	}

	mockSongRepo.On("CreateSongWithDocuments", mock.Anything, mock.Anything).Return("", errors.New("database error"))

	songID, err := service.CreateSongWithDocuments(song, nil)

	assert.Error(t, err)
	assert.Empty(t, songID)
	assert.Equal(t, "database error", err.Error())
	mockSongRepo.AssertExpectations(t)
}

func TestCreateSongWithDocuments_DocumentCreationError(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	song := models.Song{
		Title:  "Stairway to Heaven",
		Author: "Led Zeppelin",
		Genres: []string{"Rock"},
	}

	documents := []models.Document{
		{
			Type:       "sheet_music",
			Instrument: []string{"Guitar"},
			PDFURL:     "https://s3.amazonaws.com/zeppelin/stairway.pdf",
		},
	}

	mockSongRepo.On("CreateSongWithDocuments", mock.Anything, documents).Return("123", nil)
	mockDocRepo.On("CreateDocument", mock.Anything).Return("", errors.New("failed to create document"))

	songID, err := service.CreateSongWithDocuments(song, documents)

	assert.Error(t, err)
	assert.Empty(t, songID)
	assert.Equal(t, "failed to create document", err.Error())

	mockSongRepo.AssertExpectations(t)
	mockDocRepo.AssertExpectations(t)
}

func TestUpdateSong_Success(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	updates := map[string]interface{}{
		"title": "Let It Be",
	}

	mockSongRepo.On("UpdateSong", "1", mock.Anything).Return(nil)

	err := service.UpdateSong("1", updates)

	assert.NoError(t, err)
	mockSongRepo.AssertExpectations(t)
}

func TestUpdateSong_Error(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	updates := map[string]interface{}{
		"title": "Let It Be",
	}

	mockSongRepo.On("UpdateSong", "999", mock.Anything).Return(errors.New("song not found"))

	err := service.UpdateSong("999", updates)

	assert.Error(t, err)
	assert.Equal(t, "song not found", err.Error())
	mockSongRepo.AssertExpectations(t)
}

func TestDeleteSongWithDocuments_Success(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	mockSongRepo.On("DeleteSongWithDocuments", "1").Return(nil)

	err := service.DeleteSongWithDocuments("1")

	assert.NoError(t, err)
	mockSongRepo.AssertExpectations(t)
}

func TestDeleteSongWithDocuments_Error(t *testing.T) {
	mockSongRepo := new(mocks.MockSongRepository)
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewSongService(mockSongRepo, mockDocRepo)

	mockSongRepo.On("DeleteSongWithDocuments", "999").Return(errors.New("song not found"))

	err := service.DeleteSongWithDocuments("999")

	assert.Error(t, err)
	assert.Equal(t, "song not found", err.Error())
	mockSongRepo.AssertExpectations(t)
}
