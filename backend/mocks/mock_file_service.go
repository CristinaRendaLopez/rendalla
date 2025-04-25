package mocks

import (
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadPDFToS3(file multipart.File, header *multipart.FileHeader, songID string) (string, error) {
	args := m.Called(file, header, songID)
	return args.String(0), args.Error(1)
}
