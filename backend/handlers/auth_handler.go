package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: *authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleAPIError(c, err, "Invalid request data")
		return
	}

	token, err := h.authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		utils.HandleAPIError(c, err, "Invalid credentials")
		return
	}

	logrus.WithField("username", req.Username).Info("User authenticated successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.HandleAPIError(c, nil, "Unauthorized")
		return
	}

	logrus.WithField("username", username).Info("User details retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"username": username, "role": "admin"})
}
