package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

type SongServiceInterface interface {
	GetAllSongs() ([]models.Song, error)
	GetSongByID(id string) (*models.Song, error)
	CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error)
	UpdateSong(id string, updates map[string]interface{}) error
	DeleteSongWithDocuments(id string) error
}

type IDGenerator interface {
	NewID() string
}

type TimeProvider interface {
	Now() string
}

type SongService struct {
	songRepo     repository.SongRepository
	docRepo      repository.DocumentRepository
	idGen        IDGenerator
	timeProvider TimeProvider
}

func NewSongService(
	songRepo repository.SongRepository,
	docRepo repository.DocumentRepository,
	idGen IDGenerator,
	timeProvider TimeProvider,
) *SongService {
	return &SongService{
		songRepo:     songRepo,
		docRepo:      docRepo,
		idGen:        idGen,
		timeProvider: timeProvider,
	}
}

func (s *SongService) GetAllSongs() ([]models.Song, error) {
	return s.songRepo.GetAllSongs()
}

func (s *SongService) GetSongByID(id string) (*models.Song, error) {
	return s.songRepo.GetSongByID(id)
}

func (s *SongService) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	song.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	song.CreatedAt, song.UpdatedAt = now, now

	songID, err := s.songRepo.CreateSongWithDocuments(song, documents)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		doc.SongID = songID
		doc.ID = s.idGen.NewID()
		doc.CreatedAt, doc.UpdatedAt = now, now

		if _, err := s.docRepo.CreateDocument(doc); err != nil {
			return "", err
		}
	}

	return songID, nil
}

func (s *SongService) UpdateSong(id string, updates map[string]interface{}) error {
	song, err := s.songRepo.GetSongByID(id)
	if err != nil {
		return err
	}
	if song == nil {
		return utils.ErrResourceNotFound
	}

	updates["updated_at"] = s.timeProvider.Now()
	return s.songRepo.UpdateSong(id, updates)
}

func (s *SongService) DeleteSongWithDocuments(songID string) error {
	return s.songRepo.DeleteSongWithDocuments(songID)
}
