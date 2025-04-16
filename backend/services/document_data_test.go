package services_test

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

var ValidCreateDocumentRequest = dto.CreateDocumentRequest{
	Type:       "score",
	Instrument: []string{"piano"},
	PDFURL:     "https://example.com/bohemian-piano.pdf",
	SongID:     "song-123",
}

var InvalidCreateDocumentRequest_EmptyInstrument = dto.CreateDocumentRequest{
	Type:       "score",
	Instrument: []string{},
	PDFURL:     "https://example.com/invalid.pdf",
	SongID:     "song-123",
}

var MockedDocument = models.Document{
	ID:              "doc-1",
	SongID:          "song-123",
	Type:            "score",
	Instrument:      []string{"piano"},
	PDFURL:          "https://example.com/bohemian-piano.pdf",
	TitleNormalized: "bohemian rhapsody",
	CreatedAt:       "now",
	UpdatedAt:       "now",
}

var DocumentResponse = dto.DocumentResponseItem{
	ID:         "doc-1",
	SongID:     "song-123",
	Type:       "score",
	Instrument: []string{"piano"},
	PDFURL:     "https://example.com/bohemian-piano.pdf",
	CreatedAt:  "now",
	UpdatedAt:  "now",
}

var ValidUpdateDocumentRequest = dto.UpdateDocumentRequest{
	Type:       "tablature",
	Instrument: []string{"guitar"},
}

var ValidUpdateDocumentRequestPDFAndAudio = dto.UpdateDocumentRequest{
	PDFURL:   "https://example.com/bohemian-piano-2.pdf",
	AudioURL: "https://example.com/bohemian-piano.mp3",
}

var InvalidUpdateDocumentRequest = dto.UpdateDocumentRequest{
	Instrument: []string{}, // empty
}

var RelatedSong = models.Song{
	ID:    "song-123",
	Title: "Bohemian Rhapsody",
}
