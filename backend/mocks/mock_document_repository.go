package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) GetDocumentByID(id string) (*models.Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *MockDocumentRepository) CreateDocument(doc models.Document) (string, error) {
	args := m.Called(doc)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentRepository) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	args := m.Called(songID)
	return args.Get(0).([]models.Document), args.Error(1)
}

func (m *MockDocumentRepository) UpdateDocument(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteDocument(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
