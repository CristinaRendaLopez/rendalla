package mocks

import "github.com/stretchr/testify/mock"

type MockTimeProvider struct {
	mock.Mock
}

func (m *MockTimeProvider) Now() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockTimeProvider) NowUnix() int64 {
	args := m.Called()
	return int64(args.Int(0))
}
