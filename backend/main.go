package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rendalla/backend/bootstrap"
	"github.com/rendalla/backend/handlers"
)

var svc = bootstrap.InitAWS()

func main() {
	log.Println("Lambda initialized successfully")
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handlers.GetSongHandler(svc, req)
	})
}
