package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
	ErrBadRequest            = errors.New("bad request")
	ErrValidationFailed      = errors.New("validation failed")
	ErrResourceNotFound      = errors.New("resource not found")
	ErrOperationNotAllowed   = errors.New("operation not allowed")
	ErrThroughputExceeded    = errors.New("throughput limit exceeded")
	ErrInternalServer        = errors.New("internal server error")
	ErrUnauthorized          = errors.New("unauthorized access")
)

func HandleDynamoError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, dynamo.ErrNotFound) {
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
			logrus.WithError(awsErr).Error("DynamoDB resource not found")
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
	return errors.Is(err, ErrResourceNotFound)
}

func HandleAPIError(c *gin.Context, err error, message string) {
	if err == nil {
		return
	}

	statusCode := MapErrorToStatus(err)
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
	statusCode := MapErrorToStatus(err)
	message := MapErrorToMessage(err)

	logrus.WithError(err).WithField("status", statusCode).Error("API Gateway Error")

	body, _ := json.Marshal(map[string]string{
		"error": message,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrResourceNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrOperationNotAllowed):
		return http.StatusForbidden
	case errors.Is(err, ErrThroughputExceeded):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func MapErrorToMessage(err error) string {
	switch {
	case errors.Is(err, ErrBadRequest):
		return "Bad request"
	case errors.Is(err, ErrValidationFailed):
		return "Validation failed"
	case errors.Is(err, ErrResourceNotFound):
		return "Resource not found"
	case errors.Is(err, ErrOperationNotAllowed):
		return "Operation not allowed"
	case errors.Is(err, ErrThroughputExceeded):
		return "Too many requests, try again later"
	case errors.Is(err, ErrInvalidCredentials):
		return "Invalid credentials"
	case errors.Is(err, ErrUnauthorized):
		return "Unauthorized access"
	case errors.Is(err, ErrTokenGenerationFailed):
		return "Failed to generate token"
	case errors.Is(err, ErrInternalServer):
		return "Internal server error"
	default:
		return "An unexpected error occurred"
	}
}
