package mocks

import (
	"bytes"
)

type MockMultipartFile struct {
	*bytes.Reader
}

func NewFakeMultipartFile(content []byte) *MockMultipartFile {
	return &MockMultipartFile{
		Reader: bytes.NewReader(content),
	}
}

func (f *MockMultipartFile) Close() error {
	return nil
}

func (f *MockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return f.Reader.ReadAt(p, off)
}

func (f *MockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return f.Reader.Seek(offset, whence)
}
