package repository

// AuthCredentials holds the login credentials for the administrator.
type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthRepository defines how administrator credentials are retrieved for authentication.
type AuthRepository interface {

	// GetAuthCredentials returns the stored admin username and hashed password.
	GetAuthCredentials() (*AuthCredentials, error)
}
