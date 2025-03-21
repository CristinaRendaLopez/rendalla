package repository

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRepository interface {
	GetAuthCredentials() (*AuthCredentials, error)
}
