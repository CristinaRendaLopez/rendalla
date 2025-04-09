package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

type SongServiceInterface interface {
	GetAllSongs() ([]models.Song, error)
	GetSongByID(songID string) (*models.Song, error)
	CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error)
	UpdateSong(songID string, updates map[string]interface{}) error
	DeleteSongWithDocuments(songID string) error
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
		return "", err
	}

	return song.ID, nil
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

	if title, ok := updates["title"].(string); ok {
		updates["title_normalized"] = utils.Normalize(title)
	}

	return s.songRepo.UpdateSong(id, updates)
}

func (s *SongService) DeleteSongWithDocuments(songID string) error {
	return s.songRepo.DeleteSongWithDocuments(songID)
}
