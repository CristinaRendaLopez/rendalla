package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	SongTableName     string
	DocumentTableName string
	AWSRegion         string
	AppPort           string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	SongTableName = getEnv("SONG_DYNAMODB_TABLE", "default_songs_table")
	DocumentTableName = getEnv("DOCUMENT_DYNAMODB_TABLE", "default_documents_table")
	AWSRegion = getEnv("AWS_REGION", "eu-north-1")
	AppPort = getEnv("APP_PORT", "8080")

	logrus.WithFields(logrus.Fields{
		"SongTableName":     SongTableName,
		"DocumentTableName": DocumentTableName,
		"AWSRegion":         AWSRegion,
		"AppPort":           AppPort,
	}).Info("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
