package dto

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

func ToSongModel(dto SongRequest) models.Song {
	return models.Song{
		Title:  dto.Title,
		Author: dto.Author,
		Genres: dto.Genres,
	}
}

func ToDocumentModels(dtos []DocumentRequest) []models.Document {
	documents := make([]models.Document, len(dtos))
	for i, d := range dtos {
		documents[i] = models.Document{
			Type:       d.Type,
			Instrument: d.Instrument,
			PDFURL:     d.PDFURL,
			AudioURL:   d.AudioURL,
		}
	}
	return documents
}
