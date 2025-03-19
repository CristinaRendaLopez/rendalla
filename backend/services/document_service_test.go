package services_test

import (
	"testing"
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
)

func TestGetDocumentByID(t *testing.T) {
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
