package services

import (
	"strings"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

func SearchSongsByTitle(title string, limit int, nextToken dynamo.PagingKey) ([]models.Song, dynamo.PagingKey, error) {
	title = strings.TrimSpace(strings.ToLower(title))
	if title == "" {
		return nil, nil, nil
	}

	var songs []models.Song
	query := bootstrap.DB.Table(bootstrap.SongTableName).
		Scan().
		Filter("contains(title, ?)", title).
		Limit(int64(limit))

	if nextToken != nil {
		query = query.StartFrom(nextToken)
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&songs)
	if err != nil {
		logrus.WithFields(logrus.Fields{"title": title, "error": err}).Error("Failed to search songs by title")
		return nil, nil, handleDynamoError(err)
	}

	return songs, nextKey, nil
}

func FilterDocumentsByInstrument(instrument string, limit int, nextToken dynamo.PagingKey) ([]models.Document, dynamo.PagingKey, error) {
	instrument = strings.TrimSpace(strings.ToLower(instrument))
	if instrument == "" {
		return nil, nil, nil
	}

	var documents []models.Document
	query := bootstrap.DB.Table(bootstrap.DocumentTableName).
		Scan().
		Filter("contains(instrument, ?)", instrument).
		Limit(int64(limit))

	if nextToken != nil {
		query = query.StartFrom(nextToken)
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&documents)
	if err != nil {
		logrus.WithFields(logrus.Fields{"instrument": instrument, "error": err}).Error("Failed to filter documents by instrument")
		return nil, nil, handleDynamoError(err)
	}

	return documents, nextKey, nil
}

func SearchDocumentsByTitle(title string, limit int, nextToken dynamo.PagingKey) ([]models.Document, dynamo.PagingKey, error) {
	title = strings.TrimSpace(strings.ToLower(title))
	if title == "" {
		return nil, nil, nil
	}

	var documents []models.Document
	query := bootstrap.DB.Table(bootstrap.DocumentTableName).
		Scan().
		Filter("contains(title, ?)", title).
		Limit(int64(limit))

	if nextToken != nil {
		query = query.StartFrom(nextToken)
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&documents)
	if err != nil {
		logrus.WithFields(logrus.Fields{"title": title, "error": err}).Error("Failed to search documents by title")
		return nil, nil, handleDynamoError(err)
	}

	return documents, nextKey, nil
}
