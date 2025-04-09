package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

type SearchRepository interface {
	ListSongs(title, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error)
	ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error)
}
