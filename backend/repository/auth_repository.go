package repository

// AuthCredentials holds the login credentials for the administrator.
type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthRepository defines how administrator credentials are retrieved for authentication.
type AuthRepository interface {

	// GetAuthCredentials returns the stored admin username and hashed password.
	// Returns:
	//   - (*AuthCredentials, nil) on success
	//   - (nil, errors.ErrInternalServer) if retrieval or parsing fails
	GetAuthCredentials() (*AuthCredentials, error)
}
