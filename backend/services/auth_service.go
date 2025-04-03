package services

import (
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	AuthenticateUser(username, password string) (string, error)
	GetAuthCredentials() (*repository.AuthCredentials, error)
}

type Clock interface {
	NowUnix() int64
}

type TokenGenerator interface {
	GenerateToken(claims jwt.MapClaims) (string, error)
}

type AuthService struct {
	repo           repository.AuthRepository
	clock          Clock
	tokenGenerator TokenGenerator
}

var _ AuthServiceInterface = (*AuthService)(nil)

func NewAuthService(repo repository.AuthRepository, clock Clock, tokenGen TokenGenerator) *AuthService {
	return &AuthService{
		repo:           repo,
		clock:          clock,
		tokenGenerator: tokenGen,
	}
}

func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
	creds, err := s.repo.GetAuthCredentials()
	if err != nil {
		return "", err
	}

	if username != creds.Username {
		return "", utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password)); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	exp := s.clock.NowUnix() + 72*3600
	claims := jwt.MapClaims{
		"username": username,
		"exp":      exp,
	}

	token, err := s.tokenGenerator.GenerateToken(claims)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token")
		return "", utils.ErrTokenGenerationFailed
	}

	return token, nil
}

func (s *AuthService) GetAuthCredentials() (*repository.AuthCredentials, error) {
	return s.repo.GetAuthCredentials()
}
