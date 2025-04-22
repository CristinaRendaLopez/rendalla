package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockSearchService struct {
	mock.Mock
}

var _ services.SearchServiceInterface = (*MockSearchService)(nil)

func (m *MockSearchService) ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {

	args := m.Called(title, sortField, sortOrder, limit, nextToken)
	return args.Get(0).([]models.Song), args.Get(1), args.Error(2)
}

func (m *MockSearchService) ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	args := m.Called(title, instrument, docType, sortField, sortOrder, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1), args.Error(2)
}
