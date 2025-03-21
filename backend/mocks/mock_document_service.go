package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockDocumentService struct {
	mock.Mock
}

var _ services.DocumentServiceInterface = (*MockDocumentService)(nil)

func (m *MockDocumentService) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	args := m.Called(songID)
	return args.Get(0).([]models.Document), args.Error(1)
}

func (m *MockDocumentService) GetDocumentByID(id string) (*models.Document, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *MockDocumentService) CreateDocument(document models.Document) (string, error) {
	args := m.Called(document)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentService) UpdateDocument(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteDocument(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
