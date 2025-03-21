package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

type DocumentRepository interface {
	GetDocumentByID(id string) (*models.Document, error)
	CreateDocument(doc models.Document) (string, error)
	UpdateDocument(id string, updates map[string]interface{}) error
	DeleteDocument(id string) error
	GetDocumentsBySongID(songID string) ([]models.Document, error)
}
