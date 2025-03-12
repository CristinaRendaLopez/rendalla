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

func GetAllSongs() ([]models.Song, error) {
	var songs []models.Song
	err := bootstrap.DB.Table(bootstrap.SongTableName).Scan().All(&songs)

	if err != nil {
		logrus.WithError(err).Error("Failed to retrieve songs")
		return nil, handleDynamoError(err)
	}
	return songs, nil
}

func GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	err := bootstrap.DB.Table(bootstrap.SongTableName).Get("id", id).One(&song)

	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Song not found")
		return nil, handleDynamoError(err)
	}
	return &song, nil
}

func CreateSongWithDocuments(song models.Song, documents []models.Document) error {
	song.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	song.CreatedAt, song.UpdatedAt = now, now

	songItem, err := dynamodbattribute.MarshalMap(song)
	if err != nil {
		return err
	}

	var transactItems []*dynamodb.TransactWriteItem
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Put: &dynamodb.Put{
			TableName: aws.String(bootstrap.SongTableName),
			Item:      songItem,
		},
	})

	for i := range documents {
		documents[i].ID = uuid.New().String()
		documents[i].SongID = song.ID
		documents[i].CreatedAt, documents[i].UpdatedAt = now, now

		docItem, err := dynamodbattribute.MarshalMap(documents[i])
		if err != nil {
			return err
		}

		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				TableName: aws.String(bootstrap.DocumentTableName),
				Item:      docItem,
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	_, err = bootstrap.DB.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to create song with documents")
		return handleDynamoError(err)
	}

	logrus.WithField("song_id", song.ID).Info("Song and documents created successfully")
	return nil
}

func UpdateSong(id string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now().UTC().Format(time.RFC3339)
	update := bootstrap.DB.Table(bootstrap.SongTableName).Update("id", id)

	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to update song")
		return handleDynamoError(err)
	}

	logrus.WithField("song_id", id).Info("Song updated successfully")
	return nil
}

func DeleteSongWithDocuments(songID string) error {
	_, err := GetSongByID(songID)
	if err != nil {
		logrus.WithField("song_id", songID).Warn("Attempted to delete a non-existing song")
		return handleDynamoError(err)
	}

	documents, err := GetDocumentsBySongID(songID)
	if err != nil {
		return handleDynamoError(err)
	}

	var transactItems []*dynamodb.TransactWriteItem
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{
		Delete: &dynamodb.Delete{
			TableName: aws.String(bootstrap.SongTableName),
			Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(songID)}},
		},
	})

	for _, doc := range documents {
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Delete: &dynamodb.Delete{
				TableName: aws.String(bootstrap.DocumentTableName),
				Key:       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(doc.ID)}},
			},
		})
	}

	input := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	_, err = bootstrap.DB.Client().TransactWriteItems(input)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete song and documents")
		return handleDynamoError(err)
	}

	logrus.WithField("song_id", songID).Info("Song and associated documents deleted successfully")
	return nil
}
