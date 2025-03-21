package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) GetAuthCredentials() (*repository.AuthCredentials, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.AuthCredentials), args.Error(1)
}
