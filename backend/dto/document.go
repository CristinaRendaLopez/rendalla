package dto

type DocumentRequest struct {
	Type       string   `json:"type" binding:"required"`
	Instrument []string `json:"instrument" binding:"required,min=1,dive,min=1"`
	PDFURL     string   `json:"pdf_url" binding:"required,url"`
	AudioURL   string   `json:"audio_url,omitempty"`
}

type DocumentUpdateRequest struct {
	Type       *string  `json:"type,omitempty"`
	Instrument []string `json:"instrument,omitempty"`
	PDFURL     *string  `json:"pdf_url,omitempty"`
	AudioURL   *string  `json:"audio_url,omitempty"`
}

type DocumentResponse struct {
	Message    string `json:"message"`
	DocumentID string `json:"document_id"`
}
