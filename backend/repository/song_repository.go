package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

// SongRepository defines operations for accessing and manipulating songs in storage.
type SongRepository interface {

	// CreateSongWithDocuments stores a new song along with its associated documents.
	CreateSongWithDocuments(song models.Song, documents []models.Document) error

	// GetAllSongs returns a list of all songs.
	GetAllSongs() ([]models.Song, error)

	// GetSongByID retrieves a song by its unique identifier.
	GetSongByID(songID string) (*models.Song, error)

	// UpdateSong applies partial updates to a song by its ID.
	UpdateSong(songID string, updates map[string]interface{}) error

	// DeleteSongWithDocuments removes a song and all related documents.
	DeleteSongWithDocuments(songID string) error
}
