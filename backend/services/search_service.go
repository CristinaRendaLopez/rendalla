package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

// SearchService provides high-level operations for searching songs and documents.
// It delegates the actual querying to a SearchRepository, and enforces defaults for sorting.
type SearchService struct {
	repo repository.SearchRepository
}

// Ensure SearchService implements SearchServiceInterface.
var _ SearchServiceInterface = (*SearchService)(nil)

// NewSearchService returns a new instance of SearchService.
func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

// ListSongs delegates the search to the repository, applying default sorting if needed.
// It filters by normalized title and supports sorting and pagination.
// Returns:
//   - a slice of songs
//   - the next pagination token, or nil if no more results
//   - error if the operation fails
func (s *SearchService) ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	if sortField != "title" && sortField != "created_at" {
		sortField = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	return s.repo.ListSongs(title, sortField, sortOrder, limit, nextToken)
}

// ListDocuments delegates the search to the repository, applying default sorting if needed.
// Supports filters by title, instrument, and document type, as well as pagination and sorting.
// Returns:
//   - a slice of documents
//   - the next pagination token, or nil if no more results
//   - error if the operation fails
func (s *SearchService) ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	if sortField != "title" && sortField != "created_at" {
		sortField = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	return s.repo.ListDocuments(title, instrument, docType, sortField, sortOrder, limit, nextToken)
}
