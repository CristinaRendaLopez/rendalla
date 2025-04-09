package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// DocumentRepository defines operations for managing musical documents (scores or tablatures).
type DocumentRepository interface {

	// CreateDocument stores a new document.
	CreateDocument(doc models.Document) error

	// GetDocumentsBySongID returns all documents linked to a given song.
	GetDocumentsBySongID(songID string) ([]models.Document, error)

	// GetDocumentByID retrieves a document by its song ID and document ID.
	GetDocumentByID(songID string, documentID string) (*models.Document, error)

	// UpdateDocument applies partial updates to a document by its song ID and document ID.
	UpdateDocument(songID string, documentID string, updates map[string]interface{}) error

	// DeleteDocument removes a document by its song ID and document ID.
	DeleteDocument(songID string, documentID string) error
}
