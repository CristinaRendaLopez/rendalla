package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
)

type DocumentServiceInterface interface {
	GetDocumentsBySongID(songID string) ([]models.Document, error)
	GetDocumentByID(id string) (*models.Document, error)
	CreateDocument(document models.Document) (string, error)
	UpdateDocument(id string, updates map[string]interface{}) error
	DeleteDocument(id string) error
}

type DocumentService struct {
	repo         repository.DocumentRepository
	idGen        IDGenerator
	timeProvider TimeProvider
}

func NewDocumentService(repo repository.DocumentRepository, idGen IDGenerator, timeProvider TimeProvider) *DocumentService {
	return &DocumentService{
		repo:         repo,
		idGen:        idGen,
		timeProvider: timeProvider,
	}
}

func (s *DocumentService) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	return s.repo.GetDocumentsBySongID(songID)
}

func (s *DocumentService) GetDocumentByID(id string) (*models.Document, error) {
	return s.repo.GetDocumentByID(id)
}

func (s *DocumentService) CreateDocument(document models.Document) (string, error) {
	document.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	return s.repo.CreateDocument(document)
}

func (s *DocumentService) UpdateDocument(id string, updates map[string]interface{}) error {
	updates["updated_at"] = s.timeProvider.Now()
	return s.repo.UpdateDocument(id, updates)
}

func (s *DocumentService) DeleteDocument(id string) error {
	return s.repo.DeleteDocument(id)
}
