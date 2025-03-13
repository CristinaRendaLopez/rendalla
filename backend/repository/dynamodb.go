package repository

import (
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/sirupsen/logrus"
)

type DynamoDB struct{}

func (d *DynamoDB) PutSong(song models.Song) error {
	err := bootstrap.DB.Table(bootstrap.SongTableName).Put(song).Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": song.ID, "error": err}).Error("Failed to insert song into DynamoDB")
		return err
	}
	logrus.WithField("song_id", song.ID).Info("Song inserted successfully")
	return nil
}

func (d *DynamoDB) GetSongByID(id string) (*models.Song, error) {
	var song models.Song
	err := bootstrap.DB.Table(bootstrap.SongTableName).Get("id", id).One(&song)
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to retrieve song")
		return nil, err
	}
	logrus.WithField("song_id", id).Info("Song retrieved successfully")
	return &song, nil
}

func (d *DynamoDB) UpdateSong(id string, updates map[string]interface{}) error {
	update := bootstrap.DB.Table(bootstrap.SongTableName).Update("id", id)
	for key, value := range updates {
		update = update.Set(key, value)
	}

	err := update.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to update song")
		return err
	}
	logrus.WithField("song_id", id).Info("Song updated successfully")
	return nil
}

func (d *DynamoDB) DeleteSong(id string) error {
	err := bootstrap.DB.Table(bootstrap.SongTableName).Delete("id", id).Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": id, "error": err}).Error("Failed to delete song")
		return err
	}
	logrus.WithField("song_id", id).Info("Song deleted successfully")
	return nil
}
