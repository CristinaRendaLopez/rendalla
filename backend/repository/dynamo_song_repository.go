package repository

import (
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

// DynamoSongRepository implements SongRepository using DynamoDB as backend.
// Songs are stored in the "SongTable", and documents are stored separately in the "DocumentTable".
// Each song can be created or deleted transactionally along with its associated documents.
type DynamoSongRepository struct {
	db      *dynamo.DB
	docRepo DocumentRepository
}

// NewDynamoSongRepository returns a new instance of DynamoSongRepository.
func NewDynamoSongRepository(db *dynamo.DB, docRepo DocumentRepository) *DynamoSongRepository {
	return &DynamoSongRepository{
		db:      db,
		docRepo: docRepo,
	}
}

// CreateSongWithDocuments stores a new song and its associated documents in a single transactional write.
// Returns errors.ErrInternalServer on marshalling errors or any write failure.
func (d *DynamoSongRepository) CreateSongWithDocuments(song models.Song, documents []models.Document) error {
	var transactItems []*dynamodb.TransactWriteItem

	songItem, err := dynamodbattribute.MarshalMap(song)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   song.ID,
			"operation": "create",
		}).WithError(err).Error("Failed to marshal song item")
		return fmt.Errorf("failed to marshal song %s: %w", song.ID, errors.ErrInternalServer)
	}

	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName: aws.String(bootstrap.SongTableName),
			Item:      songItem,
		},
	})

	for i, doc := range documents {
		docItem, err := dynamodbattribute.MarshalMap(doc)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"song_id":   song.ID,
				"doc_index": i,
				"doc_id":    doc.ID,
				"operation": "create",
			}).WithError(err).Error("Failed to marshal document item")
			return fmt.Errorf("failed to marshal document %d (doc_id=%s) for song %s: %w", i, doc.ID, song.ID, errors.ErrInternalServer)
		}

		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(bootstrap.DocumentTableName),
				Item:      docItem,
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	}

	_, err = d.db.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   song.ID,
			"operation": "create",
		}).WithError(err).Error("Failed to execute transactional write")
		return fmt.Errorf("transactional creation of song %s failed: %w", song.ID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"song_id":   song.ID,
		"operation": "create_song",
	}).Info("Song and documents created transactionally")
	return nil
}

// GetAllSongs retrieves all songs from the SongTable.
// Returns a list of songs or an internal error if the query fails.
func (d *DynamoSongRepository) GetAllSongs() ([]models.Song, error) {
	var songs []models.Song
	err := d.db.Table(bootstrap.SongTableName).Scan().All(&songs)
	if err != nil {
		logrus.WithField("operation", "get_all").WithError(err).Error("Failed to retrieve songs")
		return nil, fmt.Errorf("retrieving all songs: %w", errors.HandleDynamoError(err))
	}
	return songs, nil
}

// GetSongByID retrieves a song by its ID.
// Returns:
//   - (*models.Song, nil) on success
//   - (nil, errors.ErrResourceNotFound) if the song does not exist
//   - (nil, errors.ErrInternalServer) for marshalling or database access errors
func (d *DynamoSongRepository) GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	err := d.db.Table(bootstrap.SongTableName).Get("id", id).One(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   id,
			"operation": "get_by_id",
		}).WithError(err).Warn("Song not found or failed to retrieve")
		return nil, fmt.Errorf("retrieving song %s: %w", id, errors.HandleDynamoError(err))
	}
	logrus.WithFields(logrus.Fields{
		"song_id":   id,
		"operation": "get_song",
	}).Info("Song retrieved successfully")
	return &song, nil
}

// UpdateSong applies partial updates to a song by its ID.
// Automatically sets the updated_at field to the current timestamp.
// Returns an error if the update fails or the item does not exist.
func (d *DynamoSongRepository) UpdateSong(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := d.db.Table(bootstrap.SongTableName).Update("id", id)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   id,
			"operation": "update",
		}).WithError(err).Error("Failed to update song")
		return fmt.Errorf("updating song %s: %w", id, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"song_id":   id,
		"operation": "update_song",
	}).Info("Song updated successfully")
	return nil
}

// DeleteSongWithDocuments removes a song and all of its associated documents in a single transaction.
// Returns:
//   - errors.ErrResourceNotFound if the song does not exist
//   - errors.ErrInternalServer if deletion fails at any point
func (d *DynamoSongRepository) DeleteSongWithDocuments(songID string) error {
	_, err := d.GetSongByID(songID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   songID,
			"operation": "delete",
		}).WithError(err).Warn("Cannot delete song: not found or retrieval failed")
		return fmt.Errorf("verifying existence of song %s before delete: %w", songID, err)
	}

	documents, err := d.docRepo.GetDocumentsBySongID(songID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   songID,
			"operation": "delete",
		}).WithError(err).Error("Failed to retrieve documents before deletion")
		return fmt.Errorf("retrieving documents for song %s: %w", songID, errors.HandleDynamoError(err))
	}

	var transactItems []*dynamodb.TransactWriteItem

	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Delete: &dynamodb.Delete{
			TableName: aws.String(bootstrap.SongTableName),
			Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(songID)}},
		},
	})

	if len(documents) > 0 {
		for _, doc := range documents {
			transactItems = append(transactItems, &dynamodb.TransactWriteItem{
				Delete: &dynamodb.Delete{
					TableName: aws.String(bootstrap.DocumentTableName),
					Key: map[string]*dynamodb.AttributeValue{
						"song_id": {S: aws.String(doc.SongID)},
						"id":      {S: aws.String(doc.ID)},
					},
				},
			})
		}
	}

	input := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	_, err = d.db.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id":   songID,
			"operation": "delete",
		}).WithError(err).Error("Failed to delete song and documents")
		return fmt.Errorf("deleting song %s transactionally: %w", songID, errors.HandleDynamoError(err))
	}

	logrus.WithFields(logrus.Fields{
		"song_id":   songID,
		"operation": "delete_song",
	}).Info("Song and associated documents deleted successfully")
	return nil
}
