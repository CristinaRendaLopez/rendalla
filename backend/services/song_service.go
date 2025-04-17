package services

import (
	"fmt"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
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

// Ensure SongService implements SongServiceInterface.
var _ SongServiceInterface = (*SongService)(nil)

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
func (s *SongService) CreateSongWithDocuments(req dto.CreateSongRequest) (string, error) {
	song, documents := dto.ToSongAndDocuments(req)

	if err := dto.ValidateCreateSongRequest(req); err != nil {
		return "", fmt.Errorf("validating song and documents: %w", err)
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
		return "", fmt.Errorf("creating song with documents: %w", err)
	}

	return song.ID, nil
}

// GetAllSongs retrieves all songs from the repository.
// Returns:
//   - ([]dto.SongResponseItem, nil) on success
//   - (nil, error) if the operation fails
func (s *SongService) GetAllSongs() ([]dto.SongResponseItem, error) {
	songs, err := s.songRepo.GetAllSongs()
	if err != nil {
		return nil, fmt.Errorf("retrieving all songs: %w", err)
	}
	return dto.ToSongResponseList(songs), nil
}

// GetSongByID retrieves a song by its unique identifier.
// Returns:
//   - (dto.SongResponseItem, nil) if found
//   - (nil, errors.ErrNotFound) if the song does not exist
//   - (nil, error) for unexpected failures
func (s *SongService) GetSongByID(id string) (dto.SongResponseItem, error) {
	song, err := s.songRepo.GetSongByID(id)
	if err != nil {
		return dto.SongResponseItem{}, fmt.Errorf("retrieving song with ID %s: %w", id, err)
	}
	return dto.ToSongResponseItem(*song), nil
}

// UpdateSong applies partial updates to a song, normalizing the title if provided.
// It also updates the 'updated_at' timestamp.
// Returns:
//   - nil on success
//   - errors.ErrResourceNotFound if the song does not exist
//   - error if the update operation fails
func (s *SongService) UpdateSong(id string, updates dto.UpdateSongRequest) error {
	_, err := s.songRepo.GetSongByID(id)
	if err != nil {
		return fmt.Errorf("checking existence of song %s: %w", id, err)
	}

	updateMap := make(map[string]interface{})
	if updates.Title != nil {
		updateMap["title"] = *updates.Title
		updateMap["title_normalized"] = utils.Normalize(*updates.Title)
	}
	if updates.Author != nil {
		updateMap["author"] = *updates.Author
	}
	if updates.Genres != nil {
		updateMap["genres"] = updates.Genres
	}
	updateMap["updated_at"] = s.timeProvider.Now()

	if err := dto.ValidateUpdateSongRequest(updates); err != nil {
		return fmt.Errorf("validating song update: %w", err)
	}

	if err := s.songRepo.UpdateSong(id, updateMap); err != nil {
		return fmt.Errorf("updating song %s: %w", id, err)
	}

	return nil
}

// DeleteSongWithDocuments removes a song and all associated documents.
// Returns:
//   - nil on success
//   - errors.ErrResourceNotFound if the song does not exist
//   - error if the deletion fails
func (s *SongService) DeleteSongWithDocuments(songID string) error {
	_, err := s.songRepo.GetSongByID(songID)
	if err != nil {
		return fmt.Errorf("checking existence of song %s: %w", songID, err)
	}

	if err := s.songRepo.DeleteSongWithDocuments(songID); err != nil {
		return fmt.Errorf("deleting song %s with documents: %w", songID, err)
	}
	return nil
}
