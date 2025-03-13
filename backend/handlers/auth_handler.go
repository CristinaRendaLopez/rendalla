package handlers

import (
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
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
		utils.HandleAPIError(c, err, "Invalid login request")
		return
	}

	if req.Username == "" || req.Password == "" {
		logrus.Warn("Login attempt with missing credentials")
		utils.HandleAPIError(c, nil, "Username and password are required")
		return
	}

	token, err := services.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		utils.HandleAPIError(c, err, "Invalid credentials")
		return
	}

	logrus.WithField("username", req.Username).Info("User authenticated successfully")
	c.JSON(200, gin.H{"token": token})
}

func MeHandler(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		utils.HandleAPIError(c, nil, "Unauthorized")
		return
	}

	logrus.WithField("username", username).Info("User details retrieved successfully")
	c.JSON(200, gin.H{"username": username, "role": "admin"})
}
