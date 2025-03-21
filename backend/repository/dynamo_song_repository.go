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
	"github.com/sirupsen/logrus"
)

type DynamoSongRepository struct {
	docRepo DocumentRepository
}

func NewDynamoSongRepository(docRepo DocumentRepository) *DynamoSongRepository {
	return &DynamoSongRepository{docRepo: docRepo}
}

func (d *DynamoSongRepository) CreateSongWithDocuments(song models.Song, documents []models.Document) (string, error) {
	song.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	song.CreatedAt, song.UpdatedAt = now, now

	songItem, err := dynamodbattribute.MarshalMap(song)
	if err != nil {
		return "", err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(bootstrap.SongTableName),
		Item:      songItem,
	}

	_, err = bootstrap.DB.Client().PutItem(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to create song")
		return "", utils.HandleDynamoError(err)
	}

	for i := range documents {
		documents[i].SongID = song.ID
		_, err := d.docRepo.CreateDocument(documents[i])
		if err != nil {
			logrus.WithError(err).Error("Failed to create document for song")
			return "", err
		}
	}

	logrus.WithField("song_id", song.ID).Info("Song and documents created successfully")
	return song.ID, nil
}

func (d *DynamoSongRepository) GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	err := bootstrap.DB.Table(bootstrap.SongTableName).Get("id", id).One(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to retrieve song")
		return nil, utils.HandleDynamoError(err)
	}
	logrus.WithField("song_id", id).Info("Song retrieved successfully")
	return &song, nil
}

func (d *DynamoSongRepository) UpdateSong(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := bootstrap.DB.Table(bootstrap.SongTableName).Update("id", id)

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
		return err
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
					Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(doc.ID)}},
				},
			})
		}
	}

	input := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	_, err = bootstrap.DB.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song and documents")
		return err
	}

	logrus.WithField("song_id", songID).Info("Song and associated documents deleted successfully")
	return nil
}

func (d *DynamoSongRepository) GetAllSongs() ([]models.Song, error) {
	var songs []models.Song
	err := bootstrap.DB.Table(bootstrap.SongTableName).Scan().All(&songs)
	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve songs")
		return nil, utils.HandleDynamoError(err)
	}
	return songs, nil
}
