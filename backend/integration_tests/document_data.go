package integration_tests

import "github.com/CristinaRendaLopez/rendalla-backend/dto"

type DocumentDetailResponse struct {
	Data dto.DocumentResponseItem `json:"data"`
}

type DocumentListResponse struct {
	Data []dto.DocumentResponseItem `json:"data"`
}

// Valid document creation
var ViolinScore = dto.CreateDocumentRequest{
	Type:       "score",
	Instrument: []string{"violin"},
	PDFURL:     "https://s3.test/violin.pdf",
	AudioURL:   "https://s3.test/violin.mp3",
}

var FluteScore = dto.CreateDocumentRequest{
	Type:       "score",
	Instrument: []string{"flute"},
	PDFURL:     "https://s3.test/test_flute.pdf",
}

// Valid document update
var TablatureUpdate = dto.UpdateDocumentRequest{
	Type:       "tablature",
	Instrument: []string{"guitar", "lead"},
	PDFURL:     "https://test-updated.com/guitar-lead.pdf",
	AudioURL:   "https://test-updated.com/audio.mp3",
}

// Malformed JSON
var InvalidJSONDocument = `{"type":`

// Valid JSON, invalid content (empty fields)
var InvalidFieldsDocument = `{
	"type": "",
	"instrument": [],
	"pdf_url": ""
}`
