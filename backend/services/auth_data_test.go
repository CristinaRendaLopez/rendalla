package services_test

import (
	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"golang.org/x/crypto/bcrypt"
)

var HashedSecretPassword, _ = bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)

var ValidStoredCredentials = repository.AuthCredentials{
	Username: "admin",
	Password: string(HashedSecretPassword),
}

var ValidLoginInput = dto.LoginRequest{
	Username: "admin",
	Password: "secret",
}

var InvalidUsernameInput = dto.LoginRequest{
	Username: "wrong",
	Password: "secret",
}

var InvalidPasswordInput = dto.LoginRequest{
	Username: "admin",
	Password: "wrongpass",
}

var GeneratedToken = "mocked.jwt.token"
