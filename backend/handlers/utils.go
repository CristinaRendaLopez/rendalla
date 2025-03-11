package handlers

import "github.com/aws/aws-lambda-go/events"

func CreateErrorResponse(statusCode int, message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       `{"error": "` + message + `"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
