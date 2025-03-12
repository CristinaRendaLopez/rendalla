package bootstrap

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/sirupsen/logrus"
)

var DB *dynamo.DB

func InitDB() {
	LoadConfig()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWSRegion),
	})
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to AWS")
	}

	DB = dynamo.New(sess)
	logrus.Info("Connected to DynamoDB successfully")
}
