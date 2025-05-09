package handlers

import (
	stdErrors "errors"
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthHandler handles HTTP requests related to admin authentication.
// It delegates logic to the AuthServiceInterface.
type AuthHandler struct {
	authService services.AuthServiceInterface
}

// NewAuthHandler returns a new instance of AuthHandler.
func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginHandler handles POST /auth/login.
// Validates credentials and returns a signed JWT token upon successful authentication.
func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Invalid request data")
		return
	}

	token, err := h.authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		message := "Authentication failed"
		switch {
		case stdErrors.Is(err, errors.ErrInvalidCredentials):
			message = "Invalid credentials"
		case stdErrors.Is(err, errors.ErrTokenGenerationFailed):
			message = "Failed to generate token"
		case stdErrors.Is(err, errors.ErrInternalServer):
			message = "Server error during authentication"
		}
		errors.HandleAPIError(c, err, message)
		return
	}

	logrus.WithField("username", req.Username).Info("User authenticated successfully")

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// MeHandler handles GET /auth/me.
// Returns the username of the authenticated user and a hardcoded role.
func (h *AuthHandler) MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		errors.HandleAPIError(c, errors.ErrUnauthorized, "Unauthorized")
		return
	}

	strUsername, ok := username.(string)
	if !ok || utils.IsEmptyString(strUsername) {
		errors.HandleAPIError(c, errors.ErrUnauthorized, "Unauthorized")
		return
	}

	logrus.WithField("username", strUsername).Info("User details retrieved successfully")

	c.JSON(http.StatusOK, gin.H{"username": strUsername, "role": "admin"})
}
