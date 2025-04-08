package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

type DocumentRepository interface {
	GetDocumentByID(songID string, documentID string) (*models.Document, error)
	CreateDocument(doc models.Document) (string, error)
	UpdateDocument(songID string, documentID string, updates map[string]interface{}) error
	DeleteDocument(songID string, documentID string) error
	GetDocumentsBySongID(songID string) ([]models.Document, error)
}
