package models

type Song struct {
	ID              string   `json:"id" dynamodbav:"id" dynamo:"id"`
	Title           string   `json:"title" dynamodbav:"title" dynamo:"title"`
	TitleNormalized string   `json:"-"`
	Author          string   `json:"author" dynamodbav:"author" dynamo:"author"`
	Genres          []string `json:"genres" dynamodbav:"genres" dynamo:"genres"`
	YoutubeURL      string   `json:"youtube_url,omitempty" dynamodbav:"youtube_url" dynamo:"youtube_url"`
	CreatedAt       string   `json:"created_at" dynamodbav:"created_at" dynamo:"created_at"`
	UpdatedAt       string   `json:"updated_at" dynamodbav:"updated_at" dynamo:"updated_at"`
}
