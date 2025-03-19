package services

import (
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/google/uuid"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	return s.repo.GetDocumentsBySongID(songID)
}

func (s *DocumentService) GetDocumentByID(id string) (*models.Document, error) {
	return s.repo.GetDocumentByID(id)
}

func (s *DocumentService) CreateDocument(document models.Document) (string, error) {
	document.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	document.CreatedAt, document.UpdatedAt = now, now

	return s.repo.CreateDocument(document)
}

func (s *DocumentService) UpdateDocument(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	return s.repo.UpdateDocument(id, updates)
}

func (s *DocumentService) DeleteDocument(id string) error {
	return s.repo.DeleteDocument(id)
}
