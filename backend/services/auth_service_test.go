package services_test

import (
	"errors"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthenticateUser_Success(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedCredentials := &repository.AuthCredentials{
		Username: "admin",
		Password: string(hashedPassword),
	}

	mockAuthRepo.On("GetAuthCredentials").Return(expectedCredentials, nil)

	token, err := service.AuthenticateUser("admin", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidPassword(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedCredentials := &repository.AuthCredentials{
		Username: "admin",
		Password: string(hashedPassword),
	}

	mockAuthRepo.On("GetAuthCredentials").Return(expectedCredentials, nil)

	token, err := service.AuthenticateUser("admin", "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidUsername(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedCredentials := &repository.AuthCredentials{
		Username: "admin",
		Password: string(hashedPassword),
	}

	mockAuthRepo.On("GetAuthCredentials").Return(expectedCredentials, nil)

	token, err := service.AuthenticateUser("wronguser", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticateUser_TokenGenerationError(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	expectedCredentials := &repository.AuthCredentials{
		Username: "admin",
		Password: string(hashedPassword),
	}

	mockAuthRepo.On("GetAuthCredentials").Return(expectedCredentials, nil)

	token, err := service.AuthenticateUser("admin", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.True(t, errors.Is(err, utils.ErrTokenGenerationFailed))
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthenticateUser_UserNotFound(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	mockAuthRepo.On("GetAuthCredentials").Return(nil, errors.New("user not found"))

	token, err := service.AuthenticateUser("admin", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "user not found", err.Error())
	mockAuthRepo.AssertExpectations(t)
}

func TestGetAuthCredentials_Success(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	expectedCredentials := &repository.AuthCredentials{
		Username: "admin",
		Password: "$2a$10$hashedpassword",
	}

	mockAuthRepo.On("GetAuthCredentials").Return(expectedCredentials, nil)

	creds, err := service.GetAuthCredentials()

	assert.NoError(t, err)
	assert.Equal(t, expectedCredentials, creds)
	mockAuthRepo.AssertExpectations(t)
}

func TestGetAuthCredentials_Failure(t *testing.T) {
	mockAuthRepo := new(mocks.MockAuthRepository)
	service := services.NewAuthService(mockAuthRepo, "mock-secret")

	mockAuthRepo.On("GetAuthCredentials").Return(nil, errors.New("failed to retrieve credentials"))

	creds, err := service.GetAuthCredentials()

	assert.Error(t, err)
	assert.Nil(t, creds)
	assert.Equal(t, "failed to retrieve credentials", err.Error())
	mockAuthRepo.AssertExpectations(t)
}
