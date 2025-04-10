package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

// SongService provides application-level operations for managing songs and their associated documents.
// It uses repositories for persistence and utility interfaces for time and ID generation.
type SongService struct {
	songRepo     repository.SongRepository
	docRepo      repository.DocumentRepository
	idGen        utils.IDGenerator
	timeProvider utils.TimeProvider
}

// NewSongService returns a new instance of SongService with its required dependencies.
func NewSongService(
	songRepo repository.SongRepository,
	docRepo repository.DocumentRepository,
	idGen utils.IDGenerator,
	timeProvider utils.TimeProvider,
) *SongService {
	return &SongService{
		songRepo:     songRepo,
		docRepo:      docRepo,
		idGen:        idGen,
		timeProvider: timeProvider,
	}
}

// CreateSongWithDocuments creates a new song and all associated documents.
// It generates UUIDs and timestamps, and normalizes the title before saving.
// Returns:
//   - the generated song ID on success
//   - error if the creation fails at any point
func (s *SongService) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	if err := utils.ValidateSongAndDocuments(song, documents); err != nil {
		return "", err
	}

	song.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	song.CreatedAt = now
	song.UpdatedAt = now
	song.TitleNormalized = utils.Normalize(song.Title)

	for i := range documents {
		documents[i].ID = s.idGen.NewID()
		documents[i].SongID = song.ID
		documents[i].TitleNormalized = song.TitleNormalized
		documents[i].CreatedAt = now
		documents[i].UpdatedAt = now
	}

	err := s.songRepo.CreateSongWithDocuments(song, documents)
	if err != nil {
		return "", err
	}

	return song.ID, nil
}

// GetAllSongs retrieves all songs from the repository.
// Returns:
//   - ([]models.Song, nil) on success
//   - (nil, error) if the operation fails
func (s *SongService) GetAllSongs() ([]models.Song, error) {
	return s.songRepo.GetAllSongs()
}

// GetSongByID retrieves a song by its unique identifier.
// Returns:
//   - (*models.Song, nil) if found
//   - (nil, errors.ErrNotFound) if the song does not exist
//   - (nil, error) for unexpected failures
func (s *SongService) GetSongByID(id string) (*models.Song, error) {
	return s.songRepo.GetSongByID(id)
}

// UpdateSong applies partial updates to a song, normalizing the title if provided.
// It also updates the 'updated_at' timestamp.
// Returns:
//   - nil on success
//   - errors.ErrResourceNotFound if the song does not exist
//   - error if the update operation fails
func (s *SongService) UpdateSong(id string, updates map[string]interface{}) error {
	song, err := s.songRepo.GetSongByID(id)
	if err != nil {
		return err
	}

	if song == nil {
		return errors.ErrResourceNotFound
	}

	updates["updated_at"] = s.timeProvider.Now()

	if title, ok := updates["title"].(string); ok {
		updates["title_normalized"] = utils.Normalize(title)
	}

	if err := utils.ValidateSongUpdate(updates); err != nil {
		return err
	}

	return s.songRepo.UpdateSong(id, updates)
}

// DeleteSongWithDocuments removes a song and all associated documents.
// Returns:
//   - nil on success
//   - error if the deletion fails
func (s *SongService) DeleteSongWithDocuments(songID string) error {
	return s.songRepo.DeleteSongWithDocuments(songID)
}
