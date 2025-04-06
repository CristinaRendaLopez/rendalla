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

	// Intentar borrar si existen
	for _, table := range []string{bootstrap.SongTableName, bootstrap.DocumentTableName} {
		_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(table)})
		if err != nil {
			logrus.WithError(err).Warnf("Could not delete table %s", table)
		}
	}

	tables := []struct {
		Name       string
		PrimaryKey string
	}{
		{bootstrap.SongTableName, "id"},
		{bootstrap.DocumentTableName, "id"},
	}

	for _, t := range tables {
		_, err := svc.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(t.Name),
			KeySchema: []*dynamodb.KeySchemaElement{
				{AttributeName: aws.String(t.PrimaryKey), KeyType: aws.String("HASH")},
			},
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{AttributeName: aws.String(t.PrimaryKey), AttributeType: aws.String("S")},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits: aws.Int64(5), WriteCapacityUnits: aws.Int64(5),
			},
		})

		if err != nil && !strings.Contains(err.Error(), dynamodb.ErrCodeResourceInUseException) {
			logrus.WithField("table", t.Name).WithError(err).Error("Failed to create table")
			return err
		}

		if err := waitForTableToBeActive(svc, t.Name); err != nil {
			logrus.WithError(err).Error("Timed out waiting for table to be active")
			return err
		}

		logrus.WithField("table", t.Name).Info("Table is ready")
	}

	return nil
}

func waitForTableToBeActive(svc *dynamodb.DynamoDB, tableName string) error {
	return svc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
}
