package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.WithError(err).Warn("Invalid login request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
		return
	}

	if req.Username == "" || req.Password == "" {
		logrus.Warn("Login attempt with missing credentials")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	token, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		logrus.WithField("username", req.Username).Warn("Failed login attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	logrus.WithField("username", req.Username).Info("User authenticated successfully")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		logrus.Warn("Unauthorized access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	logrus.WithField("username", username).Info("User details retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"username": username, "role": "admin"})
}
