package bootstrap

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	SongTableName     string
	DocumentTableName string
	AWSRegion         string
	AppPort           string
	S3BucketName      string
	MaxPDFSize        int64
)

func LoadConfig() {

	SongTableName = getEnv("SONGS_TABLE", "default_songs_table")
	DocumentTableName = getEnv("DOCUMENTS_TABLE", "default_documents_table")
	AWSRegion = getEnv("AWS_REGION", "eu-north-1")
	AppPort = getEnv("APP_PORT", "8080")
	S3BucketName = getEnv("S3_BUCKET_NAME", "default-bucket")
	maxSizeStr := getEnv("MAX_PDF_SIZE", "10485760")
	maxSize, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil {
		logrus.WithError(err).Warn("Invalid MAX_PDF_SIZE, defaulting to 10MB")
		maxSize = 10 * 1024 * 1024
	}
	MaxPDFSize = maxSize

	logrus.WithFields(logrus.Fields{
		"SongTableName":     SongTableName,
		"DocumentTableName": DocumentTableName,
		"AWSRegion":         AWSRegion,
		"AppPort":           AppPort,
		"S3BucketName":      S3BucketName,
		"MaxPDFSize":        MaxPDFSize,
	}).Info("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
