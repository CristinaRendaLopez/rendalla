package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
)

// SongServiceInterface defines application-level operations for managing songs and their associated documents.
type SongServiceInterface interface {

	// CreateSongWithDocuments creates a new song and stores all associated documents.
	// Automatically generates IDs and timestamps, and normalizes the title.
	// Returns:
	//   - the generated song ID on success
	//   - error if the operation fails
	CreateSongWithDocuments(dto.CreateSongRequest) (string, error)

	// GetAllSongs returns all songs stored in the system.
	// Returns:
	//   - ([]dto.SongResponseItem, nil) on success
	//   - (nil, error) if the retrieval fails
	GetAllSongs() ([]dto.SongResponseItem, error)

	// GetSongByID retrieves a single song by its unique identifier.
	// Returns:
	//   - (dto.SongResponseItem, nil) if found
	//   - (nil, errors.ErrNotFound) if the song does not exist
	//   - (nil, error) for unexpected errors
	GetSongByID(songID string) (dto.SongResponseItem, error)

	// UpdateSong applies partial updates to a song, including optional title normalization.
	// Returns:
	//   - nil on success
	//   - errors.ErrNotFound if the song does not exist
	//   - error if the update fails
	UpdateSong(songID string, updates dto.UpdateSongRequest) error

	// DeleteSongWithDocuments removes a song and all documents linked to it.
	// Returns:
	//   - nil on success
	//   - error if the operation fails
	DeleteSongWithDocuments(songID string) error
}
