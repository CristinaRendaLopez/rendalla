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

func TestListSongs(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		sortField    string
		sortOrder    string
		limit        int
		nextToken    map[string]string
		mockSongs    []models.Song
		mockNext     map[string]string
		mockError    error
		expectError  bool
		expectedSize int
	}{
		{
			name:         "no filters, default sort",
			title:        "",
			sortField:    "",
			sortOrder:    "",
			limit:        10,
			nextToken:    nil,
			mockSongs:    []models.Song{SongBohemianRhapsody},
			mockNext:     nil,
			mockError:    nil,
			expectedSize: 1,
		},
		{
			name:         "filter by title",
			title:        "radio",
			sortField:    "",
			sortOrder:    "",
			limit:        10,
			nextToken:    nil,
			mockSongs:    []models.Song{SongRadioGaGa},
			mockNext:     nil,
			mockError:    nil,
			expectedSize: 1,
		},
		{
			name:         "sort by title asc",
			title:        "",
			sortField:    "title",
			sortOrder:    "asc",
			limit:        10,
			nextToken:    nil,
			mockSongs:    []models.Song{SongBohemianRhapsody, SongRadioGaGa},
			mockNext:     nil,
			mockError:    nil,
			expectedSize: 2,
		},
		{
			name:         "with next token",
			title:        "",
			sortField:    "",
			sortOrder:    "",
			limit:        10,
			nextToken:    map[string]string{"last_id": "2"},
			mockSongs:    []models.Song{SongRadioGaGa},
			mockNext:     ReturnedNextToken,
			mockError:    nil,
			expectedSize: 1,
		},
		{
			name:        "repository error",
			title:       "queen",
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockSearchRepository)
			service := services.NewSearchService(repo)

			repo.On("ListSongs", tt.title, mock.Anything, mock.Anything, tt.limit, tt.nextToken).
				Return(tt.mockSongs, tt.mockNext, tt.mockError)

			songs, next, err := service.ListSongs(tt.title, tt.sortField, tt.sortOrder, tt.limit, tt.nextToken)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, songs, tt.expectedSize)
				assert.Equal(t, tt.mockNext, next)
			}
		})
	}
}

func TestListDocuments(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		instrument   string
		docType      string
		sortField    string
		sortOrder    string
		limit        int
		nextToken    map[string]string
		mockDocs     []models.Document
		mockNext     map[string]string
		mockError    error
		expectError  bool
		expectedSize int
	}{
		{
			name:         "no filters, default sort",
			limit:        10,
			mockDocs:     []models.Document{DocumentPianoScore},
			mockNext:     nil,
			expectedSize: 1,
		},
		{
			name:         "filter by instrument",
			instrument:   "guitar",
			mockDocs:     []models.Document{DocumentGuitarTab},
			expectedSize: 1,
		},
		{
			name:         "combined filters and sort",
			title:        "love",
			instrument:   "violin",
			docType:      "sheet_music",
			sortField:    "title",
			sortOrder:    "asc",
			mockDocs:     []models.Document{DocumentPianoScore, DocumentGuitarTab},
			expectedSize: 2,
		},
		{
			name:         "with next token",
			nextToken:    map[string]string{"last_id": "d1"},
			mockDocs:     []models.Document{DocumentGuitarTab},
			mockNext:     ReturnedNextToken,
			expectedSize: 1,
		},
		{
			name:        "repository error",
			title:       "queen",
			mockError:   errors.ErrInternalServer,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.MockSearchRepository)
			service := services.NewSearchService(repo)

			repo.On("ListDocuments", tt.title, tt.instrument, tt.docType, mock.Anything, mock.Anything, tt.limit, tt.nextToken).
				Return(tt.mockDocs, tt.mockNext, tt.mockError)

			docs, next, err := service.ListDocuments(
				tt.title,
				tt.instrument,
				tt.docType,
				tt.sortField,
				tt.sortOrder,
				tt.limit,
				tt.nextToken,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, docs, tt.expectedSize)
				assert.Equal(t, tt.mockNext, next)
			}
		})
	}
}
