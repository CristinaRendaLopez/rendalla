package services

import (
	"fmt"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// AuthService provides authentication logic for the admin user.
// It verifies credentials and generates JWT tokens for access.
type AuthService struct {
	repo           repository.AuthRepository
	timeProvider   utils.TimeProvider
	tokenGenerator utils.TokenGenerator
}

// Ensure AuthService implements AuthServiceInterface.
var _ AuthServiceInterface = (*AuthService)(nil)

// NewAuthService returns a new instance of AuthService with its required dependencies.
func NewAuthService(
	repo repository.AuthRepository,
	timeProvider utils.TimeProvider,
	tokenGenerator utils.TokenGenerator,
) *AuthService {
	return &AuthService{
		repo:           repo,
		timeProvider:   timeProvider,
		tokenGenerator: tokenGenerator,
	}
}

// AuthenticateUser verifies the given username and password against stored credentials.
// If valid, it generates a signed JWT token valid for 72 hours.
// Returns:
//   - the signed JWT token string on success
//   - errors.ErrInvalidCredentials if the credentials are incorrect
//   - errors.ErrTokenGenerationFailed if token signing fails
//   - other repository errors if credential retrieval fails
func (s *AuthService) AuthenticateUser(username, password string) (string, error) {
	creds, err := s.repo.GetAuthCredentials()
	if err != nil {
		return "", fmt.Errorf("retrieving auth credentials: %w", err)
	}

	if username != creds.Username {
		return "", errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password)); err != nil {
		return "", errors.ErrInvalidCredentials
	}

	exp := s.timeProvider.NowUnix() + 72*3600
	claims := jwt.MapClaims{
		"username": username,
		"exp":      exp,
	}

	token, err := s.tokenGenerator.GenerateToken(claims)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token")
		return "", fmt.Errorf("generating JWT token: %w", errors.ErrTokenGenerationFailed)
	}

	return token, nil
}

// GetAuthCredentials retrieves the current admin credentials from the repository.
// Returns:
//   - the stored credentials on success
//   - errors.ErrInternalServer if retrieval fails
func (s *AuthService) GetAuthCredentials() (*repository.AuthCredentials, error) {
	creds, err := s.repo.GetAuthCredentials()
	if err != nil {
		return nil, fmt.Errorf("retrieving stored auth credentials: %w", err)
	}
	return creds, nil
}
