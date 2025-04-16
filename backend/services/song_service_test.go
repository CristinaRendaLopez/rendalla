package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
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

func TestCreateSongWithDocuments(t *testing.T) {
	tests := []struct {
		name          string
		request       dto.CreateSongRequest
		mockSongError error
		expectError   bool
	}{
		{
			name:          "success",
			request:       ValidCreateSongRequest,
			mockSongError: nil,
			expectError:   false,
		},
		{
			name:        "validation fails",
			request:     InvalidCreateSongRequest,
			expectError: true,
		},
		{
			name:          "repo error",
			request:       ValidCreateSongRequest,
			mockSongError: errors.ErrInternalServer,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, songRepo, _, idGen, timeProvider := setupSongServiceTest()
			idGen.On("NewID").Return("id").Maybe()
			timeProvider.On("Now").Return("now").Maybe()

			if tt.mockSongError == nil && !tt.expectError {
				songRepo.On("CreateSongWithDocuments", mock.Anything, mock.Anything).
					Return(tt.mockSongError)
			} else if tt.mockSongError != nil {
				songRepo.On("CreateSongWithDocuments", mock.Anything, mock.Anything).
					Return(tt.mockSongError)
			}

			_, err := service.CreateSongWithDocuments(tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllSongs(t *testing.T) {
	tests := []struct {
		name         string
		mockSongs    []models.Song
		mockError    error
		expectError  bool
		expectedSize int
	}{
		{
			name:         "returns songs",
			mockSongs:    []models.Song{*MockedSong},
			mockError:    nil,
			expectError:  false,
			expectedSize: 1,
		},
		{
			name:         "empty list",
			mockSongs:    []models.Song{},
			mockError:    nil,
			expectError:  false,
			expectedSize: 0,
		},
		{
			name:         "repository error",
			mockError:    errors.ErrInternalServer,
			expectError:  true,
			expectedSize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, songRepo, _, _, _ := setupSongServiceTest()

			songRepo.On("GetAllSongs").Return(tt.mockSongs, tt.mockError)

			songs, err := service.GetAllSongs()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, songs, tt.expectedSize)
			}
		})
	}
}

func TestGetSongByID(t *testing.T) {
	tests := []struct {
		name        string
		songID      string
		mockSong    *models.Song
		mockError   error
		expectError bool
		expectedID  string
	}{
		{
			name:        "found",
			songID:      "1",
			mockSong:    MockedSong,
			expectError: false,
			expectedID:  "1",
		},
		{
			name:        "not found",
			songID:      "2",
			mockError:   errors.ErrResourceNotFound,
			expectError: true,
		},
		{
			name:        "repository error",
			songID:      "3",
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, songRepo, _, _, _ := setupSongServiceTest()

			songRepo.On("GetSongByID", tt.songID).Return(tt.mockSong, tt.mockError)

			song, err := service.GetSongByID(tt.songID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, song.ID)
			}
		})
	}
}

func TestUpdateSong(t *testing.T) {
	tests := []struct {
		name            string
		songID          string
		updates         dto.UpdateSongRequest
		mockGetError    error
		mockUpdateError error
		expectError     bool
	}{
		{
			name:            "successful update",
			songID:          "1",
			updates:         ValidUpdateSongRequest,
			mockGetError:    nil,
			mockUpdateError: nil,
			expectError:     false,
		},
		{
			name:            "successful author update",
			songID:          "1",
			updates:         ValidAuthorUpdateRequest,
			mockGetError:    nil,
			mockUpdateError: nil,
			expectError:     false,
		},
		{
			name:         "validation fails",
			songID:       "1",
			updates:      InvalidUpdateSongRequest,
			mockGetError: nil,
			expectError:  true,
		},
		{
			name:         "song not found",
			songID:       "2",
			updates:      SongNotFoundInvalidUpdateSongRequest,
			mockGetError: errors.ErrResourceNotFound,
			expectError:  true,
		},
		{
			name:            "repository update fails",
			songID:          "3",
			updates:         ValidUpdateSongRequest,
			mockGetError:    nil,
			mockUpdateError: errors.ErrInternalServer,
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, songRepo, _, _, timeProvider := setupSongServiceTest()

			timeProvider.On("Now").Return("mocked-time").Maybe()

			if tt.mockGetError == nil {
				songRepo.On("GetSongByID", tt.songID).Return(&models.Song{ID: tt.songID}, nil)
			} else {
				songRepo.On("GetSongByID", tt.songID).Return(nil, tt.mockGetError)
			}

			if tt.mockUpdateError != nil {
				songRepo.On("UpdateSong", tt.songID, mock.Anything).Return(tt.mockUpdateError)
			} else if tt.mockGetError == nil {
				songRepo.On("UpdateSong", tt.songID, mock.Anything).Return(nil)
			}

			err := service.UpdateSong(tt.songID, tt.updates)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteSongWithDocuments(t *testing.T) {
	tests := []struct {
		name        string
		songID      string
		mockError   error
		expectError bool
	}{
		{
			name:        "success",
			songID:      "1",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repository error",
			songID:      "2",
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, songRepo, _, _, _ := setupSongServiceTest()

			songRepo.On("DeleteSongWithDocuments", tt.songID).Return(tt.mockError)

			err := service.DeleteSongWithDocuments(tt.songID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
