package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

type DynamoSearchRepository struct {
	db      *dynamo.DB
	docRepo DocumentRepository
}

func NewDynamoSearchRepository(db *dynamo.DB, docRepo DocumentRepository) *DynamoSearchRepository {
	return &DynamoSearchRepository{
		db:      db,
		docRepo: docRepo,
	}
}

func (d *DynamoSearchRepository) SearchSongsByTitle(title string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error) {
	var songs []models.Song
	normalizedTitle := utils.Normalize(title)

	query := d.db.Table(bootstrap.SongTableName).
		Scan().
		Filter("contains(title_normalized, ?)", normalizedTitle).
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

func (d *DynamoSearchRepository) SearchDocumentsByTitle(title string, limit int, _ PagingKey) ([]models.Document, PagingKey, error) {
	var documents []models.Document
	var matchingSongs []models.Song

	normalizedTitle := utils.Normalize(title)
	songQuery := d.db.Table(bootstrap.SongTableName).
		Scan().
		Filter("contains(title_normalized, ?)", normalizedTitle).
		Limit(int64(limit))

	err := songQuery.All(&matchingSongs)
	if err != nil {
		logrus.WithFields(logrus.Fields{"title": title, "error": err}).Error("Failed to search songs by title")
		return nil, nil, err
	}

	for _, song := range matchingSongs {
		songDocs, err := d.docRepo.GetDocumentsBySongID(song.ID)
		if err != nil {
			logrus.WithFields(logrus.Fields{"song_id": song.ID, "error": err}).Warn("Failed to get documents for song")
			continue
		}
		documents = append(documents, songDocs...)
	}

	return documents, nil, nil
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
