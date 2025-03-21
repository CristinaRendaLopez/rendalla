package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

type SearchServiceInterface interface {
	SearchSongsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error)
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

func (s *SearchService) SearchSongsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Song, repository.PagingKey, error) {
	return s.repo.SearchSongsByTitle(title, limit, nextToken)
}

func (s *SearchService) SearchDocumentsByTitle(title string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	return s.repo.SearchDocumentsByTitle(title, limit, nextToken)
}

func (s *SearchService) FilterDocumentsByInstrument(instrument string, limit int, nextToken repository.PagingKey) ([]models.Document, repository.PagingKey, error) {
	return s.repo.FilterDocumentsByInstrument(instrument, limit, nextToken)
}
