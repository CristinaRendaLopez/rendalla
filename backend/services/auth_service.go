package services

import (
	"errors"
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.AuthRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.AuthRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
	creds, err := s.repo.GetAuthCredentials()
	if err != nil {
		return "", err
	}

	if username != creds.Username {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(72 * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      expirationTime,
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetAuthCredentials() (*repository.AuthCredentials, error) {
	return s.repo.GetAuthCredentials()
}
