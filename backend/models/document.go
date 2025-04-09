package models

type Document struct {
	ID              string   `json:"id" dynamodbav:"id" dynamo:"id"`
	SongID          string   `json:"song_id" dynamodbav:"song_id" dynamo:"song_id"`
	TitleNormalized string   `json:"title_normalized" dynamodbav:"title_normalized" dynamo:"title_normalized"`
	Type            string   `json:"type" dynamodbav:"type" dynamo:"type"`
	Instrument      []string `json:"instrument" dynamodbav:"instrument" dynamo:"instrument"`
	PDFURL          string   `json:"pdf_url" dynamodbav:"pdf_url" dynamo:"pdf_url"`
	AudioURL        string   `json:"audio_url,omitempty" dynamodbav:"audio_url" dynamo:"audio_url"`
	CreatedAt       string   `json:"created_at" dynamodbav:"created_at" dynamo:"created_at"`
	UpdatedAt       string   `json:"updated_at" dynamodbav:"updated_at" dynamo:"updated_at"`
}
