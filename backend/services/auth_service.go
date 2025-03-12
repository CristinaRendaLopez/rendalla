package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Could not load .env file, using default values")
	}
}

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Warn("JWT_SECRET is not defined in .env, using default secret")
		return "default_secret"
	}
	return secret
}

func AuthenticateUser(username, password string) (string, error) {
	if username != "admin" || password != "password123" {
		logrus.WithField("username", username).Warn("Authentication failed: Invalid credentials")
		return "", errors.New("invalid credentials")
	}

	expirationHours := getJWTExpirationHours()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token")
		return "", err
	}

	logrus.WithField("username", username).Info("User authenticated successfully")
	return tokenString, nil
}

func getJWTExpirationHours() int {
	exp := os.Getenv("JWT_EXPIRATION_HOURS")
	if exp == "" {
		return 72
	}
	return 72
}
