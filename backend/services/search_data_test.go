package services_test

import "github.com/CristinaRendaLopez/rendalla-backend/models"

var SongBohemianRhapsody = models.Song{
	ID:        "1",
	Title:     "Bohemian Rhapsody",
	Author:    "Queen",
	Genres:    []string{"rock"},
	CreatedAt: "1975-10-31T00:00:00Z",
}

var SongRadioGaGa = models.Song{
	ID:        "2",
	Title:     "Radio Ga Ga",
	Author:    "Queen",
	Genres:    []string{"pop"},
	CreatedAt: "1984-01-01T00:00:00Z",
}

var DocumentPianoScore = models.Document{
	ID:              "d1",
	SongID:          "s1",
	Type:            "score",
	Instrument:      []string{"piano"},
	PDFURL:          "https://example.com/piano.pdf",
	CreatedAt:       "2020-01-01T00:00:00Z",
	UpdatedAt:       "2020-01-02T00:00:00Z",
	TitleNormalized: "bohemian rhapsody",
}

var DocumentGuitarTab = models.Document{
	ID:              "d2",
	SongID:          "s2",
	Type:            "tablature",
	Instrument:      []string{"guitar"},
	PDFURL:          "https://example.com/guitar.pdf",
	CreatedAt:       "2021-01-01T00:00:00Z",
	UpdatedAt:       "2021-01-02T00:00:00Z",
	TitleNormalized: "somebody to love",
}

var ValidNextToken = map[string]string{"last_id": "2"}
var ReturnedNextToken = map[string]string{"last_id": "5"}
