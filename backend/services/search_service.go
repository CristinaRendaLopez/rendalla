package services

import (
	"fmt"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

// Ensure SearchService implements SearchServiceInterface.
var _ SearchServiceInterface = (*SearchService)(nil)

// SearchService provides business-level search functionality
// for songs and documents with optional filters and sorting.
type SearchService struct {
	repo repository.SearchRepository
}

// NewSearchService returns a new instance of SearchService.
func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// ListSongs returns a filtered and sorted list of songs with pagination support.
// It validates sorting parameters before forwarding the request to the repository.
func (s *SearchService) ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	sortField, sortOrder = applySortingDefaults(sortField, sortOrder)

	songs, next, err := s.repo.ListSongs(title, sortField, sortOrder, limit, nextToken)
	if err != nil {
		return nil, nil, fmt.Errorf("listing songs: %w", err)
	}

	return songs, next, nil
}

// ListDocuments returns a filtered and sorted list of documents with pagination support.
// It validates sorting parameters before forwarding the request to the repository.
func (s *SearchService) ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	sortField, sortOrder = applySortingDefaults(sortField, sortOrder)

	documents, next, err := s.repo.ListDocuments(title, instrument, docType, sortField, sortOrder, limit, nextToken)
	if err != nil {
		return nil, nil, fmt.Errorf("listing documents: %w", err)
	}

	return documents, next, nil
}

// applySortingDefaults normalizes invalid or empty sortField and sortOrder values.
func applySortingDefaults(sortField, sortOrder string) (string, string) {
	if sortField != "title" && sortField != "created_at" {
		sortField = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	return sortField, sortOrder
}
