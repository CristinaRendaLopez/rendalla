package services_test

import (
	"errors"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthServiceTest() (*services.AuthService, *mocks.MockAuthRepository, *mocks.MockClock, *mocks.MockTokenGenerator) {
	authRepo := new(mocks.MockAuthRepository)
	clock := new(mocks.MockClock)
	tokenGen := new(mocks.MockTokenGenerator)
	service := services.NewAuthService(authRepo, clock, tokenGen)
	return service, authRepo, clock, tokenGen
}

func TestAuthenticateUser_Success(t *testing.T) {
	service, authRepo, clock, tokenGen := setupAuthServiceTest()

	hashedPwd, _ := bcryptHash("secret")
	creds := &repository.AuthCredentials{Username: "admin", Password: hashedPwd}

	authRepo.On("GetAuthCredentials").Return(creds, nil)
	clock.On("NowUnix").Return(1000)
	expectedClaims := jwt.MapClaims{"username": "admin", "exp": int64(1000 + 72*3600)}
	tokenGen.On("GenerateToken", expectedClaims).Return("jwt-token", nil)

	token, err := service.AuthenticateUser("admin", "secret")

	assert.NoError(t, err)
	assert.Equal(t, "jwt-token", token)

	authRepo.AssertExpectations(t)
	clock.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidUsername(t *testing.T) {
	service, authRepo, _, _ := setupAuthServiceTest()

	creds := &repository.AuthCredentials{Username: "admin", Password: "xxx"}
	authRepo.On("GetAuthCredentials").Return(creds, nil)

	token, err := service.AuthenticateUser("wrong", "secret")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)

	authRepo.AssertExpectations(t)
}

func TestAuthenticateUser_InvalidPassword(t *testing.T) {
	service, authRepo, _, _ := setupAuthServiceTest()

	hashedPwd, _ := bcryptHash("realpass")
	creds := &repository.AuthCredentials{Username: "admin", Password: hashedPwd}

	authRepo.On("GetAuthCredentials").Return(creds, nil)

	token, err := service.AuthenticateUser("admin", "wrongpass")

	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
	assert.Empty(t, token)

	authRepo.AssertExpectations(t)
}

func TestAuthenticateUser_RepoError(t *testing.T) {
	service, authRepo, _, _ := setupAuthServiceTest()

	authRepo.On("GetAuthCredentials").Return(nil, errors.New("db error"))

	token, err := service.AuthenticateUser("admin", "secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	assert.Empty(t, token)

	authRepo.AssertExpectations(t)
}

func TestAuthenticateUser_TokenGenerationError(t *testing.T) {
	service, authRepo, clock, tokenGen := setupAuthServiceTest()

	hashedPwd, _ := bcryptHash("secret")
	creds := &repository.AuthCredentials{Username: "admin", Password: hashedPwd}

	authRepo.On("GetAuthCredentials").Return(creds, nil)
	clock.On("NowUnix").Return(1000)
	claims := jwt.MapClaims{"username": "admin", "exp": int64(1000 + 72*3600)}
	tokenGen.On("GenerateToken", claims).Return("", errors.New("token error"))

	token, err := service.AuthenticateUser("admin", "secret")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token error")
	assert.Empty(t, token)

	authRepo.AssertExpectations(t)
	clock.AssertExpectations(t)
	tokenGen.AssertExpectations(t)
}

func TestGetAuthCredentials_Success(t *testing.T) {
	service, authRepo, _, _ := setupAuthServiceTest()

	expected := &repository.AuthCredentials{Username: "admin", Password: "pass"}
	authRepo.On("GetAuthCredentials").Return(expected, nil)

	result, err := service.GetAuthCredentials()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	authRepo.AssertExpectations(t)
}

func bcryptHash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(hash), err
}
