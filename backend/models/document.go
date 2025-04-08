package models

type Document struct {
	ID         string   `json:"id" dynamodbav:"id"`
	SongID     string   `json:"song_id" dynamodbav:"song_id"`
	Type       string   `json:"type" dynamodbav:"type"`
	Instrument []string `json:"instrument" dynamodbav:"instrument"`
	PDFURL     string   `json:"pdf_url" dynamodbav:"pdf_url"`
	AudioURL   string   `json:"audio_url,omitempty" dynamodbav:"audio_url"`
	CreatedAt  string   `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt  string   `json:"updated_at" dynamodbav:"updated_at"`
}
