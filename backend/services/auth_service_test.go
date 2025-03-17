package services_test

import (
	"errors"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticateUser_Success(t *testing.T) {
	mockDB := new(services.MockDB)

	username := "admin"
	password := "password123"
	expectedToken := "mocked-jwt-token"

	mockDB.On("AuthenticateUser", username, password).Return(expectedToken, nil)

	token, err := mockDB.AuthenticateUser(username, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	mockDB := new(services.MockDB)

	username := "admin"
	wrongPassword := "wrongpassword"

	mockDB.On("AuthenticateUser", username, wrongPassword).Return("", errors.New("invalid credentials"))

	token, err := mockDB.AuthenticateUser(username, wrongPassword)

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)
}

func TestGetAuthCredentials_Success(t *testing.T) {
	mockDB := new(services.MockDB)

	expectedCredentials := &services.AuthCredentials{
		Username: "admin",
		Password: "$2a$10$hashedpassword",
	}

	mockDB.On("GetAuthCredentials").Return(expectedCredentials, nil)

	creds, err := mockDB.GetAuthCredentials()

	assert.NoError(t, err)
	assert.Equal(t, expectedCredentials, creds)
}

func TestGetAuthCredentials_Failure(t *testing.T) {
	mockDB := new(services.MockDB)

	mockDB.On("GetAuthCredentials").Return(nil, errors.New("failed to retrieve credentials"))

	creds, err := mockDB.GetAuthCredentials()

	assert.Error(t, err)
	assert.Nil(t, creds)
	assert.Equal(t, "failed to retrieve credentials", err.Error())
}
