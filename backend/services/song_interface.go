package services

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// SongServiceInterface defines application-level operations for managing songs and their associated documents.
type SongServiceInterface interface {

	// CreateSongWithDocuments creates a new song and stores all associated documents.
	// Automatically generates IDs and timestamps, and normalizes the title.
	// Returns:
	//   - the generated song ID on success
	//   - error if the operation fails
	CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error)

	// GetAllSongs returns all songs stored in the system.
	// Returns:
	//   - ([]models.Song, nil) on success
	//   - (nil, error) if the retrieval fails
	GetAllSongs() ([]models.Song, error)

	// GetSongByID retrieves a single song by its unique identifier.
	// Returns:
	//   - (*models.Song, nil) if found
	//   - (nil, errors.ErrNotFound) if the song does not exist
	//   - (nil, error) for unexpected errors
	GetSongByID(songID string) (*models.Song, error)

	// UpdateSong applies partial updates to a song, including optional title normalization.
	// Returns:
	//   - nil on success
	//   - errors.ErrNotFound if the song does not exist
	//   - error if the update fails
	UpdateSong(songID string, updates map[string]interface{}) error

	// DeleteSongWithDocuments removes a song and all documents linked to it.
	// Returns:
	//   - nil on success
	//   - error if the operation fails
	DeleteSongWithDocuments(songID string) error
}
