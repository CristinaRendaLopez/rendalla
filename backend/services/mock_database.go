package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

// Song-related methods
func (m *MockDB) GetSongByID(id string) (*models.Song, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Song), args.Error(1)
}

func (m *MockDB) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	args := m.Called(song, documents)
	return args.String(0), args.Error(1)
}

func (m *MockDB) UpdateSong(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockDB) DeleteSongWithDocuments(songID string) error {
	args := m.Called(songID)
	return args.Error(0)
}

// Document-related methods
func (m *MockDB) GetDocumentByID(id string) (*models.Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Document), args.Error(1)
}

func (m *MockDB) CreateDocument(document models.Document) error {
	args := m.Called(document)
	return args.Error(0)
}

func (m *MockDB) UpdateDocument(id string, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockDB) DeleteDocument(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// Search related methods
func (m *MockDB) SearchSongsByTitle(title string, limit int, nextToken interface{}) ([]models.Song, interface{}, error) {
	args := m.Called(title, limit, nextToken)
	return args.Get(0).([]models.Song), args.Get(1), args.Error(2)
}

func (m *MockDB) SearchDocumentsByTitle(title string, limit int, nextToken interface{}) ([]models.Document, interface{}, error) {
	args := m.Called(title, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1), args.Error(2)
}

func (m *MockDB) FilterDocumentsByInstrument(instrument string, limit int, nextToken interface{}) ([]models.Document, interface{}, error) {
	args := m.Called(instrument, limit, nextToken)
	return args.Get(0).([]models.Document), args.Get(1), args.Error(2)
}

// Authentication related methods
func (m *MockDB) AuthenticateUser(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockDB) GetAuthCredentials() (*AuthCredentials, error) {
	args := m.Called()

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*AuthCredentials), args.Error(1)
}
