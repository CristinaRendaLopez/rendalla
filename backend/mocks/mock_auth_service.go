package mocks

import (
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

var _ services.AuthServiceInterface = (*MockAuthService)(nil)

func (m *MockAuthService) AuthenticateUser(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GetAuthCredentials() (*repository.AuthCredentials, error) {
	args := m.Called()
	return args.Get(0).(*repository.AuthCredentials), args.Error(1)
}
