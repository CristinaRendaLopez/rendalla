package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDocumentServiceTest() (*services.DocumentService, *mocks.MockDocumentRepository, *mocks.MockIDGenerator, *mocks.MockTimeProvider) {
	docRepo := new(mocks.MockDocumentRepository)
	idGen := new(mocks.MockIDGenerator)
	timeProv := new(mocks.MockTimeProvider)
	service := services.NewDocumentService(docRepo, idGen, timeProv)
	return service, docRepo, idGen, timeProv
}

func TestGetDocumentByID_Success(t *testing.T) {
	service, docRepo, _, _ := setupDocumentServiceTest()

	expectedDoc := &models.Document{
		ID:         "doc1",
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Guitar"},
		PDFURL:     "https://s3.amazonaws.com/queen/wewillrockyou.pdf",
		CreatedAt:  "2023-03-21T00:00:00Z",
		UpdatedAt:  "2023-03-21T00:00:00Z",
	}

	docRepo.On("GetDocumentByID", "doc1").Return(expectedDoc, nil)

	doc, err := service.GetDocumentByID("doc1")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, expectedDoc, doc)

	docRepo.AssertExpectations(t)
}

func TestGetDocumentByID_NotFound(t *testing.T) {
	service, docRepo, _, _ := setupDocumentServiceTest()

	docRepo.On("GetDocumentByID", "unknown").Return(nil, utils.ErrResourceNotFound)

	doc, err := service.GetDocumentByID("unknown")

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.ErrorIs(t, err, utils.ErrResourceNotFound)

	docRepo.AssertExpectations(t)
}

func TestGetDocumentsBySongID_Success(t *testing.T) {
	service, docRepo, _, _ := setupDocumentServiceTest()

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
	service, docRepo, idGen, timeProv := setupDocumentServiceTest()

	doc := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Piano"},
		PDFURL:     "url",
	}

	idGen.On("NewID").Return("doc123")
	timeProv.On("Now").Return("2023-03-21T00:00:00Z").Maybe()

	docRepo.On("CreateDocument", mock.Anything).Return("doc123", nil)

	docID, err := service.CreateDocument(doc)

	assert.NoError(t, err)
	assert.Equal(t, "doc123", docID)

	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
}

func TestCreateDocument_Error(t *testing.T) {
	service, docRepo, idGen, timeProv := setupDocumentServiceTest()

	doc := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Violin"},
		PDFURL:     "url",
	}

	idGen.On("NewID").Return("doc123")
	timeProv.On("Now").Return("2023-03-21T00:00:00Z").Maybe()

	docRepo.On("CreateDocument", mock.Anything).Return("", utils.ErrOperationNotAllowed)

	docID, err := service.CreateDocument(doc)

	assert.Error(t, err)
	assert.Empty(t, docID)
	assert.ErrorIs(t, err, utils.ErrOperationNotAllowed)

	idGen.AssertExpectations(t)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
}

func TestUpdateDocument_Success(t *testing.T) {
	service, docRepo, _, timeProv := setupDocumentServiceTest()

	updates := map[string]interface{}{
		"type": "tablature",
	}

	timeProv.On("Now").Return("2023-03-21T00:00:00Z")
	docRepo.On("UpdateDocument", "doc1", mock.Anything).Return(nil)

	err := service.UpdateDocument("doc1", updates)

	assert.NoError(t, err)
	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
}

func TestUpdateDocument_Error(t *testing.T) {
	service, docRepo, _, timeProv := setupDocumentServiceTest()

	updates := map[string]interface{}{
		"type": "sheet_music",
	}

	timeProv.On("Now").Return("2023-03-21T00:00:00Z")
	docRepo.On("UpdateDocument", "doc1", mock.Anything).Return(utils.ErrInternalServer)

	err := service.UpdateDocument("doc1", updates)

	assert.Error(t, err)
	assert.ErrorIs(t, err, utils.ErrInternalServer)

	timeProv.AssertExpectations(t)
	docRepo.AssertExpectations(t)
}

func TestDeleteDocument_Success(t *testing.T) {
	service, docRepo, _, _ := setupDocumentServiceTest()

	docRepo.On("DeleteDocument", "doc1").Return(nil)

	err := service.DeleteDocument("doc1")

	assert.NoError(t, err)
	docRepo.AssertExpectations(t)
}

func TestDeleteDocument_Error(t *testing.T) {
	service, docRepo, _, _ := setupDocumentServiceTest()

	docRepo.On("DeleteDocument", "unknown").Return(utils.ErrResourceNotFound)

	err := service.DeleteDocument("unknown")

	assert.Error(t, err)
	assert.ErrorIs(t, err, utils.ErrResourceNotFound)

	docRepo.AssertExpectations(t)
}
