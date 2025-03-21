package mocks

import "github.com/stretchr/testify/mock"

type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() string {
	args := m.Called()
	return args.String(0)
}
