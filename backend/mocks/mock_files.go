package mocks

import (
	"mime/multipart"
)

type MockFileService struct{}

func (m *MockFileService) UploadPDFToS3(file multipart.File, header *multipart.FileHeader, songID string) (string, error) {
	return "https://fake-s3.com/songs/" + songID + "/doc_mock.pdf", nil
}
