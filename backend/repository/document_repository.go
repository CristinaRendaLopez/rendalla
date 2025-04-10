package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// DocumentRepository defines operations for managing musical documents (scores or tablatures) stored with a composite key (song_id, id) in the database.
type DocumentRepository interface {

	// CreateDocument stores a new document.
	// Returns:
	//   - nil on success
	//   - errors.ErrInternalServer if marshalling or persistence fails
	CreateDocument(doc models.Document) error

	// GetDocumentsBySongID returns all documents linked to a given song.
	// Returns:
	//   - ([]models.Document, nil) on success
	//   - (nil, errors.ErrInternalServer) if the query fails
	GetDocumentsBySongID(songID string) ([]models.Document, error)

	// GetDocumentByID retrieves a document by its song ID and document ID.
	// Returns:
	//   - (*models.Document, nil) if found
	//   - (nil, errors.ErrNotFound) if the document does not exist
	//   - (nil, errors.ErrInternalServer) if retrieval fails
	GetDocumentByID(songID string, documentID string) (*models.Document, error)

	// UpdateDocument applies partial updates to a document by its song ID and document ID.
	// Returns:
	//   - nil on success
	//   - errors.ErrInternalServer if the update operation fails
	UpdateDocument(songID string, documentID string, updates map[string]interface{}) error

	// DeleteDocument removes a document by its song ID and document ID.
	// Returns:
	//   - nil on success
	//   - errors.ErrInternalServer if the deletion fails
	DeleteDocument(songID string, documentID string) error
}
