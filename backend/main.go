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

func HandleRequest() (map[string]interface{}, error) {
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
		return nil, fmt.Errorf("unable to find item %v", err)
	}

	if result.Item == nil {
		return nil, fmt.Errorf("could not find the item")
	}

	var song Song

	err = dynamodbattribute.UnmarshalMap(result.Item, &song)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal item: %v", err)
	}

	response := map[string]interface{}{
		"statusCode": 200,
		"body":       song,
	}

	return response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
