package services

import "github.com/CristinaRendaLopez/rendalla-backend/repository"

// AuthServiceInterface defines authentication-related operations for the admin user.
type AuthServiceInterface interface {

	// AuthenticateUser verifies the provided username and password against stored credentials.
	// Returns:
	//   - a signed JWT token string on successful authentication
	//   - utils.ErrInvalidCredentials if authentication fails
	//   - utils.ErrInternalServer if token generation or credential retrieval fails
	AuthenticateUser(username, password string) (string, error)

	// GetAuthCredentials retrieves the stored admin credentials from the repository.
	// Returns:
	//   - (*repository.AuthCredentials, nil) on success
	//   - (nil, utils.ErrInternalServer) if retrieval fails
	GetAuthCredentials() (*repository.AuthCredentials, error)
}
