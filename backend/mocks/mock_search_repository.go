package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) ListSongs(title string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	args := m.Called(title, limit, nextToken)
	return args.Get(0).([]models.Song), args.Get(1).(repository.PagingKey), args.Error(2)
}

func (m *MockSearchRepository) SearchDocumentsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	args := m.Called(title, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1).(repository.PagingKey), args.Error(2)
}

func (m *MockSearchRepository) FilterDocumentsByInstrument(instrument string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	args := m.Called(instrument, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1).(repository.PagingKey), args.Error(2)
}
