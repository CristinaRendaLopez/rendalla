package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

type DynamoSearchRepository struct {
	db *dynamo.DB
}

func NewDynamoSearchRepository(db *dynamo.DB) *DynamoSearchRepository {
	return &DynamoSearchRepository{db: db}
}

func (d *DynamoSearchRepository) SearchSongsByTitle(title string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error) {
	var songs []models.Song
	query := d.db.Table(bootstrap.SongTableName).
		Scan().
		Filter("contains(title, ?)", title).
		Limit(int64(limit))

	if nextToken != nil {
		if nt, ok := nextToken.(dynamo.PagingKey); ok {
			query = query.StartFrom(nt)
		}
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&songs)
	if err != nil {
		logrus.WithFields(logrus.Fields{"title": title, "error": err}).Error("Failed to search songs by title")
		return nil, nil, err
	}

	return songs, nextKey, nil
}

func (d *DynamoSearchRepository) SearchDocumentsByTitle(title string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error) {
	var documents []models.Document
	query := d.db.Table(bootstrap.DocumentTableName).
		Scan().
		Filter("contains(title, ?)", title).
		Limit(int64(limit))

	if nextToken != nil {
		if nt, ok := nextToken.(dynamo.PagingKey); ok {
			query = query.StartFrom(nt)
		}
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&documents)
	if err != nil {
		logrus.WithFields(logrus.Fields{"title": title, "error": err}).Error("Failed to search documents by title")
		return nil, nil, err
	}

	return documents, nextKey, nil
}

func (d *DynamoSearchRepository) FilterDocumentsByInstrument(instrument string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error) {
	var documents []models.Document
	query := d.db.Table(bootstrap.DocumentTableName).
		Scan().
		Filter("contains(instrument, ?)", instrument).
		Limit(int64(limit))

	if nextToken != nil {
		if nt, ok := nextToken.(dynamo.PagingKey); ok {
			query = query.StartFrom(nt)
		}
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&documents)
	if err != nil {
		logrus.WithFields(logrus.Fields{"instrument": instrument, "error": err}).Error("Failed to filter documents by instrument")
		return nil, nil, err
	}

	return documents, nextKey, nil
}
