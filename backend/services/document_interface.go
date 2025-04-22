package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
)

// DocumentServiceInterface defines application-level operations for managing musical documents (scores or tablatures) associated with a song.
type DocumentServiceInterface interface {

	// CreateDocument stores a new document.
	// Automatically assigns a UUID and sets creation/update timestamps.
	// Returns:
	//   - the generated document ID on success
	//   - error if creation fails
	CreateDocument(document dto.CreateDocumentRequest) (string, error)

	// GetDocumentsBySongID retrieves all documents linked to the specified song ID.
	// Returns:
	//   - ([]DocumentResponseItem, nil) on success
	//   - (nil, error) if the retrieval fails
	GetDocumentsBySongID(songID string) ([]dto.DocumentResponseItem, error)

	// GetDocumentByID retrieves a single document by song ID and document ID.
	// Returns:
	//   - (*dto.DocumentResponseItem, nil) if found
	//   - (nil, errors.ErrNotFound) if the document does not exist
	//   - (nil, error) for unexpected errors
	GetDocumentByID(songID string, docID string) (dto.DocumentResponseItem, error)

	// UpdateDocument applies partial updates to a document identified by song ID and document ID.
	// Also updates the 'updated_at' timestamp.
	// Returns:
	//   - nil on success
	//   - error if the update operation fails
	UpdateDocument(songID string, docID string, updates dto.UpdateDocumentRequest) error

	// DeleteDocument removes a document by song ID and document ID.
	// Returns:
	//   - nil on success
	//   - error if the deletion fails
	DeleteDocument(songID string, docID string) error
}
