package repository

import (
	"encoding/json"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
)

type AWSAuthRepository struct{}

func NewAWSAuthRepository() *AWSAuthRepository {
	return &AWSAuthRepository{}
}

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
		return nil, utils.ErrInternalServer
	}

	var credentials AuthCredentials
	err = json.Unmarshal([]byte(*result.SecretString), &credentials)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse secret JSON")
		return nil, utils.ErrInternalServer
	}

	return &credentials, nil
}
