package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

type SearchServiceInterface interface {
	ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error)
	SearchDocumentsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error)
	FilterDocumentsByInstrument(instrument string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error)
}

type SearchService struct {
	repo repository.SearchRepository
}

var _ SearchServiceInterface = (*SearchService)(nil)

func NewSearchService(repo repository.SearchRepository) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) ListSongs(title, sortField, sortOrder string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	if sortField != "title" && sortField != "created_at" {
		sortField = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	return s.repo.ListSongs(title, sortField, sortOrder, limit, nextToken)
}

func (s *SearchService) SearchDocumentsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	return s.repo.SearchDocumentsByTitle(title, limit, nextToken)
}

func (s *SearchService) FilterDocumentsByInstrument(instrument string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	return s.repo.FilterDocumentsByInstrument(instrument, limit, nextToken)
}
