package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDocumentByID_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	expectedDoc := &models.Document{
		ID:         "doc1",
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Guitar"},
		PDFURL:     "https://s3.amazonaws.com/queen/wewillrockyou.pdf",
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	mockDocRepo.On("GetDocumentByID", "doc1").Return(expectedDoc, nil)

	doc, err := service.GetDocumentByID("doc1")

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, expectedDoc, doc)

	mockDocRepo.AssertExpectations(t)
}

func TestGetDocumentByID_NotFound(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	mockDocRepo.On("GetDocumentByID", "unknown").Return(nil, errors.New("document not found"))

	doc, err := service.GetDocumentByID("unknown")

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Equal(t, "document not found", err.Error())

	mockDocRepo.AssertExpectations(t)
}

func TestGetDocumentsBySongID_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	documents := []models.Document{
		{ID: "doc1", SongID: "song1", Type: "sheet_music", Instrument: []string{"Guitar"}, PDFURL: "https://s3.amazonaws.com/queen/wewillrockyou.pdf"},
		{ID: "doc2", SongID: "song1", Type: "sheet_music", Instrument: []string{"Piano"}, PDFURL: "https://s3.amazonaws.com/queen/bohemian.pdf"},
	}

	mockDocRepo.On("GetDocumentsBySongID", "song1").Return(documents, nil)

	result, err := service.GetDocumentsBySongID("song1")

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockDocRepo.AssertExpectations(t)
}

func TestUpdateDocument_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	updates := map[string]interface{}{
		"type": "tablature",
	}

	mockDocRepo.On("UpdateDocument", "doc1", mock.Anything).Return(nil)

	err := service.UpdateDocument("doc1", updates)

	assert.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
}

func TestCreateDocument_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	document := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Piano"},
		PDFURL:     "https://s3.amazonaws.com/queen/bohemian.pdf",
	}

	mockDocRepo.On("CreateDocument", mock.Anything).Return("doc123", nil)

	docID, err := service.CreateDocument(document)

	assert.NoError(t, err)
	assert.Equal(t, "doc123", docID)

	mockDocRepo.AssertExpectations(t)
}

func TestCreateDocument_Error(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	document := models.Document{
		SongID:     "song1",
		Type:       "sheet_music",
		Instrument: []string{"Violin"},
		PDFURL:     "https://s3.amazonaws.com/music/violin.pdf",
	}

	mockDocRepo.On("CreateDocument", mock.Anything).Return("", errors.New("failed to create document"))

	docID, err := service.CreateDocument(document)

	assert.Error(t, err)
	assert.Empty(t, docID)
	assert.Equal(t, "failed to create document", err.Error())

	mockDocRepo.AssertExpectations(t)
}

func TestDeleteDocument_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	mockDocRepo.On("DeleteDocument", "doc1").Return(nil)

	err := service.DeleteDocument("doc1")

	assert.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
}

func TestDeleteDocument_Error(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	service := services.NewDocumentService(mockDocRepo)

	mockDocRepo.On("DeleteDocument", "unknown").Return(errors.New("document not found"))

	err := service.DeleteDocument("unknown")

	assert.Error(t, err)
	assert.Equal(t, "document not found", err.Error())

	mockDocRepo.AssertExpectations(t)
}
