package repository

import (
	"sort"

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

func (d *DynamoSearchRepository) ListSongs(title, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Song, PagingKey, error) {
	var songs []models.Song

	query := d.db.Table(bootstrap.SongTableName).Scan().Limit(int64(limit))

	if title != "" {
		normalizedTitle := utils.Normalize(title)
		query = query.Filter("contains(title_normalized, ?)", normalizedTitle)
	}

	if nextToken != nil {
		if nt, ok := nextToken.(dynamo.PagingKey); ok {
			query = query.StartFrom(nt)
		}
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&songs)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"title": title,
			"error": err,
		}).Error("Failed to list songs")
		return nil, nil, err
	}

	if sortField == "" {
		sortField = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.SliceStable(songs, func(i, j int) bool {
		switch sortField {
		case "title":
			if sortOrder == "asc" {
				return songs[i].Title < songs[j].Title
			}
			return songs[i].Title > songs[j].Title
		case "created_at":
			if sortOrder == "asc" {
				return songs[i].CreatedAt < songs[j].CreatedAt
			}
			return songs[i].CreatedAt > songs[j].CreatedAt
		default:
			return true
		}
	})

	return songs, nextKey, nil
}

func (d *DynamoSearchRepository) ListDocuments(title, instrument, docType, sortField, sortOrder string, limit int, nextToken PagingKey) ([]models.Document, PagingKey, error) {
	var documents []models.Document

	query := d.db.Table(bootstrap.DocumentTableName).Scan().Limit(int64(limit))

	if title != "" {
		normalizedTitle := utils.Normalize(title)
		query = query.Filter("contains(title_normalized, ?)", normalizedTitle)
	}
	if instrument != "" {
		query = query.Filter("contains(instrument, ?)", instrument)
	}
	if docType != "" {
		query = query.Filter("'type' = ?", docType)
	}
	if nextToken != nil {
		if nt, ok := nextToken.(dynamo.PagingKey); ok {
			query = query.StartFrom(nt)
		}
	}

	nextKey, err := query.AllWithLastEvaluatedKey(&documents)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"title":      title,
			"instrument": instrument,
			"type":       docType,
			"error":      err,
		}).Error("Failed to list documents")
		return nil, nil, err
	}

	if sortField == "" {
		sortField = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	sort.SliceStable(documents, func(i, j int) bool {
		switch sortField {
		case "title":
			if sortOrder == "asc" {
				return documents[i].TitleNormalized < documents[j].TitleNormalized
			}
			return documents[i].TitleNormalized > documents[j].TitleNormalized
		case "created_at":
			if sortOrder == "asc" {
				return documents[i].CreatedAt < documents[j].CreatedAt
			}
			return documents[i].CreatedAt > documents[j].CreatedAt
		default:
			return true
		}
	})

	return documents, nextKey, nil
}
