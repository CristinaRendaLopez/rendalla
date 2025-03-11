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
		log.Fatalf("Error starting AWS session: %v", err)
	}
	svc = dynamodb.New(sess)
}

type Song struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func HandleRequest() (events.APIGatewayProxyResponse, error) {
	id := "some-id"
	log.Printf("Fetching song with ID: %s", id)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("RendallaTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Printf("Error accessing DynamoDB: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error accessing DynamoDB: ` + err.Error() + `"}`,
		}, nil
	}

	if result.Item == nil {
		log.Println("No data found in DynamoDB")
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
		log.Printf("Error processing data: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error processing data"}`,
		}, nil
	}

	body, err := json.Marshal(song)
	if err != nil {
		log.Printf("Error generating JSON response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Error generating JSON response"}`,
		}, nil
	}

	log.Printf("Song found: %+v", song)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	log.Println("Lambda initialized successfully")
	lambda.Start(HandleRequest)
}
