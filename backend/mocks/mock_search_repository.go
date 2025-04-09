package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	args := m.Called(title, sortField, sortOrder, limit, nextToken)
	return args.Get(0).([]models.Song), args.Get(1).(repository.PagingKey), args.Error(2)
}

func (m *MockSearchRepository) ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	args := m.Called(title, instrument, docType, sortField, sortOrder, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1).(repository.PagingKey), args.Error(2)
}
