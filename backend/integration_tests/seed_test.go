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

	bohemianRhapsody := models.Song{
		ID:              "queen-001",
		Title:           "Bohemian Rhapsody",
		TitleNormalized: utils.Normalize("Bohemian Rhapsody"),
		Author:          "Queen",
		Genres:          []string{"Rock", "Progressive"},
		YoutubeURL:      "https://youtube.com/watch?v=fJ9rUzIMcZQ",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	bohemianDocuments := []models.Document{
		{
			ID:         "doc-br-piano",
			SongID:     bohemianRhapsody.ID,
			Type:       "score",
			Instrument: []string{"piano"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_piano.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "doc-br-voice",
			SongID:     bohemianRhapsody.ID,
			Type:       "tablatura",
			Instrument: []string{"voz"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_voz.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	dontStopMeNow := models.Song{
		ID:              "queen-002",
		Title:           "Don't Stop Me Now",
		TitleNormalized: utils.Normalize("Don't Stop Me Now"),
		Author:          "Queen",
		Genres:          []string{"Rock", "Pop"},
		YoutubeURL:      "https://youtube.com/watch?v=HgzGwKwLmgM",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	dontStopMeNowDocs := []models.Document{
		{
			ID:         "doc-dsmn-guitar",
			SongID:     dontStopMeNow.ID,
			Type:       "score",
			Instrument: []string{"guitar"},
			PDFURL:     "https://s3.test/dontstop_guitar.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := db.Table(bootstrap.SongTableName).Put(bohemianRhapsody).Run(); err != nil {
		logrus.WithError(err).Error("Failed to insert Bohemian Rhapsody")
		return err
	}
	for _, doc := range bohemianDocuments {
		if err := db.Table(bootstrap.DocumentTableName).Put(doc).Run(); err != nil {
			logrus.WithField("doc_id", doc.ID).WithError(err).Error("Failed to insert Bohemian document")
			return err
		}
	}

	if err := db.Table(bootstrap.SongTableName).Put(dontStopMeNow).Run(); err != nil {
		logrus.WithError(err).Error("Failed to insert Don't Stop Me Now")
		return err
	}
	for _, doc := range dontStopMeNowDocs {
		if err := db.Table(bootstrap.DocumentTableName).Put(doc).Run(); err != nil {
			logrus.WithField("doc_id", doc.ID).WithError(err).Error("Failed to insert Don't Stop Me Now document")
			return err
		}
	}

	logrus.Info("Test data seeded")
	return nil
}
