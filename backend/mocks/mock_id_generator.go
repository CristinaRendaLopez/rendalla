package mocks

import "github.com/stretchr/testify/mock"

type MockIDGenerator struct {
	mock.Mock
}

func (m *MockIDGenerator) NewID() string {
	args := m.Called()
	return args.String(0)
}
