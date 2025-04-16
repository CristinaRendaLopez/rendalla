package services_test

import "github.com/CristinaRendaLopez/rendalla-backend/dto"

var ValidCreateDocumentRequest = dto.CreateDocumentRequest{
	Type:       "score",
	Instrument: []string{"piano"},
	PDFURL:     "https://example.com",
}
