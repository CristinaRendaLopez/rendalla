package handlers

import (
	"encoding/json"
	"log"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func GetSongHandler(svc *dynamodb.DynamoDB, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := req.PathParameters["id"]
	if id == "" {
		log.Println("Missing song ID in request")
		return createErrorResponse(400, "Missing song ID")
	}

	log.Printf("Fetching song with ID: %s", id)
	song, err := services.GetSongByID(svc, id)
	if err != nil {
		log.Printf("Error fetching song from DynamoDB: %v", err)
		return createErrorResponse(500, "Error fetching song")
	}

	if song == nil {
		log.Printf("Song not found: %s", id)
		return createErrorResponse(404, "Song not found")
	}

	body, err := json.Marshal(song)
	if err != nil {
		log.Printf("Error generating JSON response: %v", err)
		return createErrorResponse(500, "Error generating JSON response")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
