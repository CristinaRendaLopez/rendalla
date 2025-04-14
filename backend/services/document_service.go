package services

import (
	"fmt"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
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

// Ensure DocumentService implements DocumentServiceInterface.
var _ DocumentServiceInterface = (*DocumentService)(nil)

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
func (s *DocumentService) CreateDocument(req dto.CreateDocumentRequest) (string, error) {
	document := dto.ToDocumentModel(req)

	song, err := s.songRepo.GetSongByID(document.SongID)
	if err != nil {
		return "", fmt.Errorf("retrieving song for document creation (song_id=%s): %w", document.SongID, err)
	}

	document.TitleNormalized = utils.Normalize(song.Title)
	document.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	if err := s.repo.CreateDocument(document); err != nil {
		return "", fmt.Errorf("creating document %s: %w", document.ID, err)
	}

	return document.ID, nil
}

// GetDocumentsBySongID returns all documents associated with the specified song ID.
// Returns:
//   - ([]dto.DocumentResponseItem, nil) on success
//   - (nil, error) if the retrieval fails
func (s *DocumentService) GetDocumentsBySongID(songID string) ([]dto.DocumentResponseItem, error) {
	documents, err := s.repo.GetDocumentsBySongID(songID)
	if err != nil {
		return nil, fmt.Errorf("retrieving documents for song %s: %w", songID, err)
	}

	return dto.ToDocumentResponseList(documents), nil
}

// GetDocumentByID retrieves a document by its song ID and document ID.
// Returns:
//   - (dto.DocumentResponseItem, nil) if found
//   - (nil, errors.ErrNotFound) if not found
//   - (nil, error) if the retrieval fails
func (s *DocumentService) GetDocumentByID(songID string, docID string) (dto.DocumentResponseItem, error) {
	doc, err := s.repo.GetDocumentByID(songID, docID)
	if err != nil {
		return dto.DocumentResponseItem{}, fmt.Errorf("retrieving document %s for song %s: %w", docID, songID, err)
	}
	return dto.ToDocumentResponseItem(*doc), nil
}

// UpdateDocument applies updates to a document and refreshes the title_normalized and updated_at fields.
// If title_normalized is not explicitly provided, it is recalculated from the song's title.
// Returns:
//   - nil on success
//   - error if the update fails or the song does not exist
func (s *DocumentService) UpdateDocument(songID, docID string, updates dto.UpdateDocumentRequest) error {
	updateMap := make(map[string]interface{})

	if updates.Type != "" {
		updateMap["type"] = updates.Type
	}
	if len(updates.Instrument) > 0 {
		updateMap["instrument"] = updates.Instrument
	}
	if updates.PDFURL != "" {
		updateMap["pdf_url"] = updates.PDFURL
	}
	if updates.AudioURL != "" {
		updateMap["audio_url"] = updates.AudioURL
	}

	// Calculamos title_normalized si no viene explícito
	if _, ok := updateMap["title_normalized"]; !ok {
		song, err := s.songRepo.GetSongByID(songID)
		if err != nil {
			return fmt.Errorf("retrieving song for update of document %s: %w", docID, err)
		}
		updateMap["title_normalized"] = utils.Normalize(song.Title)
	}

	updateMap["updated_at"] = s.timeProvider.Now()

	if err := s.repo.UpdateDocument(songID, docID, updateMap); err != nil {
		return fmt.Errorf("updating document %s for song %s: %w", docID, songID, err)
	}

	return nil
}

// DeleteDocument deletes a document identified by song ID and document ID.
// Returns:
//   - nil on success
//   - error if the deletion fails
func (s *DocumentService) DeleteDocument(songID string, docID string) error {
	if err := s.repo.DeleteDocument(songID, docID); err != nil {
		return fmt.Errorf("deleting document %s for song %s: %w", docID, songID, err)
	}
	return nil
}
