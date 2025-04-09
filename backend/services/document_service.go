package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

type DocumentServiceInterface interface {
	GetDocumentsBySongID(songID string) ([]models.Document, error)
	GetDocumentByID(songID string, docID string) (*models.Document, error)
	CreateDocument(document models.Document) (string, error)
	UpdateDocument(songID string, docID string, updates map[string]interface{}) error
	DeleteDocument(songID string, docID string) error
}

type DocumentService struct {
	repo         repository.DocumentRepository
	songRepo     repository.SongRepository
	idGen        IDGenerator
	timeProvider TimeProvider
}

func NewDocumentService(repo repository.DocumentRepository, songRepo repository.SongRepository, idGen IDGenerator, timeProvider TimeProvider) *DocumentService {
	return &DocumentService{
		repo:         repo,
		songRepo:     songRepo,
		idGen:        idGen,
		timeProvider: timeProvider,
	}
}

func (s *DocumentService) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	return s.repo.GetDocumentsBySongID(songID)
}

func (s *DocumentService) GetDocumentByID(songID string, docID string) (*models.Document, error) {
	return s.repo.GetDocumentByID(songID, docID)
}

func (s *DocumentService) CreateDocument(document models.Document) (string, error) {
	song, err := s.songRepo.GetSongByID(document.SongID)
	if err != nil {
		return "", err
	}

	document.TitleNormalized = utils.Normalize(song.Title)
	document.ID = s.idGen.NewID()
	now := s.timeProvider.Now()
	document.CreatedAt = now
	document.UpdatedAt = now

	return document.ID, s.repo.CreateDocument(document)
}

func (s *DocumentService) UpdateDocument(songID, docID string, updates map[string]interface{}) error {
	if _, ok := updates["title_normalized"]; !ok {
		song, err := s.songRepo.GetSongByID(songID)
		if err != nil {
			return err
		}
		updates["title_normalized"] = utils.Normalize(song.Title)
	}
	updates["updated_at"] = s.timeProvider.Now()
	return s.repo.UpdateDocument(songID, docID, updates)
}

func (s *DocumentService) DeleteDocument(songID string, docID string) error {
	return s.repo.DeleteDocument(songID, docID)
}
