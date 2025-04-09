package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

// SearchRepository defines methods to search and filter songs and documents with support for pagination.
type SearchRepository interface {

	// ListSongs returns a paginated list of songs filtered by title, with optional sorting.
	ListSongs(title, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error)

	// ListDocuments returns a paginated list of documents filtered by title, instrument, and type, with optional sorting.
	ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error)
}
