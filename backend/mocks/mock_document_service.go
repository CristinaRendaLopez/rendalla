package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockDocumentService struct {
	mock.Mock
}

var _ services.DocumentServiceInterface = (*MockDocumentService)(nil)

func (m *MockDocumentService) CreateDocument(document dto.CreateDocumentRequest) (string, error) {
	args := m.Called(document)
	return args.String(0), args.Error(1)
}

func (m *MockDocumentService) GetDocumentsBySongID(songID string) ([]dto.DocumentResponseItem, error) {
	args := m.Called(songID)
	return args.Get(0).([]dto.DocumentResponseItem), args.Error(1)
}

func (m *MockDocumentService) GetDocumentByID(songID string, docID string) (dto.DocumentResponseItem, error) {
	args := m.Called(songID, docID)
	return args.Get(0).(dto.DocumentResponseItem), args.Error(1)
}

func (m *MockDocumentService) UpdateDocument(songID string, docID string, updates dto.UpdateDocumentRequest) error {
	args := m.Called(songID, docID, updates)
	return args.Error(0)
}

func (m *MockDocumentService) DeleteDocument(songID string, docID string) error {
	args := m.Called(songID, docID)
	return args.Error(0)
}
