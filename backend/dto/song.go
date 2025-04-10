package dto

type SongRequest struct {
	Title     string            `json:"title" binding:"required,min=3"`
	Author    string            `json:"author" binding:"required"`
	Genres    []string          `json:"genres" binding:"required,dive,min=3"`
	Documents []DocumentRequest `json:"documents,omitempty"`
}

type SongUpdateRequest struct {
	Title  *string  `json:"title,omitempty"`
	Author *string  `json:"author,omitempty"`
	Genres []string `json:"genres,omitempty"`
}

type SongResponse struct {
	Message string `json:"message"`
	SongID  string `json:"song_id"`
}
