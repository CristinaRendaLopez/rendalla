package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-north-1"),
	})
	if err != nil {
		log.Fatalf("unable to start session: %v", err)
	}
	svc = dynamodb.New(sess)
}

type Song struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func HandleRequest() (events.APIGatewayProxyResponse, error) {
	id := "some-id"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("RendallaTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error al acceder a DynamoDB: ` + err.Error() + `"}`,
		}, nil
	}

	if result.Item == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "[]",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	var song Song
	err = dynamodbattribute.UnmarshalMap(result.Item, &song)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error al procesar los datos"}`,
		}, nil
	}

	body, err := json.Marshal(song)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error al generar la respuesta JSON"}`,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
