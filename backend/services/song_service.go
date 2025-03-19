package services

import (
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SongService struct {
	songRepo repository.SongRepository
	docRepo  repository.DocumentRepository
}

func NewSongService(songRepo repository.SongRepository, docRepo repository.DocumentRepository) *SongService {
	return &SongService{songRepo: songRepo, docRepo: docRepo}
}

func (s *SongService) GetAllSongs() ([]models.Song, error) {
	return s.songRepo.GetAllSongs()
}

func (s *SongService) GetSongByID(id string) (*models.Song, error) {
	return s.songRepo.GetSongByID(id)
}

func (s *SongService) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	song.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	song.CreatedAt, song.UpdatedAt = now, now

	songID, err := s.songRepo.CreateSongWithDocuments(song, documents)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		doc.SongID = songID
		_, err := s.docRepo.CreateDocument(doc)
		if err != nil {
			logrus.WithError(err).Error("Failed to create document for song")
			return "", err
		}
	}

	logrus.WithField("song_id", songID).Info("Song and documents created successfully")
	return songID, nil
}

func (s *SongService) UpdateSong(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	return s.songRepo.UpdateSong(id, updates)
}

func (s *SongService) DeleteSongWithDocuments(songID string) error {
	return s.songRepo.DeleteSongWithDocuments(songID)
}
