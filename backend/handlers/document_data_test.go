package handlers_test

import "github.com/CristinaRendaLopez/rendalla-backend/dto"

// Valid document data
const DocumentValidJSON = `
{
	"type": "score",
	"instrument": ["piano"],
	"pdf_url": "https://example.com/bohemian-piano.pdf"
}`

// Good JSON syntax but invalid data
const DocumentInvalidDataJSON = `
{
	"type": "tablature",
	"instrument": [],
	"pdf_url": "https://example.com/bad.pdf"
}`

// Bad JSON syntax
const DocumentInvalidJSON = `
{
	"type": "tablature",
	"instrument": ["guitar"]
` // missing brace

var DocumentResponseScore = dto.DocumentResponseItem{
	ID:         "doc-1",
	SongID:     "1",
	Type:       "score",
	Instrument: []string{"piano"},
	PDFURL:     "https://example.com/bohemian-piano.pdf",
}

var DocumentResponseTablature = dto.DocumentResponseItem{
	ID:         "doc-2",
	SongID:     "1",
	Type:       "tablature",
	Instrument: []string{"guitar"},
	PDFURL:     "https://example.com/bohemian-guitar.pdf",
}

const DocumentUpdateValidJSON = `
{
	"type": "tablature",
	"instrument": ["guitar"]
}`

const DocumentUpdateInvalidJSON = `
{
	"type": "tablature",
` // missing brace

const DocumentUpdateEmptyInstrumentJSON = `
{
	"instrument": []
}`

const DocumentUpdateOnlyTypeJSON = `
{
	"type": "score"
}`
