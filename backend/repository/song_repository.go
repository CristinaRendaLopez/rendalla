package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// SongRepository defines operations for accessing and manipulating songs in storage.
type SongRepository interface {

	// CreateSongWithDocuments stores a new song along with its associated documents.
	// Returns:
	//   - nil on success
	//   - errors.ErrInternalServer if marshalling or persistence fails
	CreateSongWithDocuments(song models.Song, documents []models.Document) error

	// GetAllSongs returns a list of all songs in the database.
	// Returns:
	//   - ([]models.Song, nil) on success
	//   - (nil, errors.ErrInternalServer) if the query fails
	GetAllSongs() ([]models.Song, error)

	// GetSongByID retrieves a song by its unique identifier.
	// Returns:
	//   - (*models.Song, nil) if found
	//   - (nil, errors.ErrNotFound) if the song does not exist
	//   - (nil, errors.ErrInternalServer) if retrieval fails
	GetSongByID(songID string) (*models.Song, error)

	// UpdateSong applies partial updates to a song by its ID.
	// Returns:
	//   - nil on success
	//   - errors.ErrInternalServer if the update fails
	UpdateSong(songID string, updates map[string]interface{}) error

	// DeleteSongWithDocuments removes a song and all related documents.
	// Returns:
	//   - nil on success
	//   - errors.ErrNotFound if the song does not exist
	//   - errors.ErrInternalServer if the deletion fails
	DeleteSongWithDocuments(songID string) error
}
