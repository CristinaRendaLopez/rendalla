package repository

import (
	"time"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

type DynamoSongRepository struct {
	db      *dynamo.DB
	docRepo DocumentRepository
}

func NewDynamoSongRepository(db *dynamo.DB, docRepo DocumentRepository) *DynamoSongRepository {
	return &DynamoSongRepository{
		db:      db,
		docRepo: docRepo,
	}
}

func (d *DynamoSongRepository) CreateSongWithDocuments(song models.Song, documents []models.Document) error {
	var transactItems []*dynamodb.TransactWriteItem

	songItem, err := dynamodbattribute.MarshalMap(song)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal song item")
		return utils.ErrInternalServer
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
			logrus.WithError(err).
				WithField("doc_index", i).
				Error("Failed to marshal document item")
			return utils.ErrInternalServer
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
		logrus.WithError(err).Error("Failed to execute transactional write")
		return utils.HandleDynamoError(err)
	}

	logrus.WithField("song_id", song.ID).Info("Song and documents created transactionally")
	return nil
}

func (d *DynamoSongRepository) GetAllSongs() ([]models.Song, error) {
	var songs []models.Song
	err := d.db.Table(bootstrap.SongTableName).Scan().All(&songs)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve songs")
		return nil, utils.HandleDynamoError(err)
	}
	return songs, nil
}

func (d *DynamoSongRepository) GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	err := d.db.Table(bootstrap.SongTableName).Get("id", id).One(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to retrieve song")
		return nil, utils.HandleDynamoError(err)
	}
	logrus.WithField("song_id", id).Info("Song retrieved successfully")
	return &song, nil
}

func (d *DynamoSongRepository) UpdateSong(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := d.db.Table(bootstrap.SongTableName).Update("id", id)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to update song")
		return utils.HandleDynamoError(err)
	}
	logrus.WithField("song_id", id).Info("Song updated successfully")
	return nil
}

func (d *DynamoSongRepository) DeleteSongWithDocuments(songID string) error {
	_, err := d.GetSongByID(songID)
	if err != nil {
		logrus.WithField("song_id", songID).Warn("Attempted to delete a non-existing song")
		return err
	}

	documents, err := d.docRepo.GetDocumentsBySongID(songID)
	if err != nil && !utils.IsDynamoNotFoundError(err) {
		logrus.WithField("song_id", songID).WithError(err).Error("Failed to fetch documents before deletion")
		return utils.HandleDynamoError(err)
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
		logrus.WithError(err).Error("Failed to delete song and documents")
		return utils.HandleDynamoError(err)
	}

	logrus.WithField("song_id", songID).Info("Song and associated documents deleted successfully")
	return nil
}
