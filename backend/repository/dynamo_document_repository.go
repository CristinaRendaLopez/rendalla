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

// DynamoDocumentRepository implements DocumentRepository using Amazon DynamoDB as the storage layer.
// Documents are stored in a table with a composite primary key: (song_id, id).
// This allows efficient access to all documents for a given song, and supports direct lookup by document ID.
type DynamoDocumentRepository struct {
	db *dynamo.DB
}

// NewDynamoDocumentRepository returns a new instance of DynamoDocumentRepository.
func NewDynamoDocumentRepository(db *dynamo.DB) *DynamoDocumentRepository {
	return &DynamoDocumentRepository{db: db}
}

// CreateDocument inserts a new document into the DocumentTable.
// It generates a UUID for the document ID and sets creation/update timestamps.
// Returns:
//   - utils.ErrInternalServer if marshalling fails or the write operation fails.
func (d *DynamoDocumentRepository) CreateDocument(doc models.Document) error {
	doc.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	doc.CreatedAt, doc.UpdatedAt = now, now

	docItem, err := dynamodbattribute.MarshalMap(doc)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal document")
		return utils.ErrInternalServer
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(bootstrap.DocumentTableName),
		Item:      docItem,
	}

	_, err = d.db.Client().PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to create document in DynamoDB")
		return utils.HandleDynamoError(err)
	}

	logrus.WithField("document_id", doc.ID).Info("Document created successfully")
	return nil
}

// GetDocumentsBySongID retrieves all documents associated with the specified song ID.
// Returns:
//   - ([]models.Document, nil) on success
//   - (nil, utils.ErrInternalServer) if the query fails
func (d *DynamoDocumentRepository) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	var documents []models.Document

	err := d.db.Table(bootstrap.DocumentTableName).
		Get("song_id", songID).
		All(&documents)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": songID,
			"error":   err,
		}).Error("Failed to retrieve documents")
		return nil, utils.HandleDynamoError(err)
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	return documents, nil
}

// GetDocumentByID retrieves a specific document by song ID and document ID.
// Returns:
//   - (*models.Document, nil) on success
//   - (nil, utils.ErrNotFound) if the document does not exist
//   - (nil, utils.ErrInternalServer) if retrieval or unmarshalling fails
func (d *DynamoDocumentRepository) GetDocumentByID(songID string, docID string) (*models.Document, error) {
	var document models.Document

	err := d.db.Table(bootstrap.DocumentTableName).
		Get("song_id", songID).
		Range("id", dynamo.Equal, docID).
		One(&document)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":     songID,
			"document_id": docID,
			"error":       err,
		}).Error("Failed to retrieve document")

		return nil, utils.HandleDynamoError(err)
	}

	logrus.WithField("document_id", docID).Info("Document retrieved successfully")
	return &document, nil
}

// UpdateDocument applies a partial update to the document identified by song ID and document ID.
// Automatically updates the updated_at timestamp.
// Returns:
//   - nil on success
//   - utils.ErrInternalServer if the update operation fails
func (d *DynamoDocumentRepository) UpdateDocument(songID, docID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := d.db.Table(bootstrap.DocumentTableName).Update("song_id", songID).Range("id", docID)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": songID, "document_id": docID, "error": err}).Error("Failed to update document")
		return utils.HandleDynamoError(err)
	}

	logrus.WithFields(logrus.Fields{"song_id": songID, "document_id": docID}).Info("Document updated successfully")
	return nil
}

// DeleteDocument removes a document identified by song ID and document ID from the DocumentTable.
// Returns:
//   - nil on success
//   - utils.ErrInternalServer if the delete operation fails
func (d *DynamoDocumentRepository) DeleteDocument(songID string, docID string) error {
	err := d.db.Table(bootstrap.DocumentTableName).
		Delete("song_id", songID).
		Range("id", docID).
		Run()

	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": docID, "song_id": songID, "error": err}).Error("Failed to delete document")
		return utils.HandleDynamoError(err)
	}

	logrus.WithFields(logrus.Fields{"document_id": docID, "song_id": songID}).Info("Document deleted successfully")
	return nil
}
