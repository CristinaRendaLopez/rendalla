package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

// SearchServiceInterface defines application-level operations for searching and filtering songs and documents.
// All methods support pagination and optional sorting.
type SearchServiceInterface interface {

	// ListSongs returns a paginated list of songs filtered by title and optionally sorted.
	// Parameters:
	//   - title: optional search term (normalized internally)
	//   - sortField: "title" or "created_at"
	//   - sortOrder: "asc" or "desc"
	//   - limit: max number of results to return
	//   - nextToken: pagination token to resume from last result
	// Returns:
	//   - a list of songs
	//   - a token for the next page (or nil)
	//   - error if the query fails
	ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error)

	// ListDocuments returns a paginated list of documents filtered by title, instrument, and type.
	// Parameters:
	//   - title: optional title filter (normalized)
	//   - instrument: optional instrument filter
	//   - docType: optional document type filter ("score", "tablature", etc.)
	//   - sortField: "title" or "created_at"
	//   - sortOrder: "asc" or "desc"
	//   - limit: max number of results to return
	//   - nextToken: pagination token to resume from last result
	// Returns:
	//   - a list of documents
	//   - a token for the next page (or nil)
	//   - error if the query fails
	ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error)
}
