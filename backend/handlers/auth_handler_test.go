package handlers_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
)

func setupAuthHandlerTest() (*handlers.AuthHandler, *mocks.MockAuthService) {
	mockService := new(mocks.MockAuthService)
	handler := handlers.NewAuthHandler(mockService)
	return handler, mockService
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		setupMock    bool
		mockToken    string
		mockError    error
		expectedCode int
		expectToken  bool
	}{
		{
			name:         "success",
			body:         ValidLoginJSON,
			setupMock:    true,
			mockToken:    "valid.jwt.token",
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectToken:  true,
		},
		{
			name:         "invalid JSON",
			body:         InvalidLoginJSON,
			setupMock:    false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing required field",
			body:         MissingPasswordJSON,
			setupMock:    false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid credentials",
			body:         ValidLoginJSON,
			setupMock:    true,
			mockError:    errors.ErrInvalidCredentials,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "token generation fails",
			body:         ValidLoginJSON,
			setupMock:    true,
			mockError:    errors.ErrTokenGenerationFailed,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "internal server error",
			body:         ValidLoginJSON,
			setupMock:    true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupAuthHandlerTest()

			if tt.setupMock {
				mockService.On("AuthenticateUser", "admin", "secret123").
					Return(tt.mockToken, tt.mockError)
			}

			c, w := utils.CreateTestContext(http.MethodPost, "/auth/login", strings.NewReader(tt.body))
			handler.LoginHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectToken {
				var res dto.AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockToken, res.Token)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestMeHandler(t *testing.T) {
	tests := []struct {
		name         string
		username     interface{}
		expectedCode int
		expectedRole string
	}{
		{
			name:         "valid username in context",
			username:     "admin",
			expectedCode: http.StatusOK,
			expectedRole: "admin",
		},
		{
			name:         "missing username in context",
			username:     nil,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "username not a string",
			username:     1234,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "username is empty string",
			username:     "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, _ := setupAuthHandlerTest()

			c, w := utils.CreateTestContext(http.MethodGet, "/auth/me", nil)
			if tt.username != nil {
				c.Set("username", tt.username)
			}

			handler.MeHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var res dto.MeResponse
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)
				assert.Equal(t, "admin", res.Username)
				assert.Equal(t, tt.expectedRole, res.Role)
			}
		})
	}
}
