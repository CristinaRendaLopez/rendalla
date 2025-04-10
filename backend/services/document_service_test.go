package services_test

import (
	"testing"

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

func TestGetDocumentByID_Success(t *testing.T) {
	service, docRepo, _, _, _ := setupDocumentServiceTest()

	expectedDoc := &models.Document{
		ID:         "doc1",
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Guitar"},
		PDFURL:     "https://s3.amazonaws.com/queen/wewillrockyou.pdf",
		CreatedAt:  "2023-03-21T00:00:00Z",
		UpdatedAt:  "2023-03-21T00:00:00Z",
	}

	docRepo.On("GetDocumentByID", "song1", "doc1").Return(expectedDoc, nil)

	doc, err := service.GetDocumentByID("song1", "doc1")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, expectedDoc, doc)

	docRepo.AssertExpectations(t)
}

func TestGetDocumentByID_NotFound(t *testing.T) {
	service, docRepo, _, _, _ := setupDocumentServiceTest()

	docRepo.On("GetDocumentByID", "song1", "unknown").Return(nil, errors.ErrResourceNotFound)

	doc, err := service.GetDocumentByID("song1", "unknown")

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	docRepo.AssertExpectations(t)
}

func TestGetDocumentsBySongID_Success(t *testing.T) {
	service, docRepo, _, _, _ := setupDocumentServiceTest()

	docs := []models.Document{
		{ID: "doc1", SongID: "song1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "url1"},
		{ID: "doc2", SongID: "song1", Type: "sheet_music", Instrument: []string{"Piano"}, PDFURL: "url2"},
	}

	docRepo.On("GetDocumentsBySongID", "song1").Return(docs, nil)

	result, err := service.GetDocumentsBySongID("song1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	docRepo.AssertExpectations(t)
}

func TestCreateDocument_Success(t *testing.T) {
	service, docRepo, songRepo, idGen, timeProv := setupDocumentServiceTest()

	doc := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Piano"},
		PDFURL:     "url",
	}

	songRepo.On("GetSongByID", "song1").Return(&models.Song{
		ID:    "song1",
		Title: "Bohemian Rhapsody",
	}, nil)

	idGen.On("NewID").Return("doc123")
	timeProv.On("Now").Return("2023-03-21T00:00:00Z").Maybe()

	docRepo.On("CreateDocument", mock.Anything).Return(nil)

	docID, err := service.CreateDocument(doc)

	assert.NoError(t, err)
	assert.Equal(t, "doc123", docID)

	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	songRepo.AssertExpectations(t)
}

func TestCreateDocument_Error(t *testing.T) {
	service, docRepo, songRepo, idGen, timeProv := setupDocumentServiceTest()

	doc := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Violin"},
		PDFURL:     "url",
	}

	songRepo.On("GetSongByID", "song1").Return(&models.Song{
		ID:    "song1",
		Title: "Bohemian Rhapsody",
	}, nil)

	idGen.On("NewID").Return("doc123")
	timeProv.On("Now").Return("2023-03-21T00:00:00Z").Maybe()

	docRepo.On("CreateDocument", mock.Anything).Return(errors.ErrOperationNotAllowed)

	_, err := service.CreateDocument(doc)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrOperationNotAllowed)

	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	songRepo.AssertExpectations(t)
}

func TestUpdateDocument_Success(t *testing.T) {
	service, docRepo, songRepo, _, timeProv := setupDocumentServiceTest()

	updates := map[string]interface{}{
		"type": "tablature",
	}

	songRepo.On("GetSongByID", "song1").Return(&models.Song{
		ID:    "song1",
		Title: "Don't Stop Me Now",
	}, nil)

	timeProv.On("Now").Return("2023-03-21T00:00:00Z")
	docRepo.On("UpdateDocument", "song1", "doc1", mock.Anything).Return(nil)

	err := service.UpdateDocument("song1", "doc1", updates)

	assert.NoError(t, err)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	songRepo.AssertExpectations(t)
}

func TestUpdateDocument_Error(t *testing.T) {
	service, docRepo, songRepo, _, timeProv := setupDocumentServiceTest()

	updates := map[string]interface{}{
		"type": "sheet_music",
	}

	songRepo.On("GetSongByID", "song1").Return(&models.Song{
		ID:    "song1",
		Title: "I Want It All",
	}, nil)

	timeProv.On("Now").Return("2023-03-21T00:00:00Z")
	docRepo.On("UpdateDocument", "song1", "doc1", mock.Anything).Return(errors.ErrInternalServer)

	err := service.UpdateDocument("song1", "doc1", updates)

	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrInternalServer)

	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
	songRepo.AssertExpectations(t)
}

func TestDeleteDocument_Success(t *testing.T) {
	service, docRepo, _, _, _ := setupDocumentServiceTest()

	docRepo.On("DeleteDocument", "song1", "doc1").Return(nil)

	err := service.DeleteDocument("song1", "doc1")

	assert.NoError(t, err)
	docRepo.AssertExpectations(t)
}

func TestDeleteDocument_Error(t *testing.T) {
	service, docRepo, _, _, _ := setupDocumentServiceTest()

	docRepo.On("DeleteDocument", "song1", "unknown").Return(errors.ErrResourceNotFound)

	err := service.DeleteDocument("song1", "unknown")

	assert.Error(t, err)
	assert.ErrorIs(t, err, errors.ErrResourceNotFound)

	docRepo.AssertExpectations(t)
}
