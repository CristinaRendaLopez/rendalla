package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrResourceNotFound      = errors.New("requested resource not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
	ErrOperationNotAllowed   = errors.New("operation not allowed: conditions not met")
	ErrThroughputExceeded    = errors.New("throughput limit exceeded, please try again later")
)

func HandleDynamoError(err error) error {
	if err == nil {
		return nil
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
			logrus.WithError(awsErr).Error("DynamoDB resource not found")
			return ErrResourceNotFound
		case dynamodb.ErrCodeInternalServerError:
			logrus.WithError(awsErr).Error("Internal error in DynamoDB")
			return errors.New("internal server error, please try again later")
		default:
			logrus.WithError(awsErr).Error("Unhandled DynamoDB error")
			return errors.New("unexpected database error")
		}
	}

	logrus.WithError(err).Error("Generic error")
	return errors.New("an unexpected error occurred")
}

func IsDynamoNotFoundError(err error) bool {
	return errors.Is(err, ErrResourceNotFound)
}

func HandleAPIError(c *gin.Context, err error, message string) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError

	switch {
	case IsDynamoNotFoundError(err):
		statusCode = http.StatusNotFound
	case errors.Is(err, ErrResourceNotFound):
		statusCode = http.StatusNotFound
	case errors.Is(err, ErrInvalidCredentials):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, ErrTokenGenerationFailed):
		statusCode = http.StatusInternalServerError
	case errors.Is(err, ErrOperationNotAllowed):
		statusCode = http.StatusForbidden
	case errors.Is(err, ErrThroughputExceeded):
		statusCode = http.StatusTooManyRequests
	default:
		statusCode = http.StatusInternalServerError
	}

	logrus.WithError(err).Error(message)
	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"message": message,
			"detail":  err.Error(),
			"code":    statusCode,
		},
	})
}

func CreateErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	statusCode := http.StatusInternalServerError
	message := "An unexpected error occurred"

	switch {
	case IsDynamoNotFoundError(err):
		statusCode = http.StatusNotFound
		message = "Resource not found"
	case errors.Is(err, ErrOperationNotAllowed):
		statusCode = http.StatusForbidden
		message = "Operation not allowed"
	case errors.Is(err, ErrThroughputExceeded):
		statusCode = http.StatusTooManyRequests
		message = "Too many requests, try again later"
	default:
		logrus.WithError(err).Error("Unhandled error in API Gateway")
	}

	body, _ := json.Marshal(map[string]string{"error": message})
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
