package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

type SearchRepository interface {
	ListSongs(title string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error)
	SearchDocumentsByTitle(title string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error)
	FilterDocumentsByInstrument(instrument string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error)
}
