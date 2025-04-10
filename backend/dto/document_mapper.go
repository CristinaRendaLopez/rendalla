package dto

import (
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

func ToDocumentModel(dto DocumentRequest) models.Document {
	return models.Document{
		Type:       dto.Type,
		Instrument: dto.Instrument,
		PDFURL:     dto.PDFURL,
		AudioURL:   dto.AudioURL,
	}
}
