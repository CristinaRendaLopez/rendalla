package services

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Could not load .env file, using default values")
	}
}

var jwtSecret = []byte(getJWTSecret())

func getAuthCredentials() (*AuthCredentials, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")), // Load AWS region from environment
	}))

	svc := secretsmanager.New(sess)

	secretName := os.Getenv("AUTH_SECRET_NAME")
	if secretName == "" {
		secretName = "rendalla/auth_credentials"
	}

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve secret from Secrets Manager")
		return nil, errors.New("failed to retrieve authentication credentials")
	}

	var credentials AuthCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse secret JSON")
		return nil, errors.New("invalid secret format")
	}

	return &credentials, nil
}

func AuthenticateUser(username, password string) (string, error) {
	creds, err := getAuthCredentials()
	if err != nil {
		return "", err
	}

	if username != creds.Username {
		logrus.WithField("username", username).Warn("Authentication failed: Invalid username")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password)); err != nil {
		logrus.WithField("username", username).Warn("Authentication failed: Incorrect password")
		return "", errors.New("invalid credentials")
	}

	expirationHours := getJWTExpirationHours()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		logrus.WithError(err).Error("Failed to generate JWT token")
		return "", err
	}

	logrus.WithField("username", username).Info("User authenticated successfully")
	return tokenString, nil
}

func getJWTExpirationHours() int {
	exp := os.Getenv("JWT_EXPIRATION_HOURS")
	if exp == "" {
		return 72
	}
	expInt, err := strconv.Atoi(exp)
	if err != nil {
		return 72
	}
	return expInt
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Warn("JWT_SECRET is not defined in .env, using default secret")
		return "default_secret"
	}
	return secret
}
