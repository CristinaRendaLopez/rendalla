package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

// DocumentService provides application-level operations for managing musical documents
// such as scores and tablatures, in relation to songs.
type DocumentService struct {
	repo         repository.DocumentRepository
	songRepo     repository.SongRepository
	idGen        utils.IDGenerator
	timeProvider utils.TimeProvider
}

// NewDocumentService returns a new instance of DocumentService.
func NewDocumentService(
	repo repository.DocumentRepository,
	songRepo repository.SongRepository,
	idGen utils.IDGenerator,
	timeProvider utils.TimeProvider,
) *DocumentService {
	return &DocumentService{
		repo:         repo,
		songRepo:     songRepo,
		idGen:        idGen,
		timeProvider: timeProvider,
	}
}

// CreateDocument creates and stores a new document linked to a song.
// It normalizes the song's title, assigns a UUID, and sets timestamps.
// Returns:
//   - the generated document ID on success
//   - error if the song is not found or document creation fails
func (s *DocumentService) CreateDocument(document models.Document) (string, error) {
	song, err := s.songRepo.GetSongByID(document.SongID)
	if err != nil {
		return "", err
	}

	document.TitleNormalized = utils.Normalize(song.Title)
	document.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	return document.ID, s.repo.CreateDocument(document)
}

// GetDocumentsBySongID returns all documents associated with the specified song ID.
// Returns:
//   - ([]models.Document, nil) on success
//   - (nil, error) if the retrieval fails
func (s *DocumentService) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	return s.repo.GetDocumentsBySongID(songID)
}

// GetDocumentByID retrieves a document by its song ID and document ID.
// Returns:
//   - (*models.Document, nil) if found
//   - (nil, utils.ErrNotFound) if not found
//   - (nil, error) if the retrieval fails
func (s *DocumentService) GetDocumentByID(songID string, docID string) (*models.Document, error) {
	return s.repo.GetDocumentByID(songID, docID)
}

// UpdateDocument applies updates to a document and refreshes the title_normalized and updated_at fields.
// If title_normalized is not explicitly provided, it is recalculated from the song's title.
// Returns:
//   - nil on success
//   - error if the update fails or the song does not exist
func (s *DocumentService) UpdateDocument(songID, docID string, updates map[string]interface{}) error {
	if _, ok := updates["title_normalized"]; !ok {
		song, err := s.songRepo.GetSongByID(songID)
		if err != nil {
			return err
		}
		updates["title_normalized"] = utils.Normalize(song.Title)
	}
	updates["updated_at"] = s.timeProvider.Now()
	return s.repo.UpdateDocument(songID, docID, updates)
}

// DeleteDocument deletes a document identified by song ID and document ID.
// Returns:
//   - nil on success
//   - error if the deletion fails
func (s *DocumentService) DeleteDocument(songID string, docID string) error {
	return s.repo.DeleteDocument(songID, docID)
}
