package repository

import (
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

type DynamoDocumentRepository struct {
	db *dynamo.DB
}

func NewDynamoDocumentRepository(db *dynamo.DB) *DynamoDocumentRepository {
	return &DynamoDocumentRepository{db: db}
}

func (d *DynamoDocumentRepository) GetDocumentByID(id string) (*models.Document, error) {
	var document models.Document
	err := d.db.Table(bootstrap.DocumentTableName).Get("id", id).One(&document)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Document not found")
		return nil, utils.HandleDynamoError(err)
	}
	return &document, nil
}

func (d *DynamoDocumentRepository) CreateDocument(doc models.Document) (string, error) {
	doc.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	doc.CreatedAt, doc.UpdatedAt = now, now

	docItem, err := dynamodbattribute.MarshalMap(doc)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(bootstrap.DocumentTableName),
		Item:      docItem,
	}

	_, err = d.db.Client().PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to create document")
		return "", utils.HandleDynamoError(err)
	}

	logrus.WithField("document_id", doc.ID).Info("Document created successfully")
	return doc.ID, nil
}

func (d *DynamoDocumentRepository) UpdateDocument(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := d.db.Table(bootstrap.DocumentTableName).Update("id", id)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Failed to update document")
		return utils.HandleDynamoError(err)
	}

	logrus.WithField("document_id", id).Info("Document updated successfully")
	return nil
}

func (d *DynamoDocumentRepository) DeleteDocument(id string) error {
	err := d.db.Table(bootstrap.DocumentTableName).Delete("id", id).Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": id, "error": err}).Error("Failed to delete document")
		return utils.HandleDynamoError(err)
	}

	logrus.WithField("document_id", id).Info("Document deleted successfully")
	return nil
}

func (d *DynamoDocumentRepository) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	var documents []models.Document
	err := d.db.Table(bootstrap.DocumentTableName).Scan().Filter("song_id = ?", songID).All(&documents)

	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": songID, "error": err}).Error("Failed to retrieve documents")
		return nil, utils.HandleDynamoError(err)
	}

	return documents, nil
}
