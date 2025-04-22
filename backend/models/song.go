package models

// Song represents a musical track with metadata used for display and search purposes.
type Song struct {
	ID              string   `json:"id" dynamodbav:"id" dynamo:"id"`                                      // Unique identifier for the song
	Title           string   `json:"title" dynamodbav:"title" dynamo:"title"`                             // Original title as entered by the user
	TitleNormalized string   `json:"-" dynamodbav:"title_normalized" dynamo:"title_normalized"`           // Lowercased, accent-stripped version of the title for search optimization
	Author          string   `json:"author" dynamodbav:"author" dynamo:"author"`                          // Author or composer of the song
	Genres          []string `json:"genres" dynamodbav:"genres" dynamo:"genres"`                          // List of associated genres (e.g., classical, rock)
	YoutubeURL      string   `json:"youtube_url,omitempty" dynamodbav:"youtube_url" dynamo:"youtube_url"` // Optional link to a YouTube video
	CreatedAt       string   `json:"created_at" dynamodbav:"created_at" dynamo:"created_at"`              // ISO timestamp of creation
	UpdatedAt       string   `json:"updated_at" dynamodbav:"updated_at" dynamo:"updated_at"`              // ISO timestamp of last update
}
