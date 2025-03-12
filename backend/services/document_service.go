package services

import (
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func GetDocumentsBySongID(songID string) ([]models.Document, error) {
	var documents []models.Document
	err := bootstrap.DB.Table(bootstrap.DocumentTableName).
		Scan().
		Filter("song_id = ?", songID).
		All(&documents)

	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": songID, "error": err}).Error("Failed to retrieve documents")
		return nil, handleDynamoError(err)
	}
	return documents, nil
}

func GetDocumentByID(id string) (*models.Document, error) {
	var document models.Document
	err := bootstrap.DB.Table(bootstrap.DocumentTableName).
		Get("id", id).
		One(&document)

	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Document not found")
		return nil, handleDynamoError(err)
	}
	return &document, nil
}

func CreateDocument(document models.Document) error {
	document.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	document.CreatedAt, document.UpdatedAt = now, now

	docItem, err := dynamodbattribute.MarshalMap(document)
	if err != nil {
		return err
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				ConditionCheck: &dynamodb.ConditionCheck{
					TableName:           aws.String(bootstrap.SongTableName),
					Key:                 map[string]*dynamodb.AttributeValue{"id": {S: aws.String(document.SongID)}},
					ConditionExpression: aws.String("attribute_exists(id)"),
				},
			},
			{
				Put: &dynamodb.Put{
					TableName: aws.String(bootstrap.DocumentTableName),
					Item:      docItem,
				},
			},
		},
	}

	_, err = bootstrap.DB.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": document.ID, "error": err}).Error("Failed to create document")
		return handleDynamoError(err)
	}

	logrus.WithField("document_id", document.ID).Info("Document created successfully")
	return nil
}

func UpdateDocument(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)

	update := bootstrap.DB.Table(bootstrap.DocumentTableName).Update("id", id)
	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Failed to update document")
		return handleDynamoError(err)
	}

	logrus.WithField("document_id", id).Info("Document updated successfully")
	return nil
}

func DeleteDocument(id string) error {
	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: []*dynamodb.TransactWriteItem{
			{
				ConditionCheck: &dynamodb.ConditionCheck{
					TableName:           aws.String(bootstrap.DocumentTableName),
					Key:                 map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
					ConditionExpression: aws.String("attribute_exists(id)"),
				},
			},
			{
				Delete: &dynamodb.Delete{
					TableName: aws.String(bootstrap.DocumentTableName),
					Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(id)}},
				},
			},
		},
	}

	_, err := bootstrap.DB.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Failed to delete document")
		return handleDynamoError(err)
	}

	logrus.WithField("document_id", id).Info("Document deleted successfully")
	return nil
}
