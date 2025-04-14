package dto

import (
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

func ToDocumentModel(dto CreateDocumentRequest) models.Document {
	return models.Document{
		SongID:     dto.SongID,
		Type:       dto.Type,
		Instrument: dto.Instrument,
		PDFURL:     dto.PDFURL,
		AudioURL:   dto.AudioURL,
	}
}

func ToDocumentResponseItem(m models.Document) DocumentResponseItem {
	return DocumentResponseItem{
		ID:         m.ID,
		SongID:     m.SongID,
		Type:       m.Type,
		Instrument: m.Instrument,
		PDFURL:     m.PDFURL,
		AudioURL:   m.AudioURL,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func ToDocumentResponseList(docs []models.Document) []DocumentResponseItem {
	out := make([]DocumentResponseItem, len(docs))
	for i, d := range docs {
		out[i] = ToDocumentResponseItem(d)
	}
	return out
}

// ValidateCreateDocumentRequest validates DocumentRequest DTO.
func ValidateCreateDocumentRequest(doc CreateDocumentRequest) error {
	if utils.IsEmptyString(doc.Type) {
		return errors.ErrValidationFailed
	}
	if utils.IsEmptyString(doc.PDFURL) {
		return errors.ErrValidationFailed
	}
	if len(doc.Instrument) == 0 {
		return errors.ErrValidationFailed
	}
	for _, inst := range doc.Instrument {
		if utils.IsEmptyString(inst) {
			return errors.ErrValidationFailed
		}
	}
	return nil
}

// ValidateUpdateDocumentRequest validates a partial DocumentRequest used for updates.
func ValidateUpdateDocumentRequest(doc UpdateDocumentRequest) error {
	if doc.Type != "" && utils.IsEmptyString(doc.Type) {
		return errors.ErrValidationFailed
	}
	if doc.PDFURL != "" && utils.IsEmptyString(doc.PDFURL) {
		return errors.ErrValidationFailed
	}
	if len(doc.Instrument) > 0 {
		for _, inst := range doc.Instrument {
			if utils.IsEmptyString(inst) {
				return errors.ErrValidationFailed
			}
		}
	}
	return nil
}
