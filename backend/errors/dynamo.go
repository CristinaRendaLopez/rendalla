package errors

import (
	stdErrors "errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

func HandleDynamoError(err error) error {
	if err == nil {
		return nil
	}

	if stdErrors.Is(err, dynamo.ErrNotFound) {
		logrus.WithError(err).Warn("Item not found in DynamoDB")
		return ErrResourceNotFound
	}

	if awsErr, ok := err.(awserr.Error); ok {
		switch awsErr.Code() {
		case dynamodb.ErrCodeConditionalCheckFailedException:
			logrus.WithError(awsErr).Warn("Operation not allowed: conditions not met")
			return ErrOperationNotAllowed
		case dynamodb.ErrCodeProvisionedThroughputExceededException:
			logrus.WithError(awsErr).Warn("Throughput limit exceeded in DynamoDB")
			return ErrThroughputExceeded
		case dynamodb.ErrCodeResourceNotFoundException:
			logrus.WithError(awsErr).Warn("DynamoDB resource not found")
			return ErrResourceNotFound
		case dynamodb.ErrCodeInternalServerError:
			logrus.WithError(awsErr).Error("Internal error in DynamoDB")
			return ErrInternalServer
		default:
			logrus.WithError(awsErr).Error("Unhandled DynamoDB error")
			return ErrInternalServer
		}
	}

	logrus.WithError(err).Error("Generic error")
	return ErrInternalServer
}

func IsDynamoNotFoundError(err error) bool {
	return stdErrors.Is(err, ErrResourceNotFound)
}
