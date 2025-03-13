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

// Struct to hold credentials from AWS Secrets Manager
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

// Retrieve authentication credentials from AWS Secrets Manager
func getAuthCredentials() (*AuthCredentials, error) {
	// Create AWS session
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")), // Load AWS region from environment
	}))

	// Create Secrets Manager client
	svc := secretsmanager.New(sess)

	// Retrieve secret name from environment variable
	secretName := os.Getenv("AUTH_SECRET_NAME")
	if secretName == "" {
		secretName = "rendalla/auth_credentials"
	}

	// Get secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}
	result, err := svc.GetSecretValue(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve secret from Secrets Manager")
		return nil, errors.New("failed to retrieve authentication credentials")
	}

	// Parse JSON secret
	var credentials AuthCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse secret JSON")
		return nil, errors.New("invalid secret format")
	}

	return &credentials, nil
}

// AuthenticateUser checks if the provided credentials are valid
func AuthenticateUser(username, password string) (string, error) {
	creds, err := getAuthCredentials()
	if err != nil {
		return "", err
	}

	// Validate username
	if username != creds.Username {
		logrus.WithField("username", username).Warn("Authentication failed: Invalid username")
		return "", errors.New("invalid credentials")
	}

	// Validate password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password)); err != nil {
		logrus.WithField("username", username).Warn("Authentication failed: Incorrect password")
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
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

// Retrieve JWT expiration time from environment variables
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

// Retrieve JWT secret key from environment variables
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Warn("JWT_SECRET is not defined in .env, using default secret")
		return "default_secret"
	}
	return secret
}
