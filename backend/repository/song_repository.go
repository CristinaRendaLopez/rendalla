package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

type SongRepository interface {
	GetSongByID(id string) (*models.Song, error)
	CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error)
	UpdateSong(id string, updates map[string]interface{}) error
	DeleteSongWithDocuments(id string) error
	GetAllSongs() ([]models.Song, error)
}
