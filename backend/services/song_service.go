package services

import (
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var tableName = os.Getenv("DYNAMODB_TABLE")

func GetSongByID(svc *dynamodb.DynamoDB, id string) (*models.Song, error) {
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var song models.Song
	err = dynamodbattribute.UnmarshalMap(result.Item, &song)
	if err != nil {
		return nil, err
	}

	return &song, nil
}
