package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) GetDocumentByID(songID string, docID string) (*models.Document, error) {
	args := m.Called(songID, docID)
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

func (m *MockDocumentRepository) UpdateDocument(songID string, docID string, updates map[string]interface{}) error {
	args := m.Called(songID, docID, updates)
	return args.Error(0)
}

func (m *MockDocumentRepository) DeleteDocument(songID string, docID string) error {
	args := m.Called(songID, docID)
	return args.Error(0)
}
