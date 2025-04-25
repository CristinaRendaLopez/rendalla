package mocks

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockFileExtractor struct {
	mock.Mock
}

func (m *MockFileExtractor) GetHeader(c *gin.Context, field string) (*multipart.FileHeader, error) {
	args := m.Called(c, field)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*multipart.FileHeader), args.Error(1)
}

func (m *MockFileExtractor) OpenFile(header *multipart.FileHeader) (multipart.File, error) {
	args := m.Called(header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(multipart.File), args.Error(1)
}
