package models

type Song struct {
	ID         string   `json:"id" dynamo:"id"`
	Title      string   `json:"title" dynamodbav:"title"`
	Author     string   `json:"author" dynamodbav:"author"`
	Genres     []string `json:"genres" dynamodbav:"genres"`
	YoutubeURL string   `json:"youtube_url,omitempty" dynamodbav:"youtube_url"`
	CreatedAt  string   `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt  string   `json:"updated_at" dynamodbav:"updated_at"`
}
