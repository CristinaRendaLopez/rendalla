package handlers_test

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
)

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

var MultipartFieldsValid = map[string]string{
	"type":         "score",
	"instrument[]": "piano",
}

var MultipartFieldsInvalidInstrument = map[string]string{
	"type":         "tablature",
	"instrument[]": "", // invalid (empty)
}

var MultipartFieldsInvalid = map[string]string{
	"type": "",
}

var MultipartPDFMock = utils.TestFile{
	Filename: "bohemian.pdf",
	Content:  []byte("%PDF-1.4\nfake pdf content\n..."),
}

var MultipartPDFInvalid = utils.TestFile{
	Filename: "not_a_pdf.txt",
	Content:  []byte("just some text, not a pdf"),
}
