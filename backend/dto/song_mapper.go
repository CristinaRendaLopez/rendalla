package dto

import (
	"strings"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

func ToSongAndDocuments(req CreateSongRequest) (models.Song, []models.Document) {
	song := models.Song{
		Title:  req.Title,
		Author: req.Author,
		Genres: req.Genres,
	}

	documents := make([]models.Document, len(req.Documents))
	for i, d := range req.Documents {
		documents[i] = models.Document{
			Type:       d.Type,
			Instrument: d.Instrument,
			PDFURL:     d.PDFURL,
			AudioURL:   d.AudioURL,
		}
	}

	return song, documents
}

func ToSongResponseItem(m models.Song) SongResponseItem {
	return SongResponseItem{
		ID:     m.ID,
		Title:  m.Title,
		Author: m.Author,
		Genres: m.Genres,
	}
}

func ToSongResponseList(songs []models.Song) []SongResponseItem {
	out := make([]SongResponseItem, len(songs))
	for i, s := range songs {
		out[i] = SongResponseItem{
			ID:     s.ID,
			Title:  s.Title,
			Author: s.Author,
			Genres: s.Genres,
		}
	}
	return out
}

// ValidateCreateSongRequest validates CreateSongRequest DTO.
func ValidateCreateSongRequest(req CreateSongRequest) error {
	if utils.IsEmptyString(req.Title) || len(req.Title) < 3 {
		return errors.ErrValidationFailed
	}
	if utils.IsEmptyString(req.Author) {
		return errors.ErrValidationFailed
	}
	if len(req.Genres) == 0 {
		return errors.ErrValidationFailed
	}
	for _, g := range req.Genres {
		if len(strings.TrimSpace(g)) < 3 {
			return errors.ErrValidationFailed
		}
	}
	for _, doc := range req.Documents {
		if err := ValidateCreateDocumentRequest(doc); err != nil {
			return err
		}
	}
	return nil
}

// ValidateUpdateSongRequest validates UpdateSongRequest DTO.
func ValidateUpdateSongRequest(update UpdateSongRequest) error {
	if update.Title == nil && update.Author == nil && len(update.Genres) == 0 {
		return errors.ErrValidationFailed
	}
	if update.Title != nil && (utils.IsEmptyString(*update.Title) || len(*update.Title) < 3) {
		return errors.ErrValidationFailed
	}
	if update.Author != nil && utils.IsEmptyString(*update.Author) {
		return errors.ErrValidationFailed
	}
	if len(update.Genres) > 0 {
		for _, g := range update.Genres {
			if len(strings.TrimSpace(g)) < 3 {
				return errors.ErrValidationFailed
			}
		}
	}
	return nil
}
