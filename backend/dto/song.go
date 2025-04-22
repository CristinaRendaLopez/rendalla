package dto

type CreateSongRequest struct {
	Title     string                  `json:"title" binding:"required,min=3"`
	Author    string                  `json:"author" binding:"required"`
	Genres    []string                `json:"genres" binding:"required,dive,min=3"`
	Documents []CreateDocumentRequest `json:"documents,omitempty"`
}

type UpdateSongRequest struct {
	Title  *string  `json:"title,omitempty"`
	Author *string  `json:"author,omitempty"`
	Genres []string `json:"genres,omitempty"`
}

type CreateSongResponse struct {
	Message string `json:"message"`
	SongID  string `json:"song_id"`
}

type SongResponseItem struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Genres []string `json:"genres"`
}
