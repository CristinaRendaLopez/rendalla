package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"), // Cambia a tu región
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

func HandleRequest() (*Song, error) {
	id := "some-id"

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("MyAppTable"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to find item %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("Could not find the item")
	}

	var song Song
	err = dynamodbattribute.UnmarshalMap(result.Item, &song)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item: %v", err)
	}

	// Devuelve la canción recuperada
	return &song, nil
}

func main() {
	lambda.Start(HandleRequest)
}
