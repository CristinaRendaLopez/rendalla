package repository

import (
	"encoding/json"
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

// NewAWSAuthRepository returns a new instance of AWSAuthRepository.
func NewAWSAuthRepository() *AWSAuthRepository {
	return &AWSAuthRepository{}
}

// GetAuthCredentials retrieves admin authentication credentials from AWS Secrets Manager.
// The name of the secret is defined in the AUTH_SECRET_NAME environment variable,
// or defaults to "rendalla/auth_credentials" if not set.
// Returns:
//   - (*AuthCredentials, nil) on success
//   - (nil, errors.ErrInternalServer) if the secret cannot be fetched or parsed
func (a *AWSAuthRepository) GetAuthCredentials() (*AuthCredentials, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
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
		return nil, errors.ErrInternalServer
	}

	var credentials AuthCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse secret JSON")
		return nil, errors.ErrInternalServer
	}

	return &credentials, nil
}
