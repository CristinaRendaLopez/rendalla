package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

// SearchRepository defines methods to search and filter songs and documents with support for pagination.
type SearchRepository interface {

	// ListSongs returns a paginated list of songs filtered by title and sorted by the specified field.
	// Parameters:
	//   - title: optional string to filter by normalized title
	//   - sortField: "title" or "created_at"
	//   - sortOrder: "asc" or "desc"
	//   - limit: number of results to return
	//   - nextToken: token for pagination
	// Returns:
	//   - ([]models.Song, PagingKey, nil) on success
	//   - (nil, nil, error) if the query fails
	ListSongs(title, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error)

	// ListDocuments returns a paginated list of documents filtered by title, instrument, and type.
	// Parameters:
	//   - title: optional string to filter by normalized title
	//   - instrument: optional filter by instrument
	//   - docType: optional filter by document type
	//   - sortField: "title" or "created_at"
	//   - sortOrder: "asc" or "desc"
	//   - limit: number of results to return
	//   - nextToken: token for pagination
	// Returns:
	//   - ([]models.Document, PagingKey, nil) on success
	//   - (nil, nil, error) if the query fails
	ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error)
}
