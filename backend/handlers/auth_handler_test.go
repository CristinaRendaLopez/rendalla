package handlers_test

import (
	"net/http"
	"strings"
	"testing"

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

func TestLoginHandler_Success(t *testing.T) {
	handler, mockService := setupAuthHandlerTest()

	validLogin := `{"username": "admin", "password": "securepass"}`
	mockService.On("AuthenticateUser", "admin", "securepass").Return("mocked.jwt.token", nil)

	c, w := utils.CreateTestContext(http.MethodPost, "/auth/login", strings.NewReader(validLogin))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LoginHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "mocked.jwt.token")
	mockService.AssertExpectations(t)
}

func TestLoginHandler_MissingFields(t *testing.T) {
	handler, mockService := setupAuthHandlerTest()

	invalidLogin := `{"username": "", "password": ""}`

	c, w := utils.CreateTestContext(http.MethodPost, "/auth/login", strings.NewReader(invalidLogin))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LoginHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request data")
	mockService.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	handler, mockService := setupAuthHandlerTest()

	mockService.On("AuthenticateUser", "admin", "wrongpass").Return("", utils.ErrInvalidCredentials)

	invalidLogin := `{"username": "admin", "password": "wrongpass"}`
	c, w := utils.CreateTestContext(http.MethodPost, "/auth/login", strings.NewReader(invalidLogin))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.LoginHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")
	mockService.AssertExpectations(t)
}

func TestMeHandler_Success(t *testing.T) {
	handler, _ := setupAuthHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/auth/me", nil)
	c.Set("username", "admin")

	handler.MeHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "admin")
	assert.Contains(t, w.Body.String(), "role")
}

func TestMeHandler_Unauthorized_InvalidType(t *testing.T) {
	handler, _ := setupAuthHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/auth/me", nil)
	c.Set("username", 12345)

	handler.MeHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}

func TestMeHandler_Unauthorized_MissingUsername(t *testing.T) {
	handler, _ := setupAuthHandlerTest()

	c, w := utils.CreateTestContext(http.MethodGet, "/auth/me", nil)

	handler.MeHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Unauthorized")
}
