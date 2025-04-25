package dto

type CreateDocumentRequest struct {
	Type       string   `json:"type" binding:"required"`
	Instrument []string `json:"instrument" binding:"required,min=1,dive,min=1"`
	PDFURL     string   `json:"pdf_url" binding:"required,url"`
	AudioURL   string   `json:"audio_url,omitempty"`
	SongID     string   `json:"-"`
}

type CreateDocumentResponse struct {
	Message    string `json:"message"`
	DocumentID string `json:"document_id"`
}

type UpdateDocumentRequest struct {
	Type       string   `json:"type,omitempty"`
	Instrument []string `json:"instrument,omitempty"`
	PDFURL     string   `json:"pdf_url,omitempty"`
	AudioURL   string   `json:"audio_url,omitempty"`
}

type DocumentResponseItem struct {
	ID         string   `json:"id"`
	SongID     string   `json:"song_id"`
	Type       string   `json:"type"`
	Instrument []string `json:"instrument"`
	PDFURL     string   `json:"pdf_url"`
	AudioURL   string   `json:"audio_url,omitempty"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}
