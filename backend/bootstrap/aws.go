package bootstrap

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func InitAWS() *dynamodb.DynamoDB {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("Missing required environment variable: AWS_REGION")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Error initializing AWS session: %v", err)
	}

	log.Println("AWS session initialized successfully")
	return dynamodb.New(sess)
}
