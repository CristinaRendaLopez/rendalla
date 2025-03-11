package main

import (
	"log"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var svc = bootstrap.InitAWS()

func main() {
	log.Println("Lambda initialized successfully")
	lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return handlers.GetSongHandler(svc, req)
	})
}
