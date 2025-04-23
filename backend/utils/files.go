package utils

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type FileUploader interface {
	UploadPDFToS3(file multipart.File, header *multipart.FileHeader, songID string) (string, error)
}

type FileService struct {
	client     *s3.Client
	bucketName string
}

func NewFileService() (*FileService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(bootstrap.AWSRegion))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	return &FileService{
		client:     client,
		bucketName: bootstrap.S3BucketName,
	}, nil
}

func (fs *FileService) UploadPDFToS3(file multipart.File, fileHeader *multipart.FileHeader, songID string) (string, error) {
	uid := uuid.NewString()
	key := fmt.Sprintf("songs/%s/doc_%s.pdf", songID, uid)

	_, err := fs.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &fs.bucketName,
		Key:         &key,
		Body:        file,
		ContentType: stringPtr("application/pdf"),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("uploading file to S3: %w", err)
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", fs.bucketName, bootstrap.AWSRegion, key)
	return url, nil
}

func stringPtr(s string) *string {
	return &s
}
