package errors

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HandleAPIError sends a standardized JSON error response to the client.
// It also logs the error with appropriate severity based on the HTTP status code.
func HandleAPIError(c *gin.Context, err error, message string) {
	if err == nil {
		return
	}

	statusCode := MapErrorToStatus(err)
	errorCode := MapErrorToCode(err)
	requestID := c.GetHeader("X-Request-ID")

	entry := logrus.WithFields(logrus.Fields{
		"error":      err,
		"status":     statusCode,
		"error_code": errorCode,
		"request_id": requestID,
	})

	if statusCode >= 500 {
		entry.Error(message)
	} else {
		entry.Warn(message)
	}

	c.JSON(statusCode, gin.H{
		"error": gin.H{
			"message":    message,
			"error_code": errorCode,
			"status":     statusCode,
			"request_id": requestID,
		},
	})
}

// CreateErrorResponse builds an API Gateway-compatible error response.
// It maps internal errors to meaningful HTTP status codes and serializes the error message.
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
