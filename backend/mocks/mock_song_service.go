package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockSongService struct {
	mock.Mock
}

var _ services.SongServiceInterface = (*MockSongService)(nil)

func (m *MockSongService) GetAllSongs() ([]models.Song, error) {
	args := m.Called()
	return args.Get(0).([]models.Song), args.Error(1)
}

func (m *MockSongService) GetSongByID(id string) (*models.Song, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Song), args.Error(1)
}

func (m *MockSongService) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	args := m.Called(song, documents)
	return args.String(0), args.Error(1)
}

func (m *MockSongService) UpdateSong(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockSongService) DeleteSongWithDocuments(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
