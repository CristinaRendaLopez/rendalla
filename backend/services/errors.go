package services

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sirupsen/logrus"
)

// Translates DynamoDB errors into user-friendly messages
func handleDynamoError(err error) error {
	if err == nil {
		return nil
	}

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case dynamodb.ErrCodeConditionalCheckFailedException:
			logrus.WithError(awsErr).Warn("Operation not allowed: conditions not met")
			return errors.New("operation not allowed: conditions not met")
		case dynamodb.ErrCodeProvisionedThroughputExceededException:
			logrus.WithError(awsErr).Warn("Throughput limit exceeded in DynamoDB")
			return errors.New("throughput limit exceeded in DynamoDB")
		case dynamodb.ErrCodeResourceNotFoundException:
			logrus.WithError(awsErr).Error("DynamoDB resource not found")
			return errors.New("resource not found in DynamoDB")
		case dynamodb.ErrCodeInternalServerError:
			logrus.WithError(awsErr).Error("Internal error in DynamoDB")
			return errors.New("internal error in DynamoDB")
		default:
			logrus.WithError(awsErr).Error("Unhandled DynamoDB error")
			return errors.New("unknown error in the database")
		}
	}

	logrus.WithError(err).Error("Generic DynamoDB error")
	return err
}
