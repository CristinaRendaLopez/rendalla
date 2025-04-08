package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

type SongRepository interface {
	GetSongByID(songID string) (*models.Song, error)
	CreateSongWithDocuments(song models.Song, documents []models.Document) error
	UpdateSong(songID string, updates map[string]interface{}) error
	DeleteSongWithDocuments(songID string) error
	GetAllSongs() ([]models.Song, error)
}
