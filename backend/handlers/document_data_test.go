package handlers_test

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
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

func MockPDFFile() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "mock.pdf",
		Header:   textproto.MIMEHeader{"Content-Type": []string{"application/pdf"}},
		Size:     2048,
	}
}

func MockInvalidFile() *multipart.FileHeader {
	return &multipart.FileHeader{
		Filename: "mock.txt",
		Header:   textproto.MIMEHeader{"Content-Type": []string{"text/plain"}},
		Size:     2048,
	}
}

func buildMultipartRequest(songID, docType string, instruments []string, fileHeader *multipart.FileHeader) (*http.Request, *multipart.FileHeader) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("type", docType)
	for _, inst := range instruments {
		_ = writer.WriteField("instrument[]", inst)
	}

	part, _ := writer.CreateFormFile("pdf", fileHeader.Filename)
	part.Write([]byte("PDF binary content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/songs/%s/documents", songID), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, fileHeader
}
