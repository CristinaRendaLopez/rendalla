package models

type Song struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Genres     []string `json:"genres"`
	UploadDate string   `json:"upload_date"`
}
