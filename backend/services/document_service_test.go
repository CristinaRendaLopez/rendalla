package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDocumentServiceTest() (*services.DocumentService, *mocks.MockDocumentRepository, *mocks.MockSongRepository, *mocks.MockIDGenerator, *mocks.MockTimeProvider) {
	docRepo := new(mocks.MockDocumentRepository)
	songRepo := new(mocks.MockSongRepository)
	idGen := new(mocks.MockIDGenerator)
	timeProv := new(mocks.MockTimeProvider)
	service := services.NewDocumentService(docRepo, songRepo, idGen, timeProv)
	return service, docRepo, songRepo, idGen, timeProv
}

func TestCreateDocument(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateDocumentRequest
		mockSong    *models.Song
		mockSongErr error
		mockDocErr  error
		expectError bool
	}{
		{
			name:        "success",
			request:     ValidCreateDocumentRequest,
			mockSong:    &RelatedSong,
			mockSongErr: nil,
			mockDocErr:  nil,
			expectError: false,
		},
		{
			name:        "song not found",
			request:     ValidCreateDocumentRequest,
			mockSongErr: errors.ErrResourceNotFound,
			expectError: true,
		},
		{
			name:        "document repo error",
			request:     ValidCreateDocumentRequest,
			mockSong:    &RelatedSong,
			mockSongErr: nil,
			mockDocErr:  errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, docRepo, songRepo, idGen, timeProv := setupDocumentServiceTest()

			idGen.On("NewID").Return("doc-1")
			timeProv.On("Now").Return("now")

			if tt.mockSongErr == nil {
				songRepo.On("GetSongByID", tt.request.SongID).Return(tt.mockSong, nil)
			} else {
				songRepo.On("GetSongByID", tt.request.SongID).Return(nil, tt.mockSongErr)
			}

			if tt.mockSongErr == nil {
				docRepo.On("CreateDocument", mock.Anything).Return(tt.mockDocErr)
			}

			docID, err := service.CreateDocument(tt.request)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				docRepo.AssertCalled(t, "CreateDocument", mock.MatchedBy(func(doc models.Document) bool {
					return doc.ID == docID &&
						doc.SongID == tt.request.SongID &&
						doc.Type == tt.request.Type &&
						doc.PDFURL == tt.request.PDFURL
				}))
			}
		})
	}
}

func TestGetDocumentsBySongID(t *testing.T) {
	tests := []struct {
		name           string
		songID         string
		mockDocs       []models.Document
		mockError      error
		expectError    bool
		expectedResult []dto.DocumentResponseItem
	}{
		{
			name:           "documents found",
			songID:         "song-123",
			mockDocs:       []models.Document{MockedDocument},
			mockError:      nil,
			expectError:    false,
			expectedResult: []dto.DocumentResponseItem{DocumentResponse},
		},
		{
			name:           "no documents found",
			songID:         "song-456",
			mockDocs:       []models.Document{},
			mockError:      nil,
			expectError:    false,
			expectedResult: []dto.DocumentResponseItem{},
		},
		{
			name:        "repository error",
			songID:      "song-789",
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, docRepo, _, _, _ := setupDocumentServiceTest()

			docRepo.On("GetDocumentsBySongID", tt.songID).Return(tt.mockDocs, tt.mockError)

			result, err := service.GetDocumentsBySongID(tt.songID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetDocumentByID(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		docID        string
		mockDoc      *models.Document
		mockError    error
		expectError  bool
		expectedItem dto.DocumentResponseItem
	}{
		{
			name:         "document found",
			songID:       "song-123",
			docID:        "doc-1",
			mockDoc:      &MockedDocument,
			mockError:    nil,
			expectError:  false,
			expectedItem: DocumentResponse,
		},
		{
			name:        "document not found",
			songID:      "song-999",
			docID:       "doc-404",
			mockDoc:     nil,
			mockError:   errors.ErrResourceNotFound,
			expectError: true,
		},
		{
			name:        "repository error",
			songID:      "song-123",
			docID:       "doc-err",
			mockDoc:     nil,
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, docRepo, _, _, _ := setupDocumentServiceTest()

			docRepo.On("GetDocumentByID", tt.songID, tt.docID).Return(tt.mockDoc, tt.mockError)

			result, err := service.GetDocumentByID(tt.songID, tt.docID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedItem, result)
			}
		})
	}
}

func TestUpdateDocument(t *testing.T) {
	tests := []struct {
		name          string
		songID        string
		docID         string
		updates       dto.UpdateDocumentRequest
		mockSong      *models.Song
		mockSongErr   error
		mockUpdateErr error
		expectError   bool
	}{
		{
			name:          "successful update",
			songID:        "song-123",
			docID:         "doc-1",
			updates:       ValidUpdateDocumentRequest,
			mockSong:      &RelatedSong,
			mockSongErr:   nil,
			mockUpdateErr: nil,
			expectError:   false,
		},
		{
			name:          "successful update urls",
			songID:        "song-123",
			docID:         "doc-1",
			updates:       ValidUpdateDocumentRequestPDFAndAudio,
			mockSong:      &RelatedSong,
			mockSongErr:   nil,
			mockUpdateErr: nil,
			expectError:   false,
		},
		{
			name:        "song not found",
			songID:      "missing-song",
			docID:       "doc-1",
			updates:     ValidUpdateDocumentRequest,
			mockSong:    nil,
			mockSongErr: errors.ErrResourceNotFound,
			expectError: true,
		},
		{
			name:          "document not found",
			songID:        "song-123",
			docID:         "missing-doc",
			updates:       ValidUpdateDocumentRequest,
			mockSong:      &RelatedSong,
			mockSongErr:   nil,
			mockUpdateErr: nil,
			expectError:   true,
		},
		{
			name:          "repository update fails",
			songID:        "song-123",
			docID:         "doc-1",
			updates:       ValidUpdateDocumentRequest,
			mockSong:      &RelatedSong,
			mockSongErr:   nil,
			mockUpdateErr: errors.ErrInternalServer,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, docRepo, songRepo, _, timeProv := setupDocumentServiceTest()

			timeProv.On("Now").Return("now")

			songRepo.On("GetSongByID", tt.songID).Return(tt.mockSong, tt.mockSongErr)

			if tt.mockSongErr == nil {
				if tt.docID == "missing-doc" {
					docRepo.On("GetDocumentByID", tt.songID, tt.docID).Return(nil, errors.ErrResourceNotFound)
				} else {
					docRepo.On("GetDocumentByID", tt.songID, tt.docID).Return(&MockedDocument, nil)
					docRepo.On("UpdateDocument", tt.songID, tt.docID, mock.Anything).Return(tt.mockUpdateErr)
				}
			}

			err := service.UpdateDocument(tt.songID, tt.docID, tt.updates)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteDocument(t *testing.T) {
	tests := []struct {
		name             string
		songID           string
		docID            string
		mockGetDocErr    error
		mockDeleteDocErr error
		expectError      bool
	}{
		{
			name:             "successful delete",
			songID:           "song-123",
			docID:            "doc-1",
			mockGetDocErr:    nil,
			mockDeleteDocErr: nil,
			expectError:      false,
		},
		{
			name:             "repository error",
			songID:           "song-123",
			docID:            "doc-err",
			mockGetDocErr:    nil,
			mockDeleteDocErr: errors.ErrInternalServer,
			expectError:      true,
		},
		{
			name:             "document not found",
			songID:           "song-123",
			docID:            "missing-doc",
			mockGetDocErr:    errors.ErrResourceNotFound,
			mockDeleteDocErr: nil,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, docRepo, _, _, _ := setupDocumentServiceTest()

			docRepo.On("GetDocumentByID", tt.songID, tt.docID).Return(&MockedDocument, tt.mockGetDocErr)

			if tt.mockGetDocErr == nil {
				docRepo.On("DeleteDocument", tt.songID, tt.docID).Return(tt.mockDeleteDocErr)
			}

			err := service.DeleteDocument(tt.songID, tt.docID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
