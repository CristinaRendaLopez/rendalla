package models

type Document struct {
	ID         string   `json:"id"`
	SongID     string   `json:"song_id"`
	Type       string   `json:"type"`
	Instrument []string `json:"instrument"`
	PDFURL     string   `json:"pdf_url"`
	AudioURL   string   `json:"audio_url,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}
