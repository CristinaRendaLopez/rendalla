package services_test

import (
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthServiceTest() (*services.AuthService, *mocks.MockAuthRepository, *mocks.MockTimeProvider, *mocks.MockTokenGenerator) {
	authRepo := new(mocks.MockAuthRepository)
	timeProvider := new(mocks.MockTimeProvider)
	tokenGen := new(mocks.MockTokenGenerator)
	service := services.NewAuthService(authRepo, timeProvider, tokenGen)
	return service, authRepo, timeProvider, tokenGen
}

func TestAuthenticateUser(t *testing.T) {
	const nowUnix int64 = 1000

	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte("realpass"), bcrypt.DefaultCost)
	customCreds := repository.AuthCredentials{
		Username: "admin",
		Password: string(hashedPwd),
	}

	tests := []struct {
		name           string
		input          dto.LoginRequest
		creds          *repository.AuthCredentials
		mockNow        int64
		mockToken      string
		mockTokenError error
		mockCredsError error
		expectToken    string
		expectError    error
	}{
		{
			name:        "valid credentials",
			input:       ValidLoginInput,
			creds:       &ValidStoredCredentials,
			mockNow:     nowUnix,
			mockToken:   GeneratedToken,
			expectToken: GeneratedToken,
		},
		{
			name:        "invalid username",
			input:       InvalidUsernameInput,
			creds:       &ValidStoredCredentials,
			expectError: errors.ErrInvalidCredentials,
		},
		{
			name:        "invalid password",
			input:       InvalidPasswordInput,
			creds:       &customCreds,
			expectError: errors.ErrInvalidCredentials,
		},
		{
			name:           "repository error",
			input:          ValidLoginInput,
			mockCredsError: errors.ErrInternalServer,
			expectError:    errors.ErrInternalServer,
		},
		{
			name:           "token generation fails",
			input:          ValidLoginInput,
			creds:          &ValidStoredCredentials,
			mockNow:        nowUnix,
			mockTokenError: errors.ErrTokenGenerationFailed,
			expectError:    errors.ErrTokenGenerationFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, authRepo, clock, tokenGen := setupAuthServiceTest()

			if tt.mockCredsError != nil {
				authRepo.On("GetAuthCredentials").Return(nil, tt.mockCredsError)
			} else {
				authRepo.On("GetAuthCredentials").Return(tt.creds, nil)
			}

			if tt.mockNow != 0 {
				clock.On("NowUnix").Return(tt.mockNow)
			}

			if tt.mockToken != "" || tt.mockTokenError != nil {
				claims := jwt.MapClaims{
					"username": tt.input.Username,
					"exp":      tt.mockNow + 72*3600,
				}
				tokenGen.On("GenerateToken", claims).Return(tt.mockToken, tt.mockTokenError)
			}

			token, err := service.AuthenticateUser(tt.input.Username, tt.input.Password)

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectToken, token)
			}

			authRepo.AssertExpectations(t)
			clock.AssertExpectations(t)
			tokenGen.AssertExpectations(t)
		})
	}
}

func TestGetAuthCredentials(t *testing.T) {
	tests := []struct {
		name         string
		mockCreds    *repository.AuthCredentials
		mockError    error
		expectError  error
		expectedCred *repository.AuthCredentials
	}{
		{
			name:         "success",
			mockCreds:    &ValidStoredCredentials,
			mockError:    nil,
			expectError:  nil,
			expectedCred: &ValidStoredCredentials,
		},
		{
			name:        "repository error",
			mockCreds:   nil,
			mockError:   errors.ErrInternalServer,
			expectError: errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, authRepo, _, _ := setupAuthServiceTest()

			authRepo.On("GetAuthCredentials").Return(tt.mockCreds, tt.mockError)

			result, err := service.GetAuthCredentials()

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCred, result)
			}

			authRepo.AssertExpectations(t)
		})
	}
}
