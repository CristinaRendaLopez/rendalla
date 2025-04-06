package integration_tests

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTestJWT(username string) (string, error) {
	return GenerateTestJWTWithSecret(username, os.Getenv("JWT_SECRET"))
}

func GenerateTestJWTWithSecret(username, secret string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
