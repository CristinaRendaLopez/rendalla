package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockSongService struct {
	mock.Mock
}

var _ services.SongServiceInterface = (*MockSongService)(nil)

func (m *MockSongService) GetAllSongs() ([]dto.SongResponseItem, error) {
	args := m.Called()
	return args.Get(0).([]dto.SongResponseItem), args.Error(1)
}

func (m *MockSongService) GetSongByID(id string) (dto.SongResponseItem, error) {
	args := m.Called(id)
	return args.Get(0).(dto.SongResponseItem), args.Error(1)
}

func (m *MockSongService) CreateSongWithDocuments(req dto.CreateSongRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockSongService) UpdateSong(id string, updates dto.UpdateSongRequest) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockSongService) DeleteSongWithDocuments(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
