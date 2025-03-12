package bootstrap

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Configuration variables
var (
	SongTableName     string
	DocumentTableName string
	AWSRegion         string
	AppPort           string
)

// Load configuration from .env file or system environment variables
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	SongTableName = getEnv("SONG_DYNAMODB_TABLE", "default_songs_table")
	DocumentTableName = getEnv("DOCUMENT_DYNAMODB_TABLE", "default_documents_table")
	AWSRegion = getEnv("AWS_REGION", "us-east-1")
	AppPort = getEnv("APP_PORT", "8080")

	logrus.WithFields(logrus.Fields{
		"SongTableName":     SongTableName,
		"DocumentTableName": DocumentTableName,
		"AWSRegion":         AWSRegion,
		"AppPort":           AppPort,
	}).Info("Configuration loaded successfully")
}

// Helper function to get environment variables with a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
