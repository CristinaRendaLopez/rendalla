package repository

import (
	stdErrors "errors"
	"fmt"
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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
//   - errors.ErrInternalServer if marshalling fails or the write operation fails.
func (d *DynamoDocumentRepository) CreateDocument(doc models.Document) error {

	if doc.ID == "" {
		logrus.WithField("operation", "create").Error("Missing document ID")
		return fmt.Errorf("document ID is required: %w", errors.ErrValidationFailed)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	doc.CreatedAt, doc.UpdatedAt = now, now

	docItem, err := dynamodbattribute.MarshalMap(doc)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"document_id": doc.ID,
			"song_id":     doc.SongID,
			"operation":   "create",
		}).WithError(err).Error("Failed to marshal document")
		return fmt.Errorf("failed to marshal document %s: %w", doc.ID, errors.ErrInternalServer)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(bootstrap.DocumentTableName),
		Item:      docItem,
	}

	_, err = d.db.Client().PutItem(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"document_id": doc.ID,
			"song_id":     doc.SongID,
			"operation":   "create",
		}).WithError(err).Error("Failed to create document in DynamoDB")
		return fmt.Errorf("creating document %s for song %s: %w", doc.ID, doc.SongID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"document_id": doc.ID,
		"song_id":     doc.SongID,
		"operation":   "create_document",
	}).Info("Document created successfully")
	return nil
}

// GetDocumentsBySongID retrieves all documents associated with the specified song ID.
// Returns:
//   - ([]models.Document, nil) on success
//   - (nil, errors.ErrInternalServer) if the query fails
func (d *DynamoDocumentRepository) GetDocumentsBySongID(songID string) ([]models.Document, error) {
	var documents []models.Document

	err := d.db.Table(bootstrap.DocumentTableName).
		Get("song_id", songID).
		All(&documents)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   songID,
			"operation": "get_all",
		}).WithError(err).Error("Failed to retrieve documents")
		return nil, fmt.Errorf("retrieving documents for song %s: %w", songID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"song_id":   songID,
		"operation": "get_all",
	}).Info("Documents retrieved successfully")
	return documents, nil
}

// GetDocumentByID retrieves a specific document by song ID and document ID.
// Returns:
//   - (*models.Document, nil) on success
//   - (nil, errors.ErrResourceNotFound) if the document does not exist
//   - (nil, errors.ErrInternalServer) if retrieval or unmarshalling fails
func (d *DynamoDocumentRepository) GetDocumentByID(songID string, docID string) (*models.Document, error) {
	var document models.Document

	err := d.db.Table(bootstrap.DocumentTableName).
		Get("song_id", songID).
		Range("id", dynamo.Equal, docID).
		One(&document)

	if err != nil {
		logFields := logrus.Fields{
			"song_id":     songID,
			"document_id": docID,
			"operation":   "get_by_id",
		}

		if stdErrors.Is(err, errors.ErrResourceNotFound) {
			logrus.WithFields(logFields).WithError(err).Warn("Document not found")
		} else {
			logrus.WithFields(logFields).WithError(err).Error("Failed to retrieve document")
		}

		return nil, fmt.Errorf("retrieving document %s for song %s: %w", docID, songID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"document_id": docID,
		"song_id":     songID,
		"operation":   "get_by_id",
	}).Info("Document retrieved successfully")

	return &document, nil
}

// UpdateDocument applies a partial update to the document identified by song ID and document ID.
// Automatically updates the updated_at timestamp.
// Returns:
//   - nil on success
//   - errors.ErrInternalServer if the update operation fails
func (d *DynamoDocumentRepository) UpdateDocument(songID, docID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := d.db.Table(bootstrap.DocumentTableName).Update("song_id", songID).Range("id", docID)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":     songID,
			"document_id": docID,
			"operation":   "update",
		}).WithError(err).Error("Failed to update document")
		return fmt.Errorf("updating document %s for song %s: %w", docID, songID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
		"operation":   "update_document",
	}).Info("Document updated successfully")
	return nil
}

// DeleteDocument removes a document identified by song ID and document ID from the DocumentTable.
// Returns:
//   - nil on success
//   - errors.ErrInternalServer if the delete operation fails
func (d *DynamoDocumentRepository) DeleteDocument(songID string, docID string) error {
	err := d.db.Table(bootstrap.DocumentTableName).
		Delete("song_id", songID).
		Range("id", docID).
		Run()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"document_id": docID,
			"song_id":     songID,
			"operation":   "delete",
		}).WithError(err).Error("Failed to delete document")
		return fmt.Errorf("deleting document %s for song %s: %w", docID, songID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"document_id": docID,
		"song_id":     songID,
		"operation":   "delete_document",
	}).Info("Document deleted successfully")
	return nil
}
