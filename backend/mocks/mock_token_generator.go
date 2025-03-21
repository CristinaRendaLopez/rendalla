package mocks

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) GenerateToken(claims jwt.MapClaims) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(1)
}
