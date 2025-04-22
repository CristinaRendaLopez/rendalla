package services_test

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
)

func ptr[T any](v T) *T {
	return &v
}

var ValidCreateSongRequest = dto.CreateSongRequest{
	Title:  "Bohemian Rhapsody",
	Author: "Queen",
	Genres: []string{"rock"},
	Documents: []dto.CreateDocumentRequest{
		ValidCreateDocumentRequest,
	},
}

var InvalidCreateSongRequest = dto.CreateSongRequest{
	Title:  "A",
	Author: "",
	Genres: []string{},
}

var ValidUpdateSongRequest = dto.UpdateSongRequest{
	Title:  ptr("We Will Rock You"),
	Genres: []string{"rock"},
}

var ValidAuthorUpdateRequest = dto.UpdateSongRequest{
	Author: ptr("Freddie Mercury"),
}

var InvalidUpdateSongRequest = dto.UpdateSongRequest{
	Title: ptr(""),
}

var SongNotFoundInvalidUpdateSongRequest = dto.UpdateSongRequest{
	Title: ptr("Radio"),
}

var MockedSong = models.Song{
	ID:     "1",
	Title:  "Bohemian Rhapsody",
	Author: "Queen",
}
