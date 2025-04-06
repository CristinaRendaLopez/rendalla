package integration_tests

import (
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

func SeedTestData(db *dynamo.DB, timeProvider utils.TimeProvider) error {
	now := timeProvider.Now()

	song := models.Song{
		ID:         "queen-001",
		Title:      "Bohemian Rhapsody",
		Author:     "Queen",
		Genres:     []string{"Rock", "Progressive"},
		YoutubeURL: "https://youtube.com/watch?v=fJ9rUzIMcZQ",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	docs := []models.Document{
		{
			ID:         "doc-br-piano",
			SongID:     song.ID,
			Type:       "partitura",
			Instrument: []string{"piano"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_piano.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "doc-br-voice",
			SongID:     song.ID,
			Type:       "tablatura",
			Instrument: []string{"voz"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_voz.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := db.Table(bootstrap.SongTableName).Put(song).Run(); err != nil {
		logrus.WithError(err).Error("Failed to insert test song")
		return err
	}

	for _, doc := range docs {
		if err := db.Table(bootstrap.DocumentTableName).Put(doc).Run(); err != nil {
			logrus.WithField("doc_id", doc.ID).WithError(err).Error("Failed to insert document")
			return err
		}
	}

	logrus.Info("Test data seeded")
	return nil
}
