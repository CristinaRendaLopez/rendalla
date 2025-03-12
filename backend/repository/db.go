package repository

import "github.com/CristinaRendaLopez/rendalla-backend/models"

type Database interface {
	PutSong(song models.Song) error
	GetSongByID(id string) (*models.Song, error)
	UpdateSong(id string, updates map[string]interface{}) error
	DeleteSong(id string) error
}
