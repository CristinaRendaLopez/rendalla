package models

// Document represents a musical score or tablature associated with a song.
type Document struct {
	ID              string   `json:"id" dynamodbav:"id" dynamo:"id"`                                // Unique identifier for the document
	SongID          string   `json:"song_id" dynamodbav:"song_id" dynamo:"song_id"`                 // Foreign key referencing the associated song
	TitleNormalized string   `json:"-" dynamodbav:"title_normalized" dynamo:"title_normalized"`     // Normalized title (inherited from the song) used for search and pagination
	Type            string   `json:"type" dynamodbav:"type" dynamo:"type"`                          // Document type: "score" or "tablature"
	Instrument      []string `json:"instrument" dynamodbav:"instrument" dynamo:"instrument"`        // Target instruments or voices (e.g., "guitar", "soprano")
	PDFURL          string   `json:"pdf_url" dynamodbav:"pdf_url" dynamo:"pdf_url"`                 // URL to the PDF file stored in S3
	AudioURL        string   `json:"audio_url,omitempty" dynamodbav:"audio_url" dynamo:"audio_url"` // Optional URL to an accompanying audio file
	CreatedAt       string   `json:"created_at" dynamodbav:"created_at" dynamo:"created_at"`        // ISO timestamp of creation
	UpdatedAt       string   `json:"updated_at" dynamodbav:"updated_at" dynamo:"updated_at"`        // ISO timestamp of last update
}
