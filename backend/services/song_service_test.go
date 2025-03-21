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

func setupSongServiceTest() (*services.SongService, *mocks.MockSongRepository, *mocks.MockDocumentRepository, *mocks.MockIDGenerator, *mocks.MockTimeProvider) {
	songRepo := new(mocks.MockSongRepository)
	docRepo := new(mocks.MockDocumentRepository)
	idGen := new(mocks.MockIDGenerator)
	timeProv := new(mocks.MockTimeProvider)
	service := services.NewSongService(songRepo, docRepo, idGen, timeProv)
	return service, songRepo, docRepo, idGen, timeProv
}

func TestGetAllSongs_Success(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	expectedSongs := []models.Song{
		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}},
		{ID: "2", Title: "Imagine", Author: "John Lennon", Genres: []string{"Pop"}},
	}

	songRepo.On("GetAllSongs").Return(expectedSongs, nil)

	songs, err := service.GetAllSongs()

	assert.NoError(t, err)
	assert.Len(t, songs, 2)
	assert.Equal(t, "Bohemian Rhapsody", songs[0].Title)
	songRepo.AssertExpectations(t)
}

func TestGetAllSongs_Error(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	songRepo.On("GetAllSongs").Return([]models.Song{}, errors.New("database error"))

	songs, err := service.GetAllSongs()

	assert.Error(t, err)
	assert.Empty(t, songs)
	assert.Equal(t, "database error", err.Error())
	songRepo.AssertExpectations(t)
}

func TestGetSongByID_Success(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	expectedSong := &models.Song{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen", Genres: []string{"Rock"}}

	songRepo.On("GetSongByID", "1").Return(expectedSong, nil)

	song, err := service.GetSongByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, song)
	assert.Equal(t, "Bohemian Rhapsody", song.Title)
	songRepo.AssertExpectations(t)
}

func TestGetSongByID_NotFound(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	songRepo.On("GetSongByID", "999").Return(nil, errors.New("song not found"))

	song, err := service.GetSongByID("999")

	assert.Error(t, err)
	assert.Nil(t, song)
	assert.Equal(t, "song not found", err.Error())
	songRepo.AssertExpectations(t)
}

func TestCreateSongWithDocuments_Success(t *testing.T) {
	service, songRepo, docRepo, idGen, timeProv := setupSongServiceTest()

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

	idGen.On("NewID").Return("song-123").Once()
	timeProv.On("Now").Return("2023-03-20T12:00:00Z").Maybe()

	idGen.On("NewID").Return("doc-1").Once()

	songRepo.On("CreateSongWithDocuments", mock.Anything, documents).Return("song-123", nil)
	docRepo.On("CreateDocument", mock.Anything).Return("doc-1", nil)

	songID, err := service.CreateSongWithDocuments(song, documents)

	assert.NoError(t, err)
	assert.Equal(t, "song-123", songID)

	songRepo.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
}

func TestCreateSongWithDocuments_Error(t *testing.T) {
	service, songRepo, _, idGen, timeProv := setupSongServiceTest()

	song := models.Song{
		Title:  "Yesterday",
		Author: "The Beatles",
		Genres: []string{"Rock"},
	}

	idGen.On("NewID").Return("song-123")
	timeProv.On("Now").Return("2023-03-20T12:00:00Z").Maybe()

	songRepo.On("CreateSongWithDocuments", mock.Anything, mock.Anything).Return("", errors.New("database error"))

	songID, err := service.CreateSongWithDocuments(song, nil)

	assert.Error(t, err)
	assert.Empty(t, songID)
	assert.Equal(t, "database error", err.Error())
	songRepo.AssertExpectations(t)
}

func TestCreateSongWithDocuments_DocumentCreationError(t *testing.T) {
	service, songRepo, docRepo, idGen, timeProv := setupSongServiceTest()

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

	idGen.On("NewID").Return("song-123").Once()
	timeProv.On("Now").Return("2023-03-20T12:00:00Z").Maybe()

	idGen.On("NewID").Return("doc-1").Once()

	songRepo.On("CreateSongWithDocuments", mock.Anything, documents).Return("song-123", nil)
	docRepo.On("CreateDocument", mock.Anything).Return("", errors.New("failed to create document"))

	songID, err := service.CreateSongWithDocuments(song, documents)

	assert.Error(t, err)
	assert.Empty(t, songID)
	assert.Equal(t, "failed to create document", err.Error())

	songRepo.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
}

func TestUpdateSong_Success(t *testing.T) {
	service, songRepo, _, _, timeProv := setupSongServiceTest()

	updates := map[string]interface{}{
		"title": "Let It Be",
	}

	timeProv.On("Now").Return("2023-03-20T12:00:00Z").Maybe()
	songRepo.On("UpdateSong", "1", mock.Anything).Return(nil)

	err := service.UpdateSong("1", updates)

	assert.NoError(t, err)
	songRepo.AssertExpectations(t)
	timeProv.AssertExpectations(t)
}

func TestUpdateSong_Error(t *testing.T) {
	service, songRepo, _, _, timeProv := setupSongServiceTest()

	updates := map[string]interface{}{
		"title": "Let It Be",
	}

	timeProv.On("Now").Return("2023-03-20T12:00:00Z").Maybe()
	songRepo.On("UpdateSong", "999", mock.Anything).Return(errors.New("song not found"))

	err := service.UpdateSong("999", updates)

	assert.Error(t, err)
	assert.Equal(t, "song not found", err.Error())
	songRepo.AssertExpectations(t)
	timeProv.AssertExpectations(t)
}

func TestDeleteSongWithDocuments_Success(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	songRepo.On("DeleteSongWithDocuments", "1").Return(nil)

	err := service.DeleteSongWithDocuments("1")

	assert.NoError(t, err)
	songRepo.AssertExpectations(t)
}

func TestDeleteSongWithDocuments_Error(t *testing.T) {
	service, songRepo, _, _, _ := setupSongServiceTest()

	songRepo.On("DeleteSongWithDocuments", "999").Return(errors.New("song not found"))

	err := service.DeleteSongWithDocuments("999")

	assert.Error(t, err)
	assert.Equal(t, "song not found", err.Error())
	songRepo.AssertExpectations(t)
}
