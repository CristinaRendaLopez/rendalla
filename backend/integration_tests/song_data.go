package integration_tests

import "github.com/CristinaRendaLopez/rendalla-backend/dto"

type songListResponse struct {
	Data []dto.SongResponseItem `json:"data"`
}

type songDetailResponse struct {
	Data dto.SongResponseItem `json:"data"`
}

var BohemianDocs = []dto.CreateDocumentRequest{
	{
		Type:       "score",
		Instrument: []string{"piano"},
		PDFURL:     "https://test-bucket/bohemian-piano.pdf",
	},
	{
		Type:       "tablature",
		Instrument: []string{"guitar"},
		PDFURL:     "https://test-bucket/bohemian-guitar.pdf",
	},
}

var DontStopMeDocs = []dto.CreateDocumentRequest{
	{
		Type:       "score",
		Instrument: []string{"voice"},
		PDFURL:     "https://test-bucket/dontstop-voice.pdf",
	},
	{
		Type:       "tablature",
		Instrument: []string{"bass"},
		PDFURL:     "https://test-bucket/dontstop-bass.pdf",
	},
}

var BohemianRhapsodyPayload = dto.CreateSongRequest{
	Title:     "Bohemian Rhapsody",
	Author:    "Queen",
	Genres:    []string{"Rock", "Opera"},
	Documents: BohemianDocs,
}

var DontStopMeNowPayload = dto.CreateSongRequest{
	Title:     "Don't stop me now",
	Author:    "Queen",
	Genres:    []string{"Rock", "Pop"},
	Documents: DontStopMeDocs,
}

var WeAreTheChampionsPayload = dto.CreateSongRequest{
	Title:  "We are the Champions",
	Author: "Queen",
	Genres: []string{"Rock"},
}
