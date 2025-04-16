package handlers_test

// Valid song data
const SongValidJSON = `
{
	"title": "Bohemian Rhapsody",
	"author": "Queen",
	"genres": ["rock", "opera"],
	"documents": [
		{
			"type": "score",
			"instrument": ["piano"],
			"pdf_url": "https://example.com/bohemian-piano.pdf"
		},
		{
			"type": "tablature",
			"instrument": ["guitar"],
			"pdf_url": "https://example.com/bohemian-guitar.pdf"
		}
	]
}`

// Good JSON syntax but invalid data
const SongInvalidDataJSON = `
{
	"title": "Invisible Man",
	"author": "Queen",
	"genres": ["rock"],
	"documents": [
		{
			"type": "tablature",
			"instrument": [],
			"pdf_url": "https://example.com/invalid.pdf"
		}
	]
}`

// Bad JSON syntax
const SongInvalidJSON = `
{
	"title": "Another One Bites The Dust",
	"author": "Queen"
` // missing brace
