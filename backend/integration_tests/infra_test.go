package integration_tests

import (
	"strings"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

func CreateTestTables(db *dynamo.DB) error {
	svc, ok := db.Client().(*dynamodb.DynamoDB)
	if !ok {
		logrus.Fatal("Failed to assert Dynamo client as *dynamodb.DynamoDB")
	}

	for _, table := range []string{bootstrap.SongTableName, bootstrap.DocumentTableName} {
		_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(table)})
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete table %s", table)
		}
	}

	if err := createSongsTable(svc); err != nil {
		return err
	}

	if err := createDocumentsTable(svc); err != nil {
		return err
	}

	return nil
}

func createSongsTable(svc *dynamodb.DynamoDB) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(bootstrap.SongTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String("HASH")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := svc.CreateTable(input)
	if err != nil && !strings.Contains(err.Error(), dynamodb.ErrCodeResourceInUseException) {
		logrus.WithError(err).Error("Failed to create SongsTable")
		return err
	}

	return waitForTableToBeActive(svc, bootstrap.SongTableName)
}

func createDocumentsTable(svc *dynamodb.DynamoDB) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(bootstrap.DocumentTableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("song_id"), KeyType: aws.String("HASH")},
			{AttributeName: aws.String("id"), KeyType: aws.String("RANGE")},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("song_id"), AttributeType: aws.String("S")},
			{AttributeName: aws.String("id"), AttributeType: aws.String("S")},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5),
		},
	}

	_, err := svc.CreateTable(input)
	if err != nil && !strings.Contains(err.Error(), dynamodb.ErrCodeResourceInUseException) {
		logrus.WithError(err).Error("Failed to create DocumentsTable")
		return err
	}

	return waitForTableToBeActive(svc, bootstrap.DocumentTableName)
}

func waitForTableToBeActive(svc *dynamodb.DynamoDB, tableName string) error {
	return svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
}

func ClearTestTables(db *dynamo.DB) error {
	tables := []string{bootstrap.SongTableName, bootstrap.DocumentTableName}

	for _, tableName := range tables {
		table := db.Table(tableName)

		var items []map[string]interface{}
		if err := table.Scan().All(&items); err != nil {
			logrus.WithError(err).Errorf("Failed to scan items from table %s", tableName)
			return err
		}

		for _, item := range items {
			switch tableName {
			case bootstrap.SongTableName:
				if id, ok := item["id"].(string); ok {
					if err := table.Delete("id", id).Run(); err != nil {
						logrus.WithError(err).Warnf("Failed to delete item %s from %s", id, tableName)
					}
				}
			case bootstrap.DocumentTableName:
				songID, ok1 := item["song_id"].(string)
				id, ok2 := item["id"].(string)
				if ok1 && ok2 {
					if err := table.Delete("song_id", songID).Range("id", id).Run(); err != nil {
						logrus.WithError(err).Warnf("Failed to delete document %s for song %s", id, songID)
					}
				}
			}
		}
	}

	return nil
}
