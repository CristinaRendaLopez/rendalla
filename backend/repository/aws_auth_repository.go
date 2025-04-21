package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
)

// AWSAuthRepository implements AuthRepository by retrieving admin credentials from AWS Secrets Manager.
// The secret must be a JSON object containing "username" and "password" fields.
type AWSAuthRepository struct{}

// NewAWSAuthRepository returns an appropriate AuthRepository implementation based on the environment.
//
// Parameters:
//   - env: the current runtime environment (e.g., "test", "production").
//
// Behavior:
//   - If env == "test", returns a FakeAuthRepository with fixed credentials for testing purposes.
//   - For all other environments, returns an AWSAuthRepository that fetches credentials from AWS Secrets Manager.
func NewAWSAuthRepository(env string) AuthRepository {
	if env == "test" {
		hashed := "$2a$10$xCIMIA6eYNX3lfarOEojqezWiDCMvxeJsA3kNAnNx7TX8d59sMjPy" // bcrypt("test")
		return &FakeAuthRepository{
			Credentials: AuthCredentials{
				Username: "test",
				Password: hashed,
			},
		}
	}

	return &AWSAuthRepository{}
}

// GetAuthCredentials retrieves admin authentication credentials from AWS Secrets Manager.
// The name of the secret is defined in the AUTH_SECRET_NAME environment variable,
// or defaults to "rendalla/auth_credentials" if not set.
// Returns:
//   - (*AuthCredentials, nil) on success
//   - (nil, errors.ErrInternalServer) if the secret cannot be fetched or parsed
func (a *AWSAuthRepository) GetAuthCredentials() (*AuthCredentials, error) {
	region := os.Getenv("AWS_REGION")
	secretName := os.Getenv("AUTH_SECRET_NAME")
	if secretName == "" {
		secretName = "rendalla/auth_credentials"
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation":   "get_auth_credentials",
			"secret_name": secretName,
		}).WithError(err).Error("Failed to retrieve secret from Secrets Manager")
		return nil, fmt.Errorf("retrieving secret %s from Secrets Manager: %w", secretName, errors.ErrInternalServer)
	}

	var credentials AuthCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"operation":   "get_auth_credentials",
			"secret_name": secretName,
		}).WithError(err).Error("Failed to parse secret JSON")
		return nil, fmt.Errorf("parsing secret %s: %w", secretName, errors.ErrInternalServer)
	}

	logrus.WithFields(logrus.Fields{
		"operation":   "get_auth_credentials",
		"secret_name": secretName,
	}).Info("Auth credentials retrieved successfully")

	return &credentials, nil
}
