package handlers

import (
	"errors"
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authService services.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleAPIError(c, utils.ErrValidationFailed, "Invalid request data")
		return
	}

	token, err := h.authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		message := "Authentication failed"
		switch {
		case errors.Is(err, utils.ErrInvalidCredentials):
			message = "Invalid credentials"
		case errors.Is(err, utils.ErrTokenGenerationFailed):
			message = "Failed to generate token"
		case errors.Is(err, utils.ErrInternalServer):
			message = "Server error during authentication"
		}
		utils.HandleAPIError(c, err, message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("User authenticated successfully")

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.HandleAPIError(c, utils.ErrUnauthorized, "Unauthorized")
		return
	}

	strUsername, ok := username.(string)
	if !ok || utils.IsEmptyString(strUsername) {
		utils.HandleAPIError(c, utils.ErrUnauthorized, "Unauthorized")
		return
	}

	logrus.WithFields(logrus.Fields{
		"username": strUsername,
	}).Info("User details retrieved successfully")

	c.JSON(http.StatusOK, gin.H{"username": strUsername, "role": "admin"})
}
