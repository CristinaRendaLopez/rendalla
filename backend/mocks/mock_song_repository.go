package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockSongRepository struct {
	mock.Mock
}

func (m *MockSongRepository) CreateSongWithDocuments(song models.Song, documents []models.Document) error {
	args := m.Called(song, documents)
	return args.Error(0)
}

func (m *MockSongRepository) GetAllSongs() ([]models.Song, error) {
	args := m.Called()
	return args.Get(0).([]models.Song), args.Error(1)
}

func (m *MockSongRepository) GetSongByID(id string) (*models.Song, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Song), args.Error(1)
}

func (m *MockSongRepository) UpdateSong(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockSongRepository) DeleteSongWithDocuments(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
