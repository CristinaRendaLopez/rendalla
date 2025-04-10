package services

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// DocumentServiceInterface defines application-level operations for managing musical documents (scores or tablatures) associated with a song.
type DocumentServiceInterface interface {

	// GetDocumentsBySongID retrieves all documents linked to the specified song ID.
	// Returns:
	//   - ([]models.Document, nil) on success
	//   - (nil, error) if the retrieval fails
	GetDocumentsBySongID(songID string) ([]models.Document, error)

	// GetDocumentByID retrieves a single document by song ID and document ID.
	// Returns:
	//   - (*models.Document, nil) if found
	//   - (nil, errors.ErrNotFound) if the document does not exist
	//   - (nil, error) for unexpected errors
	GetDocumentByID(songID string, docID string) (*models.Document, error)

	// CreateDocument stores a new document.
	// Automatically assigns a UUID and sets creation/update timestamps.
	// Returns:
	//   - the generated document ID on success
	//   - error if creation fails
	CreateDocument(document models.Document) (string, error)

	// UpdateDocument applies partial updates to a document identified by song ID and document ID.
	// Also updates the 'updated_at' timestamp.
	// Returns:
	//   - nil on success
	//   - error if the update operation fails
	UpdateDocument(songID string, docID string, updates map[string]interface{}) error

	// DeleteDocument removes a document by song ID and document ID.
	// Returns:
	//   - nil on success
	//   - error if the deletion fails
	DeleteDocument(songID string, docID string) error
}
